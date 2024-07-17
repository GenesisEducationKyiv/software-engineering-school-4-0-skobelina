package main

import (
	"os"

	"github.com/sirupsen/logrus"
	cronjobs "github.com/skobelina/email_sender/internal/cron-jobs"
	"github.com/skobelina/email_sender/internal/mails"
	"github.com/skobelina/email_sender/pkg/queue"
	"github.com/skobelina/email_sender/pkg/repo"
)

var (
	databaseURL = os.Getenv("DATABASE_URL")
	rabbitMQURL = os.Getenv("RABBITMQ_URL")
)

func main() {
	repo, err := repo.Connect(databaseURL)
	if err != nil {
		panic(err)
	}
	if err := repo.AutoMigrate(&cronjobs.Subscriber{}); err != nil {
		logrus.Infof("failed to migrate database: %v", err)
	}
	rabbitMQ, err := queue.NewRabbitMQ(rabbitMQURL, "events")
	if err != nil {
		logrus.Infof("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	subscriberRepo := cronjobs.NewRepository(repo)
	mailService := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost)
	cronJobService := cronjobs.NewCronJobService(mailService, rabbitMQ, subscriberRepo)
	cronJobService.InitializeSubscribers()
	logrus.Info("Starting email sender service")

	go cronJobService.ConsumeMessages()
	select {}
}
