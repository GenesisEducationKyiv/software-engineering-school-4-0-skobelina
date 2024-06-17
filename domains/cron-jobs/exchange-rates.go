package cronjobs

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/domains"
	"github.com/skobelina/currency_converter/domains/mails/templates"
	"github.com/skobelina/currency_converter/domains/subscribers"
)

func (s *CronJobService) NotificationExchangeRates() error {
	subscribersResp, err := s.Subscribers.Search(&subscribers.SearchSubscribeRequest{
		Filter: domains.DefaultFilter(),
	})
	if err != nil {
		return err
	}
	recipients := make([]string, 0, len(subscribersResp.Data))
	for _, person := range subscribersResp.Data {
		recipients = append(recipients, person.Email)
	}
	currentTime := time.Now().Format("2006-01-02")

	exchangeRate, err := s.Rates.Get()
	if err != nil {
		return err
	}
	template := templates.ExchangeRateTemplate{
		CreatedAt:    currentTime,
		ExchangeRate: strconv.FormatFloat(*exchangeRate, 'f', 2, 64),
	}
	if err := s.Mail.SendEmail(recipients, "Exchange rates notification", template); err != nil {
		logrus.Errorf("CronJob: NotificationExchangeRates: %v", err)
	}

	logrus.Infof("CronJob: NotificationExchangeRates: sent to %d subscribers", len(recipients))
	return nil
}
