package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EmailRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_requests_total",
			Help: "Total number of email requests",
		},
		[]string{"status"},
	)
	EmailRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "email_request_duration_seconds",
			Help:    "Histogram of response latency for email requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)
)

func Init() {
	prometheus.MustRegister(EmailRequestCount)
	prometheus.MustRegister(EmailRequestDuration)
}
