package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/pkg/metrics"
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

func LoggingMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			logrus.Infof("Received request: %s %s", r.Method, r.URL.Path)
			metrics.RequestCount.WithLabelValues(r.Method, r.URL.Path).Inc()
			logrus.Infof("Incremented request count metric for: %s %s", r.Method, r.URL.Path)
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
			logrus.Infof("Observed request duration for: %s %s", r.Method, r.URL.Path)
		})
	}
}
