package rates

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	utils "github.com/skobelina/currency_converter/utils/errors"
	"github.com/skobelina/currency_converter/utils/rest"
	"github.com/skobelina/currency_converter/utils/serializer"
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
			return utils.NewItemNotFoundError("rate not found")
		} else if err.Error() == "bad request" {
			return utils.NewBadRequestError("bad request")
		} else if err.Error() == "internal server error" {
			return utils.NewInternalServerError("internal server error")
		} else {
			return err
		}
	}
	return serializer.SendJSON(w, http.StatusOK, response)
}
