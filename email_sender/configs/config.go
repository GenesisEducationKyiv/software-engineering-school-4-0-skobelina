package configs

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" required:"true"`
	MailPass    string `env:"MAILPASS" required:"true"`
	RabbitMQURL string `env:"RABBITMQ_URL" required:"true"`
}
