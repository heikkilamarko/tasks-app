package emailnotifier

import (
	"context"
	"log/slog"
)

type NullEmailClient struct {
	Logger *slog.Logger
}

var _ EmailClient = (*NullEmailClient)(nil)

func (c *NullEmailClient) SendEmail(ctx context.Context, to string, subject string, templateName string, data any) error {
	c.Logger.Info("send email",
		slog.Group("email",
			slog.String("to", to),
			slog.String("subject", subject),
			slog.String("template_name", templateName),
		),
	)

	return nil
}
