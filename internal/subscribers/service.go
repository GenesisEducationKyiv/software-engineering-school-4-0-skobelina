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
	saga     *Saga
}

func NewService(repo Repository, rabbitMQ *queue.RabbitMQ, saga *Saga) *SubscriberService {
	return &SubscriberService{
		repo:     repo,
		rabbitMQ: rabbitMQ,
		saga:     saga,
	}
}

func (s *SubscriberService) Create(request *SubscriberRequest) (*string, error) {
	subscriber := request.Map()
	if err := subscriber.Validate(); err != nil {
		logrus.Warnf("SubscriberService - Validation error: %v", err)
		return nil, err
	}
	if exists, err := s.checkEmail(subscriber.Email); err != nil {
		logrus.Errorf("SubscriberService - Error checking email: %v", err)
		return nil, err
	} else if exists {
		logrus.Warnf("SubscriberService - Email already exists: %s", subscriber.Email)
		return nil, serializer.NewIsConflictError("email already exists")
	}
	if err := s.saga.StartSubscribeSaga(request.Email); err != nil {
		logrus.Errorf("SubscriberService - Error starting saga: %v", err)
		return nil, serializer.NewInternalServerErrorf("cannot start saga: %v", err)
	}
	status := constants.StatusAdded
	if err := s.createSubscribeEvent(request.Email); err != nil {
		logrus.Errorf("SubscriberService - Error creating subscribe event: %v", err)
		return nil, err
	}
	logrus.Infof("SubscriberService - Subscriber created successfully: %s", subscriber.Email)
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
		logrus.Errorf("SubscriberService - Error searching subscribers: %v", err)
		return nil, serializer.NewInternalServerErrorf("cannot get subscribers: %v", err)
	}

	logrus.Infof("SubscriberService - Found %d subscribers", count)
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
			logrus.Warnf("SubscriberService - Subscriber not found: %s", request.Email)
			return nil, serializer.NewItemNotFoundErrorf("subscriber not found")
		}
		logrus.Errorf("SubscriberService - Error finding subscriber: %v", err)
		return nil, serializer.NewInternalServerErrorf("failed to find subscriber: %v", err)
	}

	if err := s.repo.Delete(subscriber); err != nil {
		logrus.Errorf("SubscriberService - Error deleting subscriber: %v", err)
		return nil, serializer.NewInternalServerErrorf("failed to delete subscriber: %v", err)
	}

	status := constants.StatusDeleted
	if err := s.createUnsubscribeEvent(request.Email); err != nil {
		logrus.Errorf("SubscriberService - Error creating unsubscribe event: %v", err)
		return nil, err
	}
	logrus.Infof("SubscriberService - Subscriber deleted successfully: %s", request.Email)
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
