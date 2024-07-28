package rest

import (
	"net/http"

	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

func ErrorHandler(h func(w http.ResponseWriter, r *http.Request) error, middlewares ...func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, middleware := range middlewares {
			if err := middleware(w, r); err != nil {
				_ = serializer.SendError(w, err)
				return
			}
		}
		if err := h(w, r); err != nil {
			_ = serializer.SendError(w, err)
		}
	})
}
