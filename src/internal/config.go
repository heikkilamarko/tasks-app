package internal

import (
	"fmt"
	"strconv"
)

type Config struct {
	Address            string
	DBConnectionString string
	LogLevel           string
	SMTPHost           string
	SMTPPort           int
	SMTPFromName       string
	SMTPFromAddress    string
	SMTPPassword       string
}

func (c *Config) Load() error {
	var err error

	c.Address = Env("APP_ADDRESS", ":8080")
	c.DBConnectionString = Env("APP_DB_CONNECTION_STRING", "")
	c.LogLevel = Env("APP_LOG_LEVEL", "warn")
	c.SMTPHost = Env("SMTP_HOST", "")
	c.SMTPPort, err = strconv.Atoi(Env("SMTP_PORT", "587"))
	if err != nil {
		return fmt.Errorf("read SMTP_PORT env: %w", err)
	}
	c.SMTPFromName = Env("SMTP_FROM_NAME", "")
	c.SMTPFromAddress = Env("SMTP_FROM_ADDRESS", "")
	c.SMTPPassword = Env("SMTP_PASSWORD", "")

	return nil
}
