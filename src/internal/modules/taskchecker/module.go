package taskchecker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"
	"time"
)

type Module struct {
	Config          *shared.Config
	Logger          *slog.Logger
	TxManager       shared.TxManager
	MessagingClient shared.MessagingClient
}

var _ shared.AppModule = (*Module)(nil)

func (m *Module) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(m.Config.TaskChecker.CheckInterval):
			if err := m.checkTasks(ctx); err != nil {
				m.Logger.Error("run checks", "error", err)
			}
		}
	}
}

func (m *Module) checkTasks(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return errors.Join(
		m.checkCompletedTasks(ctx),
		m.checkExpiringTasks(ctx),
		m.checkExpiredTasks(ctx),
	)
}

func (m *Module) checkCompletedTasks(ctx context.Context) error {
	m.Logger.Info("check completed tasks")

	return m.TxManager.RunInTx(func(txc shared.TxContext) error {
		count, err := txc.TaskRepository.DeleteCompleted(ctx, m.Config.TaskChecker.DeleteWindow)
		if err != nil {
			return err
		}

		if 0 < count {
			m.Logger.Info("found completed tasks", slog.Int64("count", count))
		}

		return nil
	})
}

func (m *Module) checkExpiringTasks(ctx context.Context) error {
	m.Logger.Info("check expiring tasks")

	var tasks []*shared.Task
	var err error

	err = m.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err = txc.TaskRepository.GetExpiring(ctx, m.Config.TaskChecker.ExpiringWindow)
		return err
	})

	if err != nil {
		return err
	}

	count := len(tasks)
	if 0 < count {
		m.Logger.Info("found expiring tasks", slog.Int("count", count))
	}

	var errs []error

	for _, task := range tasks {
		task.SetExpiringInfoAt()

		err = m.TxManager.RunInTx(func(txc shared.TxContext) error {
			if err := txc.TaskRepository.Update(ctx, task); err != nil {
				return err
			}

			if err := m.MessagingClient.SendPersistent(ctx, fmt.Sprintf("task.%s.%d.expiring", task.UserID, task.ID), shared.TaskExpiringMsg{Task: task}); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (m *Module) checkExpiredTasks(ctx context.Context) error {
	m.Logger.Info("check expired tasks")

	var tasks []*shared.Task
	var err error

	err = m.TxManager.RunInTx(func(txc shared.TxContext) error {
		tasks, err = txc.TaskRepository.GetExpired(ctx)
		return err
	})

	count := len(tasks)
	if 0 < count {
		m.Logger.Info("found expired tasks", slog.Int("count", count))
	}

	var errs []error

	for _, task := range tasks {
		task.SetExpiredInfoAt()

		err = m.TxManager.RunInTx(func(txc shared.TxContext) error {
			if err := txc.TaskRepository.Update(ctx, task); err != nil {
				return err
			}

			if err := m.MessagingClient.SendPersistent(ctx, fmt.Sprintf("task.%s.%d.expired", task.UserID, task.ID), shared.TaskExpiredMsg{Task: task}); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
