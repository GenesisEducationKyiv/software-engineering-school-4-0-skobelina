package cronjobs

import (
	"net/http"

	"github.com/gorilla/mux"
	errors "github.com/skobelina/currency_converter/pkg/utils/errors"
	"github.com/skobelina/currency_converter/pkg/utils/rest"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

const cronKey = "E6B3C4F7"

type CronJobServiceInterface interface {
	Close() error
	NotificationExchangeRates() error
}

type handler struct {
	service CronJobServiceInterface
}

func NewHandler(s CronJobServiceInterface) *handler {
	return &handler{s}
}

func (h *handler) Register(r *mux.Router) {
	r.HandleFunc("/api/cron-jobs/notifications/exchange-rates/{key}", rest.ErrorHandler(h.notificationExchangeRates)).Methods("GET", "OPTIONS")
}

func (h *handler) notificationExchangeRates(w http.ResponseWriter, r *http.Request) error {
	if valid := validateCronjobRequest(r); !valid {
		return errors.NewForbiddenError()
	}
	err := h.service.NotificationExchangeRates()
	if err != nil {
		return err
	}
	return serializer.SendNoContent(w)
}

func validateCronjobRequest(r *http.Request) bool {
	vars := mux.Vars(r)
	return vars["key"] == cronKey
}
