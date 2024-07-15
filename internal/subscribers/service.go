package subscribers

import (
	domains "github.com/skobelina/currency_converter/internal"
	"github.com/skobelina/currency_converter/internal/constants"
	errors "github.com/skobelina/currency_converter/pkg/utils/errors"
	"gorm.io/gorm"
)

type SubscriberService struct {
	repo Repository
}

func NewService(repo Repository) *SubscriberService {
	return &SubscriberService{repo}
}

func (s *SubscriberService) Create(request *SubscriberRequest) (*string, error) {
	subscriber := request.Map()
	if err := subscriber.Validate(); err != nil {
		return nil, err
	}
	if exists, err := s.checkEmail(subscriber.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.NewIsConflictError("email already exists")
	}
	if err := s.repo.Create(subscriber); err != nil {
		return nil, errors.NewInternalServerErrorf("cannot create email db row: %v", err)
	}
	status := constants.StatusAdded
	return &status, nil
}

func (s *SubscriberService) checkEmail(email string) (bool, error) {
	subscriber, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.NewInternalServerErrorf("failed to check email: %v", err)
	}
	return subscriber != nil, nil
}

func (s *SubscriberService) Search(filter *SearchSubscribeRequest) (*SearchSubscribeResponse, error) {
	if filter == nil {
		filter = new(SearchSubscribeRequest)
		filter.Validate()
	}

	subscribers, count, err := s.repo.Search(filter)
	if err != nil {
		return nil, errors.NewInternalServerErrorf("cannot get subscribers: %v", err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewItemNotFoundErrorf("subscriber not found")
		}
		return nil, errors.NewInternalServerErrorf("failed to find subscriber: %v", err)
	}

	if err := s.repo.Delete(subscriber); err != nil {
		return nil, errors.NewInternalServerErrorf("failed to delete subscriber: %v", err)
	}

	status := constants.StatusDeleted
	return &status, nil
}
