package shared

import (
	"slices"
	"time"

	"github.com/caarlos0/env/v9"
)

type SharedConfig struct {
	Services                 []string `env:"APP_SHARED_SERVICES" envDefault:"db:postgres,attachments:nats,messaging:nats"`
	Modules                  []string `env:"APP_SHARED_MODULES" envDefault:"ui,taskchecker,emailnotifier:null"`
	LogLevel                 string   `env:"APP_SHARED_LOG_LEVEL" envDefault:"warn"`
	PostgresConnectionString string   `env:"APP_SHARED_POSTGRES_CONNECTION_STRING,notEmpty"`
	NATSURL                  string   `env:"APP_SHARED_NATS_URL,notEmpty"`
	NATSUser                 string   `env:"APP_SHARED_NATS_USER,notEmpty"`
	NATSPassword             string   `env:"APP_SHARED_NATS_PASSWORD,notEmpty"`
	AttachmentsPath          string   `env:"APP_SHARED_ATTACHMENTS_PATH" envDefault:"attachments"`
}

type UIConfig struct {
	Addr   string `env:"APP_UI_ADDR,notEmpty" envDefault:":8080"`
	HubURL string `env:"APP_UI_HUB_URL,notEmpty"`
}

type TaskCheckerConfig struct {
	CheckInterval  time.Duration `env:"APP_TASK_CHECKER_CHECK_INTERVAL,notEmpty" envDefault:"60s"`
	ExpiringWindow time.Duration `env:"APP_TASK_CHECKER_EXPIRING_WINDOW,notEmpty" envDefault:"24h"`
	DeleteWindow   time.Duration `env:"APP_TASK_CHECKER_DELETE_WINDOW,notEmpty" envDefault:"48h"`
}

type EmailNotifierConfig struct {
	ToAddress       string `env:"APP_EMAIL_NOTIFIER_TO_ADDRESS"`
	SMTPHost        string `env:"APP_EMAIL_NOTIFIER_SMTP_HOST"`
	SMTPPort        int    `env:"APP_EMAIL_NOTIFIER_SMTP_PORT" envDefault:"587"`
	SMTPFromName    string `env:"APP_EMAIL_NOTIFIER_SMTP_FROM_NAME"`
	SMTPFromAddress string `env:"APP_EMAIL_NOTIFIER_SMTP_FROM_ADDRESS"`
	SMTPPassword    string `env:"APP_EMAIL_NOTIFIER_SMTP_PASSWORD"`
}

type Config struct {
	Shared        SharedConfig
	UI            UIConfig
	TaskChecker   TaskCheckerConfig
	EmailNotifier EmailNotifierConfig
}

func (c *Config) Load() error {
	return env.Parse(c)
}

func (c *Config) IsServiceEnabled(name string) bool {
	return slices.Contains(c.Shared.Services, name)
}

func (c *Config) IsModuleEnabled(name string) bool {
	return slices.Contains(c.Shared.Modules, name)
}
