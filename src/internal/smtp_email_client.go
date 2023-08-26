package internal

import (
	"strings"

	"github.com/wneessen/go-mail"
)

type SMTPEmailClientOptions struct {
	Host        string
	Port        int
	FromName    string
	FromAddress string
	Password    string
}

type SMTPEmailClient struct {
	Options SMTPEmailClientOptions
}

func (c *SMTPEmailClient) SendEmail(to string, subject string, templateName string, data any) error {
	var bodyBuilder strings.Builder
	if err := EmailTemplates.ExecuteTemplate(&bodyBuilder, templateName, data); err != nil {
		return err
	}
	body := bodyBuilder.String()

	msg := mail.NewMsg()

	if err := msg.FromFormat(c.Options.FromName, c.Options.FromAddress); err != nil {
		return err
	}

	if err := msg.To(to); err != nil {
		return err
	}

	msg.Subject(subject)

	msg.SetBodyString(mail.TypeTextHTML, body)

	client, err := mail.NewClient(c.Options.Host,
		mail.WithPort(c.Options.Port),
		mail.WithTLSPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername(c.Options.FromAddress),
		mail.WithPassword(c.Options.Password),
	)
	if err != nil {
		return err
	}

	return client.DialAndSend(msg)
}
