package shared

import (
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Services                 []string      `env:"APP_SERVICES" envDefault:"db:postgres,messaging:nats"`
	Modules                  []string      `env:"APP_MODULES" envDefault:"ui,taskchecker,emailnotifier:null"`
	Addr                     string        `env:"APP_ADDR,notEmpty" envDefault:":8080"`
	LogLevel                 string        `env:"APP_LOG_LEVEL" envDefault:"warn"`
	PostgresConnectionString string        `env:"APP_POSTGRES_CONNECTION_STRING,notEmpty"`
	NATSURL                  string        `env:"APP_NATS_URL,notEmpty"`
	NATSToken                string        `env:"APP_NATS_TOKEN,notEmpty"`
	EmailToAddress           string        `env:"APP_EMAIL_TO_ADDRESS"`
	SMTPHost                 string        `env:"APP_SMTP_HOST"`
	SMTPPort                 int           `env:"APP_SMTP_PORT" envDefault:"587"`
	SMTPFromName             string        `env:"APP_SMTP_FROM_NAME"`
	SMTPFromAddress          string        `env:"APP_SMTP_FROM_ADDRESS"`
	SMTPPassword             string        `env:"APP_SMTP_PASSWORD"`
	TaskCheckInterval        time.Duration `env:"APP_TASK_CHECK_INTERVAL_SECONDS,notEmpty" envDefault:"60s"`
	TaskCheckExpiringWindow  time.Duration `env:"APP_TASK_CHECK_EXPIRING_WINDOW,notEmpty" envDefault:"24h"`
	TaskCheckDeleteWindow    time.Duration `env:"APP_TASK_CHECK_DELETE_WINDOW,notEmpty" envDefault:"48h"`
}

func (c *Config) Load() error {
	return env.Parse(c)
}
