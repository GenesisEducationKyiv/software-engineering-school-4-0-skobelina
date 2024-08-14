package cronjobs

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/internal/constants"
)

type Event struct {
	EventID     string    `json:"eventId"`
	EventType   string    `json:"eventType"`
	AggregateID string    `json:"aggregateId"`
	Timestamp   string    `json:"timestamp"`
	Data        EventData `json:"data"`
}

type EventData struct {
	CreatedAt    string `json:"createdAt"`
	ExchangeRate string `json:"exchangeRate"`
}

func (s *CronJobService) NotificationExchangeRates() error {
	currentTime := time.Now().Format(constants.DateFormatYMD)
	exchangeRate, err := s.rates.Get()
	if err != nil {
		logrus.Errorf("CronJobService - Error getting exchange rate: %v", err)
		return err
	}
	event := Event{
		EventID:     uuid.New().String(),
		EventType:   "CurrencyRate",
		AggregateID: "rate-" + uuid.New().String(),
		Timestamp:   time.Now().Format(time.RFC3339),
		Data: EventData{
			CreatedAt:    currentTime,
			ExchangeRate: strconv.FormatFloat(*exchangeRate, 'f', 2, 64),
		},
	}
	body, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("CronJobService - Error marshalling event: %v", err)
		return err
	}
	if err := s.rabbitMQ.PublishMessage(string(body)); err != nil {
		logrus.Errorf("CronJobService - Error publishing message: %v", err)
		return err
	}

	logrus.Infof("CronJobService: NotificationExchangeRates: sent exchange rate event")
	return nil
}
