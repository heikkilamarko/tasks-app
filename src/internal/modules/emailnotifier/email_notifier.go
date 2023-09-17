package emailnotifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"
)

type EmailNotifier struct {
	Config          *shared.Config
	Logger          *slog.Logger
	MessagingClient shared.MessagingClient
	EmailClient     EmailClient
}

func (n *EmailNotifier) Run(ctx context.Context) error {
	return n.MessagingClient.SubscribePersistent(ctx, "tasks", "tasks", n.HandleMessage)
}

func (n *EmailNotifier) Close() error {
	return nil
}

func (n *EmailNotifier) HandleMessage(ctx context.Context, msg shared.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	switch msg.Subject() {
	case shared.SubjectTasksExpiring:
		return n.HandleTaskExpiringMessage(ctx, msg)
	case shared.SubjectTasksExpired:
		return n.HandleTaskExpiredMessage(ctx, msg)
	default:
		return n.HandleUnknownMessage(ctx, msg)
	}
}

func (n *EmailNotifier) HandleTaskExpiringMessage(ctx context.Context, msg shared.Message) error {
	var data shared.TaskExpiringMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.EmailClient.SendEmail(ctx, n.Config.EmailToAddress, "Task Expiring", "task_expiring.html", data.Task)
}

func (n *EmailNotifier) HandleTaskExpiredMessage(ctx context.Context, msg shared.Message) error {
	var data shared.TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.EmailClient.SendEmail(ctx, n.Config.EmailToAddress, "Task Expired", "task_expired.html", data.Task)
}

func (n *EmailNotifier) HandleUnknownMessage(ctx context.Context, msg shared.Message) error {
	n.Logger.Warn("handle unknown message",
		slog.Group("message",
			slog.String("subject", msg.Subject()),
			slog.Any("message", string(msg.Data())),
		),
	)

	return errors.New("unknown message")
}
