package internal

import "log/slog"

type NullEmailClient struct {
	Logger *slog.Logger
}

func (c *NullEmailClient) SendEmail(to string, subject string, templateName string, data any) error {
	c.Logger.Info("send email", "to", to, "subject", subject, "template_name", templateName)
	return nil
}
