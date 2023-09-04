package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

type EmailNotifier struct {
	Config          *Config
	Logger          *slog.Logger
	MessagingClient MessagingClient
	EmailClient     EmailClient
}

func (n *EmailNotifier) Run(ctx context.Context) {
	go n.MessagingClient.SubscribePersistent(ctx, "tasks", "tasks", n.HandleMessage)
}

func (n *EmailNotifier) HandleMessage(ctx context.Context, msg Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	switch msg.Subject() {
	case SubjectTasksExpiring:
		return n.HandleTaskExpiringMessage(ctx, msg)
	case SubjectTasksExpired:
		return n.HandleTaskExpiredMessage(ctx, msg)
	default:
		return n.HandleUnknownMessage(ctx, msg)
	}
}

func (n *EmailNotifier) HandleTaskExpiringMessage(ctx context.Context, msg Message) error {
	var data TaskExpiringMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.EmailClient.SendEmail(ctx, n.Config.EmailToAddress, "Task Expiring", "task_expiring.html", data)
}

func (n *EmailNotifier) HandleTaskExpiredMessage(ctx context.Context, msg Message) error {
	var data TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.EmailClient.SendEmail(ctx, n.Config.EmailToAddress, "Task Expired", "task_expired.html", data)
}

func (n *EmailNotifier) HandleUnknownMessage(ctx context.Context, msg Message) error {
	n.Logger.Warn("handle unknown message",
		slog.Group("message",
			slog.String("subject", msg.Subject()),
			slog.Any("message", string(msg.Data())),
		),
	)

	return errors.New("unknown message")
}
