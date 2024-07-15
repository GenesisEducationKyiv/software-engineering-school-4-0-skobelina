package subscribers

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(subscriber *Subscriber) error
	FindByEmail(email string) (*Subscriber, error)
	Search(filter *SearchSubscribeRequest) ([]Subscriber, int64, error)
	Delete(subscriber *Subscriber) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(subscriber *Subscriber) error {
	return r.db.Create(subscriber).Error
}

func (r *repository) FindByEmail(email string) (*Subscriber, error) {
	var subscriber Subscriber
	if err := r.db.Where("email = ?", email).First(&subscriber).Error; err != nil {
		return nil, err
	}
	return &subscriber, nil
}

func (r *repository) Search(filter *SearchSubscribeRequest) ([]Subscriber, int64, error) {
	var subscribers []Subscriber
	q := r.db.Offset(filter.Offset).Limit(filter.Limit).Order(filter.OrderString())
	if err := q.Find(&subscribers).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := r.db.Model(&Subscriber{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return subscribers, count, nil
}

func (r *repository) Delete(subscriber *Subscriber) error {
	return r.db.Delete(subscriber).Error
}
