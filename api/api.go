package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	cronJobs "github.com/skobelina/currency_converter/domains/cron-jobs"
	"github.com/skobelina/currency_converter/domains/rates"
	"github.com/skobelina/currency_converter/domains/subscribers"
)

var (
	databaseURL = os.Getenv("DATABASE_URL")
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
	cronJobService := cronJobs.NewService(deps.Repo, deps.MailService, deps.Rates, deps.Subscribers)
	cronJobs.NewHandler(cronJobService).Register(r)

	r.Use(
		OptionsHandler(),
	)
	return &api{
		router: r,
	}
}

func (a *api) Handle() error {
	http.Handle("/", a.router)
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
