package main

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/email_sender/internal/mails"
	"github.com/skobelina/email_sender/internal/mails/templates"
	"github.com/skobelina/email_sender/pkg/queue"
)

var (
	rabbitMQURL = os.Getenv("RABBITMQ_URL")
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

func main() {
	mailPassword := os.Getenv("MAILPASS")
	if mailPassword == "" {
		logrus.Fatal("MAILPASS environment variable is required")
	}
	emailSender := mails.NewService("marisa.skobelina@gmail.com", "smtp.gmail.com")

	rabbitMQ, err := queue.NewRabbitMQ(rabbitMQURL, "events")
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.ConsumeMessages()
	if err != nil {
		logrus.Fatalf("Failed to register a consumer: %v", err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logrus.Infof("Received a message: %s", d.Body)

			var event Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				logrus.Infof("Error unmarshalling message: %v", err)
				continue
			}
			if event.EventType == "CurrencyRate" {
				logrus.Infof("Processing event: %+v", event)

				temp := templates.ExchangeRateTemplate{
					CreatedAt:    event.Data.CreatedAt,
					ExchangeRate: event.Data.ExchangeRate,
				}
				err := emailSender.SendEmail(event.Data.Recipients, "Currency Rate", temp)
				if err != nil {
					logrus.Errorf("Error sending email: %v", err)
				} else {
					logrus.Infof("Successfully sent email to %v", event.Data.Recipients)
				}
			} else {
				logrus.Infof("Ignoring event of type: %s", event.EventType)
			}
		}
	}()

	logrus.Infof("Waiting for messages")
	<-forever
}
