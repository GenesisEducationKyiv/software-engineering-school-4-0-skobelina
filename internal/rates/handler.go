package rates

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/pkg/utils/rest"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

type RateServiceInterface interface {
	Get() (*float64, error)
}

type handler struct {
	service RateServiceInterface
}

func NewHandler(s RateServiceInterface) *handler {
	return &handler{s}
}

func (h *handler) Register(r *mux.Router) {
	r.HandleFunc("/api/rate", rest.ErrorHandler(h.get)).Methods("GET", "OPTIONS")
}

// swagger:route GET /rate Rate getRate
// Get the current USD to UAH rate
//
// responses:
//
//	200: body:rateResponse ok
//	400: statusBadRequest
func (h *handler) get(w http.ResponseWriter, r *http.Request) error {
	response, err := h.service.Get()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warnf("Handler - Rate not found")
			return serializer.NewItemNotFoundError("rate not found")
		} else if err.Error() == "bad request" {
			logrus.Warnf("Handler - Bad request error")
			return serializer.NewBadRequestError("bad request")
		} else if err.Error() == "internal server error" {
			logrus.Errorf("Handler - Internal server error")
			return serializer.NewInternalServerError("internal server error")
		} else {
			logrus.Errorf("Handler - Unexpected error: %v", err)
			return err
		}
	}
	logrus.Infof("Handler - Successfully fetched rate")
	return serializer.SendJSON(w, http.StatusOK, response)
}
