package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	MailPass    string `envconfig:"MAILPASS" required:"true"`
	RabbitMQURL string `envconfig:"RABBITMQ_URL" required:"true"`
}

func LoadConfig(envFile string) (*Config, error) {
	err := godotenv.Load(envFile)
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	var config Config
	err = envconfig.Process("email_sender", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
