package currencies

import (
	"github.com/skobelina/currency_converter/configs"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

type CurrencyHandler interface {
	SetNext(handler CurrencyHandler)
	Handle(config *configs.Config) (float64, error)
}

type BaseHandler struct {
	next CurrencyHandler
}

func (b *BaseHandler) SetNext(handler CurrencyHandler) {
	b.next = handler
}

func (b *BaseHandler) Handle(config *configs.Config) (float64, error) {
	if b.next != nil {
		return b.next.Handle(config)
	}
	return 0, serializer.NewBadRequestError("no handler could handle the request")
}
