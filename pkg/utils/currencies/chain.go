package currencies

import errors "github.com/skobelina/currency_converter/pkg/utils/errors"

type CurrencyHandler interface {
	SetNext(handler CurrencyHandler)
	Handle() (float64, error)
}

type BaseHandler struct {
	next CurrencyHandler
}

func (b *BaseHandler) SetNext(handler CurrencyHandler) {
	b.next = handler
}

func (b *BaseHandler) Handle() (float64, error) {
	if b.next != nil {
		return b.next.Handle()
	}
	return 0, errors.NewBadRequestError("no handler could handle the request")
}
