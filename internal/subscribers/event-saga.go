package subscribers

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/pkg/queue"
)

type EventProcessor struct {
	repo     Repository
	rabbitMQ *queue.RabbitMQ
}

func NewEventProcessor(repo Repository, rabbitMQ *queue.RabbitMQ) *EventProcessor {
	return &EventProcessor{
		repo:     repo,
		rabbitMQ: rabbitMQ,
	}
}

func (ep *EventProcessor) ProcessSubscribeValidationRequestedEvent(event Event) error {
	if !isEmailValid(event.Data.Email) {
		ep.publishSubscribeValidationFailedEvent(event)
		return fmt.Errorf("validation failed for email: %s", event.Data.Email)
	}

	ep.publishSubscribeValidatedEvent(event)
	return nil
}

func (ep *EventProcessor) ProcessSubscribeValidatedEvent(event Event) error {
	subscriber := &Subscriber{
		Email: event.Data.Email,
	}
	if err := ep.repo.Create(subscriber); err != nil {
		ep.publishSubscribeFailedEvent(event)
		return err
	}
	ep.publishSubscribeCompletedEvent(event)
	return nil
}

func (ep *EventProcessor) publishSubscribeValidatedEvent(event Event) {
	event.EventType = "SubscribeValidated"
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error marshalling event: %v", err)
		return
	}
	if err := ep.rabbitMQ.PublishMessage(string(body)); err != nil {
		logrus.Errorf("Error publishing message: %v", err)
	}
}

func (ep *EventProcessor) publishSubscribeValidationFailedEvent(event Event) {
	event.EventType = "SubscribeValidationFailed"
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error marshalling event: %v", err)
		return
	}
	if err := ep.rabbitMQ.PublishMessage(string(body)); err != nil {
		logrus.Errorf("Error publishing message: %v", err)
	}
}

func (ep *EventProcessor) publishSubscribeCompletedEvent(event Event) {
	event.EventType = "SubscribeCompleted"
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error marshalling event: %v", err)
		return
	}
	if err := ep.rabbitMQ.PublishMessage(string(body)); err != nil {
		logrus.Errorf("Error publishing message: %v", err)
	}
}

func (ep *EventProcessor) ProcessSubscribeFailedEvent(event Event) {
	subscriber, err := ep.repo.FindByEmail(event.Data.Email)
	if err != nil {
		logrus.Errorf("Error finding subscriber: %v", err)
		return
	}

	if subscriber != nil {
		if err := ep.repo.Delete(subscriber); err != nil {
			logrus.Errorf("Error deleting subscriber: %v", err)
			return
		}
		logrus.Infof("Compensation action: Deleted subscriber %s", event.Data.Email)
	}
}

func (ep *EventProcessor) publishSubscribeFailedEvent(event Event) {
	event.EventType = "SubscribeFailed"
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error marshalling event: %v", err)
		return
	}
	if err := ep.rabbitMQ.PublishMessage(string(body)); err != nil {
		logrus.Errorf("Error publishing message: %v", err)
	}
}
