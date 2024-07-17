package subscribers

import (
	"strings"

	domains "github.com/skobelina/currency_converter/internal"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

type Subscriber struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type SubscriberRequest struct {
	Email string `json:"email"`
}

func (s *SubscriberRequest) Map() *Subscriber {
	return &Subscriber{
		Email: s.Email,
	}
}

// swagger:parameters createSubscribe
type CreateSubscribe struct {
	// in: body
	// required: true
	Body SubscriberRequest
}

// swagger:parameters deleteSubscribe
type DeleteSubscribe struct {
	// in: body
	// required: true
	Body SubscriberRequest
}

// swagger:parameters searchSubscribe
type SearchSubscribeRequest struct {
	domains.Filter
}

// swagger:model searchSubscribeResponse
type SearchSubscribeResponse struct {
	Data       []Subscriber        `json:"data"`
	Pagination *domains.Pagination `json:"pagination,omitempty"`
}

func (s *Subscriber) Validate() error {
	if strings.TrimSpace(s.Email) == "" {
		return serializer.NewBadRequestError("validation errors occurent: email must be set")
	}
	if strings.ContainsAny(s.Email, " \t\n") {
		return serializer.NewBadRequestError("validation errors occurent: email cannot contain spaces")
	}
	if !isEmailValid(s.Email) {
		return serializer.NewBadRequestError("validation errors occurent: email format is invalid")
	}
	return nil
}

func isEmailValid(email string) bool {
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at
}

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
	Email        string `json:"email"`
}
