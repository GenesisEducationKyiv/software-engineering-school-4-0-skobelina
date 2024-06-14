package subscribers

import (
	"gorm.io/gorm"

	"github.com/skobelina/currency_converter/constants"
	"github.com/skobelina/currency_converter/domains"
	errors "github.com/skobelina/currency_converter/utils/errors"
)

type SubscriberService struct {
	repo *gorm.DB
}

func NewService(repo *gorm.DB) *SubscriberService {
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
	if err := s.repo.Create(subscriber).Error; err != nil {
		return nil, errors.NewInternalServerErrorf("cannot create email db row: %v", err)
	}
	status := constants.StatusAdded
	return &status, nil
}

func (s *SubscriberService) checkEmail(email string) (bool, error) {
	var count int64
	if err := s.repo.
		Model(&Subscriber{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, errors.NewInternalServerErrorf("failed to check email: %v", err)
	}
	return count > 0, nil
}

func (s *SubscriberService) Search(filter *SearchSubscribeRequest) (*SearchSubscribeResponse, error) {
	if filter == nil {
		filter = new(SearchSubscribeRequest)
		filter.Validate()
	}

	q := s.repo.
		Table("subscribers").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Order(filter.OrderString())

	var subscribers []Subscriber
	if err := q.Find(&subscribers).Error; err != nil {
		return nil, errors.NewInternalServerErrorf("cannot get subscribers: %v", err)
	}
	var count int64
	q = s.repo.
		Table("subscribers")
	if err := q.Count(&count).Error; err != nil {
		return nil, errors.NewInternalServerErrorf("cannot count subscribers: %v", err)
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
