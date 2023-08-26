package internal

type EmailClient interface {
	SendEmail(to string, subject string, templateName string, data any) error
}
