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

func (n *UINotifier) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				n.Logger.Info("exit ui notifier")
				return
			default:
				if err := n.HandleMessages(ctx); err != nil {
					n.Logger.Error("handle messages", "error", err)
				}
			}
		}
	}()
}

func (n *UINotifier) HandleMessages(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	msgs, err := n.MessagingClient.PullPersistentMsgs(ctx, "tasks", "tasks", 10)
	if err != nil {
		return err
	}

	var errs []error

	for _, msg := range msgs {
		if err := n.HandleMessage(ctx, msg); err != nil {
			errs = append(errs, err)
			continue
		}
		msg.Ack()
	}

	return errors.Join(errs...)
}

func (n *UINotifier) HandleMessage(ctx context.Context, msg Message) error {
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

	return n.MessagingClient.SendMsg(ctx, SubjectTasksUIExpiring, data)
}

func (n *UINotifier) HandleTaskExpiredMessage(ctx context.Context, msg Message) error {
	var data TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return n.MessagingClient.SendMsg(ctx, SubjectTasksUIExpired, data)
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
