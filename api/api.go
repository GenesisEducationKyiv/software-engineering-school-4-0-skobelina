package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	cronJobs "github.com/skobelina/currency_converter/internal/cron-jobs"
	"github.com/skobelina/currency_converter/internal/rates"
	"github.com/skobelina/currency_converter/internal/subscribers"
	"github.com/skobelina/currency_converter/pkg/metrics"
)

type Api interface {
	Handle() error
}

type api struct {
	router *mux.Router
}

func New() Api {
	deps := registerDependencies()
	r := mux.NewRouter()
	rates.NewHandler(deps.Rates).Register(r)
	subscribers.NewHandler(deps.Subscribers).Register(r)
	cronJobService := cronJobs.NewService(deps.Rates, deps.RabbitMQ)
	cronJobs.NewHandler(cronJobService).Register(r)

	r.Use(
		OptionsHandler(),
		LoggingMiddleware(),
	)
	return &api{
		router: r,
	}
}

func (a *api) Handle() error {
	metrics.Init()

	http.Handle("/", a.router)
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
	}
	logrus.Info("Starting API service")
	return srv.ListenAndServe()
}
