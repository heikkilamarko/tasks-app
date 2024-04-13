package emailnotifier

import (
	"context"
	"strings"
	"tasks-app/internal/shared"

	"github.com/wneessen/go-mail"
)

type SMTPEmailClient struct {
	Config *shared.Config
}

var _ EmailClient = (*SMTPEmailClient)(nil)

func (c *SMTPEmailClient) SendEmail(ctx context.Context, to string, subject string, templateName string, data any) error {
	var bodyBuilder strings.Builder
	if err := Templates.ExecuteTemplate(&bodyBuilder, templateName, data); err != nil {
		return err
	}
	body := bodyBuilder.String()

	msg := mail.NewMsg()

	if err := msg.FromFormat(c.Config.EmailNotifier.SMTPFromName, c.Config.EmailNotifier.SMTPFromAddress); err != nil {
		return err
	}

	if err := msg.To(to); err != nil {
		return err
	}

	msg.Subject(subject)

	msg.SetBodyString(mail.TypeTextHTML, body)

	var o []mail.Option

	if c.Config.EmailNotifier.SMTPPort == 25 {
		o = append(o,
			mail.WithTLSPortPolicy(mail.NoTLS),
		)
	} else {
		o = append(o,
			mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthLogin),
			mail.WithUsername(c.Config.EmailNotifier.SMTPFromAddress),
			mail.WithPassword(c.Config.EmailNotifier.SMTPPassword),
		)
	}

	client, err := mail.NewClient(c.Config.EmailNotifier.SMTPHost, o...)
	if err != nil {
		return err
	}

	return client.DialAndSendWithContext(ctx, msg)
}
