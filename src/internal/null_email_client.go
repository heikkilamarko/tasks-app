package internal

import (
	"context"
	"log/slog"
)

type NullEmailClient struct {
	Logger *slog.Logger
}

func (c *NullEmailClient) SendEmail(ctx context.Context, to string, subject string, templateName string, data any) error {
	c.Logger.Info("send email", "to", to, "subject", subject, "template_name", templateName)
	return nil
}
