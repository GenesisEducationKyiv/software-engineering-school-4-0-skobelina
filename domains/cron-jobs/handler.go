package cronjobs

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/skobelina/currency_converter/domains/mails"
	errors "github.com/skobelina/currency_converter/utils/errors"
	"github.com/skobelina/currency_converter/utils/rest"
	"github.com/skobelina/currency_converter/utils/serializer"
)

const cronKey = "E6B3C4F7"

type CronJobServiceInterface interface {
	Close() error
	NotificationExchangeRates() error
}

type handler struct {
	config  *CronJobConfig
	service CronJobServiceInterface
}

type CronJobConfig struct {
	DatabaseURL   string
	CurrencyRates float64
	MailService   *mails.MailService
}

func NewHandler(config *CronJobConfig) rest.Registrable {
	service := NewService(config)
	return &handler{config: config, service: service}
}

func (h *handler) Register(r *mux.Router) {
	r.HandleFunc("/api/cron-jobs/notifications/exchange-rates/{key}", rest.ErrorHandler(h.notificationExchangeRates)).Methods("GET", "OPTIONS")
}

func (h *handler) notificationExchangeRates(w http.ResponseWriter, r *http.Request) error {
	if valid := validateCronjobRequest(r); !valid {
		return errors.NewForbiddenError()
	}
	service := NewService(h.config)
	defer service.Close()
	err := service.NotificationExchangeRates()
	if err != nil {
		return err
	}
	return serializer.SendNoContent(w)
}

func validateCronjobRequest(r *http.Request) bool {
	vars := mux.Vars(r)
	return vars["key"] == cronKey
}
