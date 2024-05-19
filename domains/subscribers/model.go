package subscribers

import (
	"strings"

	"github.com/skobelina/currency_converter/domains"
	errors "github.com/skobelina/currency_converter/utils/errors"
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
		return errors.NewBadRequestError("validation errors occurent: email must be set")
	}
	if strings.ContainsAny(s.Email, " \t\n") {
		return errors.NewBadRequestError("validation errors occurent: email cannot contain spaces")
	}
	if !isEmailValid(s.Email) {
		return errors.NewBadRequestError("validation errors occurent: email format is invalid")
	}
	return nil
}

func isEmailValid(email string) bool {
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at
}
