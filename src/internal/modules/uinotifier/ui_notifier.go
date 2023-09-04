package uinotifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"
)

type UINotifier struct {
	Config          *shared.Config
	Logger          *slog.Logger
	MessagingClient shared.MessagingClient
}

func (*UINotifier) Name() string { return "uinotifier" }

func (n *UINotifier) Run(ctx context.Context) error {
	return n.MessagingClient.Subscribe(ctx, "tasks.*", n.HandleMessage)
}

func (n *UINotifier) Close() error {
	return nil
}

func (n *UINotifier) HandleMessage(ctx context.Context, msg shared.Message) (err error) {
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

func (n *UINotifier) HandleTaskExpiringMessage(ctx context.Context, msg shared.Message) error {
	var data shared.TaskExpiringMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.MessagingClient.Send(ctx, shared.SubjectTasksUIExpiring, data)
}

func (n *UINotifier) HandleTaskExpiredMessage(ctx context.Context, msg shared.Message) error {
	var data shared.TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.MessagingClient.Send(ctx, shared.SubjectTasksUIExpired, data)
}

func (n *UINotifier) HandleUnknownMessage(ctx context.Context, msg shared.Message) error {
	n.Logger.Warn("handle unknown message",
		slog.Group("message",
			slog.String("subject", msg.Subject()),
			slog.Any("message", string(msg.Data())),
		),
	)

	return errors.New("unknown message")
}
