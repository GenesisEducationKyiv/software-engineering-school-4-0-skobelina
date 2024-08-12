package cronjobs

import (
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
	return r.db.Create(subscriber).Error
}

func (r *repository) FindByEmail(email string) (*Subscriber, error) {
	var subscriber Subscriber
	if err := r.db.Where("email = ?", email).First(&subscriber).Error; err != nil {
		return nil, err
	}
	return &subscriber, nil
}

func (r *repository) Search() ([]Subscriber, error) {
	var subscribers []Subscriber
	err := r.db.Find(&subscribers).Error
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}

func (r *repository) Delete(subscriber *Subscriber) error {
	return r.db.Delete(subscriber).Error
}

func (r *repository) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Subscriber{}).Error
}