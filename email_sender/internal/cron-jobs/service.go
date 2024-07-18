package cronjobs

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/email_sender/internal/mails"
	"github.com/skobelina/email_sender/internal/mails/templates"
	"github.com/skobelina/email_sender/pkg/queue"
)

type CronJobService struct {
	mailService mails.MailServiceInterface
	rabbitMQ    queue.Queue
	repo        Repository
}

func NewCronJobService(mailService mails.MailServiceInterface, rabbitMQ queue.Queue, repo Repository) *CronJobService {
	return &CronJobService{
		mailService: mailService,
		rabbitMQ:    rabbitMQ,
		repo:        repo,
	}
}

func (s *CronJobService) ConsumeMessages() {
	msgs, err := s.rabbitMQ.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	for msg := range msgs {
		var event Event
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			logrus.Errorf("Error unmarshalling message: %v", err)
			continue
		}

		switch event.EventType {
		case "Subscribe":
			subscriber := &Subscriber{Email: event.Data.Email}
			if err := s.repo.Create(subscriber); err != nil {
				logrus.Errorf("Error inserting subscriber: %v", err)
			} else {
				logrus.Infof("Successfully inserted subscriber: %s", event.Data.Email)
			}
		case "Unsubscribe":
			subscriber, err := s.repo.FindByEmail(event.Data.Email)
			if err != nil {
				logrus.Errorf("Error finding subscriber: %v", err)
				continue
			}
			if err := s.repo.Delete(subscriber); err != nil {
				logrus.Errorf("Error deleting subscriber: %v", err)
			} else {
				logrus.Infof("Successfully deleted subscriber: %s", event.Data.Email)
			}
		case "CurrencyRate":
			subscribers, err := s.repo.Search()
			if err != nil {
				logrus.Errorf("Error fetching subscribers from database: %v", err)
				continue
			}

			if len(subscribers) == 0 {
				logrus.Infof("No subscribers found, skipping email sending")
				continue
			}

			var recipients []string
			for _, subscriber := range subscribers {
				recipients = append(recipients, subscriber.Email)
			}

			template := templates.ExchangeRateTemplate{
				CreatedAt:    event.Data.CreatedAt,
				ExchangeRate: event.Data.ExchangeRate,
			}

			err = s.mailService.SendEmail(recipients, "Exchange rates notification", template)
			if err != nil {
				logrus.Errorf("Error sending email: %v", err)
			}
			logrus.Infof("CronJob: NotificationExchangeRates: sent to %d subscribers", len(recipients))
		default:
			logrus.Infof("Ignoring event of type: %s", event.EventType)
		}
	}
}

func (s *CronJobService) InitializeSubscribers() {
	logrus.Info("Starting to initialize subscribers")
	resp, err := http.Get("http://localhost:8080/api/subscribe")
	if err != nil {
		logrus.Errorf("Error fetching subscribers: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf("Error closing response body: %v", err)
		}
	}()
	logrus.Infof("Received response with status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	var subscribersResp struct {
		Data []struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&subscribersResp); err != nil {
		logrus.Errorf("Error decoding subscribers response: %v", err)
		return
	}

	if err := s.repo.DeleteAll(); err != nil {
		logrus.Errorf("Error clearing subscriber database: %v", err)
		return
	}

	for _, subscriberData := range subscribersResp.Data {
		subscriber := &Subscriber{Email: subscriberData.Email}
		if err := s.repo.Create(subscriber); err != nil {
			logrus.Errorf("Error inserting subscriber: %v", err)
		} else {
			logrus.Infof("Successfully inserted subscriber: %s", subscriberData.Email)
		}
	}
	logrus.Info("Finished initializing subscribers")
}
