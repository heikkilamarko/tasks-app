package emailnotifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"
)

type Module struct {
	Config          *shared.Config
	Logger          *slog.Logger
	MessagingClient shared.MessagingClient
	EmailClient     EmailClient
}

func (m *Module) Run(ctx context.Context) error {
	return m.MessagingClient.SubscribePersistent(ctx, "tasks", "tasks", m.handleMessage)
}

func (m *Module) handleMessage(ctx context.Context, msg shared.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	switch msg.Subject() {
	case shared.SubjectTasksExpiring:
		return m.handleTaskExpiringMessage(ctx, msg)
	case shared.SubjectTasksExpired:
		return m.handleTaskExpiredMessage(ctx, msg)
	default:
		return m.handleUnknownMessage(ctx, msg)
	}
}

func (m *Module) handleTaskExpiringMessage(ctx context.Context, msg shared.Message) error {
	var data shared.TaskExpiringMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return m.EmailClient.SendEmail(ctx, m.Config.EmailToAddress, "Task Expiring", "task_expiring", data.Task)
}

func (m *Module) handleTaskExpiredMessage(ctx context.Context, msg shared.Message) error {
	var data shared.TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}

	return m.EmailClient.SendEmail(ctx, m.Config.EmailToAddress, "Task Expired", "task_expired", data.Task)
}

func (m *Module) handleUnknownMessage(ctx context.Context, msg shared.Message) error {
	m.Logger.Warn("handle unknown message",
		slog.Group("message",
			slog.String("subject", msg.Subject()),
			slog.Any("message", string(msg.Data())),
		),
	)

	return errors.New("unknown message")
}