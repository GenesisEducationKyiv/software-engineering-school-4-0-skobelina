package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseURL            string `envconfig:"DATABASE_URL" required:"true"`
	AppCurrencyExchangeURL string `envconfig:"APP_CURRENCY_EXCHANGE_URL" required:"true"`
	AppCurrencyExchangeKey string `envconfig:"APP_CURRENCY_EXCHANGE_KEY" required:"true"`
	AppCurrencyBeaconURL   string `envconfig:"APP_CURRENCY_BEACON_URL" required:"true"`
	AppCurrencyBeaconKey   string `envconfig:"APP_CURRENCY_BEACON_KEY" required:"true"`
	MailPass               string `envconfig:"MAILPASS" required:"true"`
	RabbitMQURL            string `envconfig:"RABBITMQ_URL" required:"true"`
}

func LoadConfig(envFile string) (*Config, error) {
	err := godotenv.Load(envFile)
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	var config Config
	err = envconfig.Process("currency_converter", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
