package cronjobs

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	domains "github.com/skobelina/currency_converter/internal"
	"github.com/skobelina/currency_converter/internal/subscribers"
)

type Event struct {
	EventID     string    `json:"eventId"`
	EventType   string    `json:"eventType"`
	AggregateID string    `json:"aggregateId"`
	Timestamp   string    `json:"timestamp"`
	Data        EventData `json:"data"`
}

type EventData struct {
	CreatedAt    string   `json:"createdAt"`
	ExchangeRate string   `json:"exchangeRate"`
	Recipients   []string `json:"recipients"`
}

type ExchangeRateTemplate struct {
	CreatedAt    string
	ExchangeRate string
}

func (s *CronJobService) NotificationExchangeRates() error {
	subscribersResp, err := s.subscribers.Search(&subscribers.SearchSubscribeRequest{
		Filter: domains.DefaultFilter(),
	})
	if err != nil {
		return err
	}
	recipients := make([]string, 0, len(subscribersResp.Data))
	for _, person := range subscribersResp.Data {
		recipients = append(recipients, person.Email)
	}
	currentTime := time.Now().Format("2006-01-02")

	exchangeRate, err := s.rates.Get()
	if err != nil {
		return err
	}
	template := ExchangeRateTemplate{
		CreatedAt:    currentTime,
		ExchangeRate: strconv.FormatFloat(*exchangeRate, 'f', 2, 64),
	}

	event := Event{
		EventID:     uuid.New().String(),
		EventType:   "CurrencyRate",
		AggregateID: "rate-" + uuid.New().String(),
		Timestamp:   time.Now().Format(time.RFC3339),
		Data: EventData{
			CreatedAt:    currentTime,
			ExchangeRate: template.ExchangeRate,
			Recipients:   recipients,
		},
	}
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error marshalling event: %v", err)
		return err
	}
	if err := s.rabbitMQ.PublishMessage(string(body)); err != nil {
		logrus.Errorf("Error publishing message: %v", err)
		return err
	}

	logrus.Infof("CronJob: NotificationExchangeRates: sent to %d subscribers", len(recipients))
	return nil
}
