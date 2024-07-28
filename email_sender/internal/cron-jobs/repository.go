package cronjobs

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository interface {
	Create(subscriber *Subscriber) error
	FindByEmail(email string) (*Subscriber, error)
	Search() ([]Subscriber, error)
	Delete(subscriber *Subscriber) error
	DeleteAll() error
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

func (r *repository) Search() ([]Subscriber, error) {
	var subscribers []Subscriber
	err := r.db.Find(&subscribers).Error
	if err != nil {
		logrus.Errorf("Repository - Error searching subscribers: %v", err)
		return nil, err
	}
	logrus.Infof("Repository - Found %d subscribers", len(subscribers))
	return subscribers, nil
}

func (r *repository) Delete(subscriber *Subscriber) error {
	logrus.Infof("Repository - Deleting subscriber: %s", subscriber.Email)
	return r.db.Delete(subscriber).Error
}

func (r *repository) DeleteAll() error {
	logrus.Info("Repository - Deleting all subscribers")
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Subscriber{}).Error
}
