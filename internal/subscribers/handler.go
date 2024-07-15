package subscribers

import (
	"net/http"

	"github.com/gorilla/mux"

	domains "github.com/skobelina/currency_converter/internal"
	errors "github.com/skobelina/currency_converter/pkg/utils/errors"
	"github.com/skobelina/currency_converter/pkg/utils/rest"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

type SubscriberServiceInterface interface {
	Create(request *SubscriberRequest) (*string, error)
	Search(filter *SearchSubscribeRequest) (*SearchSubscribeResponse, error)
	Delete(request *SubscriberRequest) (*string, error)
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
	r.HandleFunc("/api/subscribe", rest.ErrorHandler(h.delete)).Methods("DELETE", "OPTIONS")
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

// swagger:route DELETE /subscribe Subscription deleteSubscribe
// Unsubscribe from receiving current exchange rates
//
// responses:
//
//	200: body:message ok
//	404: notFound
func (h *handler) delete(w http.ResponseWriter, r *http.Request) error {
	request := new(SubscriberRequest)
	if err := serializer.ParseJsonBody(r.Body, request); err != nil {
		return err
	}
	status, err := h.service.Delete(request)
	if err != nil {
		return err
	}
	return serializer.SendJSON(w, http.StatusOK, status)
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
