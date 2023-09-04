package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

type UINotifier struct {
	Config          *Config
	Logger          *slog.Logger
	MessagingClient MessagingClient
}

func (n *UINotifier) Run(ctx context.Context) error {
	return n.MessagingClient.Subscribe(ctx, "tasks.*", n.HandleMessage)
}

func (n *UINotifier) Close() error {
	return nil
}

func (n *UINotifier) HandleMessage(ctx context.Context, msg Message) (err error) {
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

func (n *UINotifier) HandleTaskExpiringMessage(ctx context.Context, msg Message) error {
	var data TaskExpiringMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.MessagingClient.Send(ctx, SubjectTasksUIExpiring, data)
}

func (n *UINotifier) HandleTaskExpiredMessage(ctx context.Context, msg Message) error {
	var data TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.MessagingClient.Send(ctx, SubjectTasksUIExpired, data)
}

func (n *UINotifier) HandleUnknownMessage(ctx context.Context, msg Message) error {
	n.Logger.Warn("handle unknown message",
		slog.Group("message",
			slog.String("subject", msg.Subject()),
			slog.Any("message", string(msg.Data())),
		),
	)

	return errors.New("unknown message")
}
