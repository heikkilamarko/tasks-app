package emailnotifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"tasks-app/internal/shared"
	"time"
)

type Module struct {
	Config          *shared.Config
	Logger          *slog.Logger
	MessagingClient shared.MessagingClient
	EmailResolver   EmailResolver
	EmailClient     EmailClient
	validator       *shared.SchemaValidator
}

func (m *Module) Run(ctx context.Context) error {
	m.validator = shared.NewSchemaValidator(SchemasFS)

	return m.MessagingClient.SubscribePersistent(ctx, "tasks", "tasks", m.handleMessage)
}

func (m *Module) handleMessage(ctx context.Context, msg shared.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			m.NakMessage(msg)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	sub := msg.Subject()

	if strings.HasPrefix(sub, "task.") && strings.HasSuffix(sub, ".expiring") {
		return m.handleTaskExpiringMessage(ctx, msg)
	} else if strings.HasPrefix(sub, "task.") && strings.HasSuffix(sub, ".expired") {
		return m.handleTaskExpiredMessage(ctx, msg)
	} else {
		return m.handleUnknownMessage(ctx, msg)
	}
}

func (m *Module) handleTaskExpiringMessage(ctx context.Context, msg shared.Message) error {
	if err := m.validator.ValidateBytes("schemas/task.expiring.json", msg.Data()); err != nil {
		m.NakMessage(msg)
		return err
	}

	var data shared.TaskExpiringMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		m.NakMessage(msg)
		return err
	}

	to, err := m.EmailResolver.ResolveEmail(data.Task.UserID)
	if err != nil {
		m.NakMessage(msg)
		return err
	}

	if err := m.EmailClient.SendEmail(ctx, to, "Task Expiring", "task_expiring.html", data.Task); err != nil {
		m.NakMessage(msg)
		return err
	}

	m.AckMessage(msg)
	return nil
}

func (m *Module) handleTaskExpiredMessage(ctx context.Context, msg shared.Message) error {
	if err := m.validator.ValidateBytes("schemas/task.expired.json", msg.Data()); err != nil {
		m.NakMessage(msg)
		return err
	}

	var data shared.TaskExpiredMsg
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		m.NakMessage(msg)
		return err
	}

	to, err := m.EmailResolver.ResolveEmail(data.Task.UserID)
	if err != nil {
		m.NakMessage(msg)
		return err
	}

	if err := m.EmailClient.SendEmail(ctx, to, "Task Expired", "task_expired.html", data.Task); err != nil {
		m.NakMessage(msg)
		return err
	}

	m.AckMessage(msg)
	return nil
}

func (m *Module) handleUnknownMessage(_ context.Context, msg shared.Message) error {
	m.AckMessage(msg)

	m.Logger.Warn("handle unknown message",
		slog.Group("message",
			slog.String("subject", msg.Subject()),
			slog.Any("message", string(msg.Data())),
		),
	)

	return errors.New("unknown message")
}

func (m *Module) AckMessage(msg shared.Message) {
	if err := msg.Ack(); err != nil {
		m.Logger.Error("message ack failed")
	}
}

func (m *Module) NakMessage(msg shared.Message) {
	if err := msg.NakWithDelay(4 * time.Second); err != nil {
		m.Logger.Error("message nak failed")
	}
}
