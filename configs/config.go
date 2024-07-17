package configs

type Config struct {
	DatabaseURL            string `env:"DATABASE_URL" required:"true"`
	AppCurrencyExchangeURL string `env:"APP_CURRENCY_EXCHANGE_URL" required:"true"`
	AppCurrencyExchangeKey string `env:"APP_CURRENCY_EXCHANGE_KEY" required:"true"`
	AppCurrencyBeaconURL   string `env:"APP_CURRENCY_BEACON_URL" required:"true"`
	AppCurrencyBeaconKey   string `env:"APP_CURRENCY_BEACON_KEY" required:"true"`
	MailPass               string `env:"MAILPASS" required:"true"`
	RabbitMQURL            string `env:"RABBITMQ_URL" required:"true"`
}
