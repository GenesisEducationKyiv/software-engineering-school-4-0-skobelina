package subscribers

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	domains "github.com/skobelina/currency_converter/internal"
	"github.com/skobelina/currency_converter/internal/constants"
	"github.com/skobelina/currency_converter/pkg/queue"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
	"gorm.io/gorm"
)

type SubscriberService struct {
	repo     Repository
	rabbitMQ *queue.RabbitMQ
}

func NewService(repo Repository, rabbitMQ *queue.RabbitMQ) *SubscriberService {
	return &SubscriberService{
		repo:     repo,
		rabbitMQ: rabbitMQ,
	}
}

func (s *SubscriberService) Create(request *SubscriberRequest) (*string, error) {
	subscriber := request.Map()
	if err := subscriber.Validate(); err != nil {
		return nil, err
	}
	if exists, err := s.checkEmail(subscriber.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, serializer.NewIsConflictError("email already exists")
	}
	if err := s.repo.Create(subscriber); err != nil {
		return nil, serializer.NewInternalServerErrorf("cannot create email db row: %v", err)
	}
	status := constants.StatusAdded
	if err := s.createSubscribeEvent(request.Email); err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *SubscriberService) checkEmail(email string) (bool, error) {
	subscriber, err := s.repo.FindByEmail(email)
	if err != nil {
		if serializer.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, serializer.NewInternalServerErrorf("failed to check email: %v", err)
	}
	return subscriber != nil, nil
}

func (s *SubscriberService) createSubscribeEvent(email string) error {
	event := Event{
		EventID:     uuid.New().String(),
		EventType:   "Subscribe",
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
	logrus.Infof("SubscriberService: Subscribe: sent subscribe event for %s", email)
	return nil
}

func (s *SubscriberService) Search(filter *SearchSubscribeRequest) (*SearchSubscribeResponse, error) {
	if filter == nil {
		filter = new(SearchSubscribeRequest)
		filter.Validate()
	}

	subscribers, count, err := s.repo.Search(filter)
	if err != nil {
		return nil, serializer.NewInternalServerErrorf("cannot get subscribers: %v", err)
	}

	return &SearchSubscribeResponse{
		Data: subscribers,
		Pagination: &domains.Pagination{
			Order:      filter.OrderString(),
			Offset:     filter.Offset,
			Limit:      filter.Limit,
			TotalItems: &count,
		},
	}, nil
}

func (s *SubscriberService) Delete(request *SubscriberRequest) (*string, error) {
	subscriber, err := s.repo.FindByEmail(request.Email)
	if err != nil {
		if serializer.Is(err, gorm.ErrRecordNotFound) {
			return nil, serializer.NewItemNotFoundErrorf("subscriber not found")
		}
		return nil, serializer.NewInternalServerErrorf("failed to find subscriber: %v", err)
	}

	if err := s.repo.Delete(subscriber); err != nil {
		return nil, serializer.NewInternalServerErrorf("failed to delete subscriber: %v", err)
	}

	status := constants.StatusDeleted
	if err := s.createUnsubscribeEvent(request.Email); err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *SubscriberService) createUnsubscribeEvent(email string) error {
	event := Event{
		EventID:     uuid.New().String(),
		EventType:   "Unsubscribe",
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
	logrus.Infof("SubscriberService: Unsubscribe: sent unsubscribe event for %s", email)
	return nil
}
