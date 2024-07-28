package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/skobelina/email_sender/configs"
	cronjobs "github.com/skobelina/email_sender/internal/cron-jobs"
	"github.com/skobelina/email_sender/internal/mails"
	"github.com/skobelina/email_sender/pkg/metrics"
	"github.com/skobelina/email_sender/pkg/queue"
	"github.com/skobelina/email_sender/pkg/repo"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	config, err := configs.LoadConfig(".env")
	if err != nil {
		logrus.Fatal("Failed to load configuration: ", err)
	}
	logrus.Info("Configuration loaded successfully")

	repo, err := repo.Connect(config.DatabaseURL)
	if err != nil {
		logrus.Fatalf("Failed to connect to the database: %v", err)
	}
	logrus.Info("Connected to the database")
	if err := repo.AutoMigrate(&cronjobs.Subscriber{}); err != nil {
		logrus.Warnf("Failed to migrate database: %v", err)
	}

	rabbitMQ, err := queue.NewRabbitMQ(config.RabbitMQURL, "events")
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	subscriberRepo := cronjobs.NewRepository(repo)
	mailService := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost, config)
	cronJobService := cronjobs.NewCronJobService(mailService, rabbitMQ, subscriberRepo)
	cronJobService.InitializeSubscribers()
	logrus.Info("Initialized cron job service")

	metrics.Init()
	logrus.Info("Metrics initialized")

	logrus.Info("Starting email sender service")

	go func() {
		log.Println("Starting metrics server on :8081")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8081", nil); err != nil {
			logrus.Fatalf("Error starting metrics server: %v", err)
		}
	}()

	go cronJobService.ConsumeMessages()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logrus.Warn("Received shutdown signal")

	logrus.Info("Email sender service stopped")
}
