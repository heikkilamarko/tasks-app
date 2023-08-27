package internal

import "github.com/caarlos0/env/v9"

type Config struct {
	Address            string `env:"APP_ADDRESS,notEmpty" envDefault:":8080"`
	DBConnectionString string `env:"APP_DB_CONNECTION_STRING,notEmpty"`
	LogLevel           string `env:"APP_LOG_LEVEL" envDefault:"warn"`
	SMTPHost           string `env:"SMTP_HOST"`
	SMTPPort           int    `env:"SMTP_PORT" envDefault:"587"`
	SMTPFromName       string `env:"SMTP_FROM_NAME"`
	SMTPFromAddress    string `env:"SMTP_FROM_ADDRESS"`
	SMTPPassword       string `env:"SMTP_PASSWORD"`
}

func (c *Config) Load() error {
	return env.Parse(c)
}
