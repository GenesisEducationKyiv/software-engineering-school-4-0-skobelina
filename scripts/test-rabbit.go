package main

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Errorf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"events",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Errorf("Failed to declare a queue: %v", err)
	}

	event := Event{
		EventID:     "1",
		EventType:   "CurrencyRate",
		AggregateID: "rate-1",
		Timestamp:   time.Now().Format(time.RFC3339),
		Data: EventData{
			CreatedAt:    time.Now().Format("2006-01-02"),
			ExchangeRate: "5.50",
			Recipients:   []string{"ms.skobelina@gmail.com"},
		},
	}
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error marshalling event: %v", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		logrus.Errorf("Failed to publish a message: %v", err)
	}
	logrus.Infof("Sent %s", body)
}
