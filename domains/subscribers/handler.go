package subscribers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/skobelina/currency_converter/domains"
	errors "github.com/skobelina/currency_converter/utils/errors"
	"github.com/skobelina/currency_converter/utils/rest"
	"github.com/skobelina/currency_converter/utils/serializer"
)

type SubscriberServiceInterface interface {
	Create(request *SubscriberRequest) (*string, error)
	Search(filter *SearchSubscribeRequest) (*SearchSubscribeResponse, error)
}

type handler struct {
	service SubscriberServiceInterface
}

func NewHandler(s SubscriberServiceInterface) *handler {
	return &handler{s}
}

func (h *handler) Register(r *mux.Router) {
	r.HandleFunc("/api/subscribe", rest.ErrorHandler(h.create)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/subscribe", rest.ErrorHandler(h.search)).Methods("GET", "OPTIONS")
}

// swagger:route POST /subscribe Subscription createSubscribe
// Sign up to receive the current exchange rates
//
// responses:
//
//	200: body:message ok
//	409: statusConflict
func (h *handler) create(w http.ResponseWriter, r *http.Request) error {
	request := new(SubscriberRequest)
	if err := serializer.ParseJsonBody(r.Body, request); err != nil {
		return err
	}
	status, err := h.service.Create(request)
	if err != nil {
		return err
	}
	return serializer.SendJSON(w, http.StatusOK, status)
}

// swagger:route GET /subscribe Subscription searchSubscribe
// Search all subscribers
//
// responses:
//
//	200: body:searchSubscribeResponse ok
func (h *handler) search(w http.ResponseWriter, r *http.Request) error {
	filter, err := getFilterFromQuery(r)
	if err != nil {
		return err
	}
	response, err := h.service.Search(filter)
	if err != nil {
		return err
	}
	return serializer.SendJSON(w, http.StatusOK, response)
}

func getFilterFromQuery(r *http.Request) (*SearchSubscribeRequest, error) {
	filter, err := domains.GetFilterFromQuery(r)
	if err != nil {
		return nil, errors.NewBadRequestError(err)
	}
	return &SearchSubscribeRequest{
		Filter: *filter,
	}, nil
}
