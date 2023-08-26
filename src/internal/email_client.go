package internal

import "context"

type EmailClient interface {
	SendEmail(ctx context.Context, to string, subject string, templateName string, data any) error
}
