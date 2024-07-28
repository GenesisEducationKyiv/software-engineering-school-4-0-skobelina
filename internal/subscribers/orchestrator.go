package subscribers

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/pkg/queue"
)

type Saga struct {
	rabbitMQ *queue.RabbitMQ
}

func NewSaga(rabbitMQ *queue.RabbitMQ) *Saga {
	return &Saga{rabbitMQ: rabbitMQ}
}

func (s *Saga) StartSubscribeSaga(email string) error {
	event := Event{
		EventID:     uuid.New().String(),
		EventType:   "SubscribeValidationRequested",
		AggregateID: "subscriber-" + uuid.New().String(),
		Timestamp:   time.Now().Format(time.RFC3339),
		Data: EventData{
			Email: email,
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
	logrus.Infof("Saga: StartSubscribeSaga: sent subscribe validation requested event for %s", email)
	return nil
}
