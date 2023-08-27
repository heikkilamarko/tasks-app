package internal

import "github.com/caarlos0/env/v9"

type Config struct {
	Addr                     string `env:"APP_ADDR,notEmpty" envDefault:":8080"`
	LogLevel                 string `env:"APP_LOG_LEVEL" envDefault:"warn"`
	PostgresConnectionString string `env:"APP_POSTGRES_CONNECTION_STRING,notEmpty"`
	SMTPHost                 string `env:"APP_SMTP_HOST"`
	SMTPPort                 int    `env:"APP_SMTP_PORT" envDefault:"587"`
	SMTPFromName             string `env:"APP_SMTP_FROM_NAME"`
	SMTPFromAddress          string `env:"APP_SMTP_FROM_ADDRESS"`
	SMTPPassword             string `env:"APP_SMTP_PASSWORD"`
}

func (c *Config) Load() error {
	return env.Parse(c)
}
