package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/email_sender/configs"
	cronjobs "github.com/skobelina/email_sender/internal/cron-jobs"
	"github.com/skobelina/email_sender/internal/mails"
	"github.com/skobelina/email_sender/pkg/queue"
	"github.com/skobelina/email_sender/pkg/repo"
)

func main() {
	config, err := configs.LoadConfig(".env")
	if err != nil {
		log.Fatal(err.Error())
	}

	repo, err := repo.Connect(config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	if err := repo.AutoMigrate(&cronjobs.Subscriber{}); err != nil {
		logrus.Infof("failed to migrate database: %v", err)
	}
	rabbitMQ, err := queue.NewRabbitMQ(config.RabbitMQURL, "events")
	if err != nil {
		logrus.Infof("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	subscriberRepo := cronjobs.NewRepository(repo)
	mailService := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost, config)
	cronJobService := cronjobs.NewCronJobService(mailService, rabbitMQ, subscriberRepo)
	cronJobService.InitializeSubscribers()
	logrus.Info("Starting email sender service")

	go cronJobService.ConsumeMessages()
	select {}
}
