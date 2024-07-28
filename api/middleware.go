package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
)

func OptionsHandler() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				serializer.SetCorsHeaders(w)
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
