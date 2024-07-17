package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/skobelina/email_sender/configs"
	cronjobs "github.com/skobelina/email_sender/internal/cron-jobs"
	"github.com/skobelina/email_sender/internal/mails"
	"github.com/skobelina/email_sender/pkg/queue"
	"github.com/skobelina/email_sender/pkg/repo"
)

func main() {
	err := godotenv.Load("email_sender/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config configs.Config
	err = envconfig.Process("email_sender", &config)
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
	mailService := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost)
	cronJobService := cronjobs.NewCronJobService(mailService, rabbitMQ, subscriberRepo)
	cronJobService.InitializeSubscribers()
	logrus.Info("Starting email sender service")

	go cronJobService.ConsumeMessages()
	select {}
}
