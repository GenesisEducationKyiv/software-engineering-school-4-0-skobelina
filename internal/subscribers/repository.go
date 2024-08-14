package subscribers

import (
	"github.com/sirupsen/logrus"
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
	logrus.Infof("Repository - Creating subscriber: %s", subscriber.Email)
	return r.db.Create(subscriber).Error
}

func (r *repository) FindByEmail(email string) (*Subscriber, error) {
	var subscriber Subscriber
	err := r.db.Where("email = ?", email).First(&subscriber).Error
	if err != nil {
		logrus.Warnf("Repository - Subscriber not found: %s", email)
		return nil, err
	}
	logrus.Infof("Repository - Subscriber found: %s", email)
	return &subscriber, nil
}

func (r *repository) Search(filter *SearchSubscribeRequest) ([]Subscriber, int64, error) {
	var subscribers []Subscriber
	q := r.db.Offset(filter.Offset).Limit(filter.Limit).Order(filter.OrderString())
	err := q.Find(&subscribers).Error
	if err != nil {
		logrus.Errorf("Repository - Error searching subscribers: %v", err)
		return nil, 0, err
	}

	var count int64
	err = r.db.Model(&Subscriber{}).Count(&count).Error
	if err != nil {
		logrus.Errorf("Repository - Error counting subscribers: %v", err)
		return nil, 0, err
	}
	logrus.Infof("Repository - Found %d subscribers", count)
	return subscribers, count, nil
}

func (r *repository) Delete(subscriber *Subscriber) error {
	logrus.Infof("Repository - Deleting subscriber: %s", subscriber.Email)
	return r.db.Delete(subscriber).Error
}
