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
	TaskRepository  shared.TaskRepository
	MessagingClient shared.MessagingClient
}

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

	count, err := m.TaskRepository.DeleteCompleted(ctx, m.Config.TaskChecker.DeleteWindow)
	if err != nil {
		return err
	}

	if 0 < count {
		m.Logger.Info("found completed tasks", slog.Int64("count", count))
	}

	return nil
}

func (m *Module) checkExpiringTasks(ctx context.Context) error {
	m.Logger.Info("check expiring tasks")

	tasks, err := m.TaskRepository.GetExpiring(ctx, m.Config.TaskChecker.ExpiringWindow)
	if err != nil {
		return err
	}

	count := len(tasks)
	if 0 < count {
		m.Logger.Info("found expiring tasks", slog.Int("count", count))
	}

	var errs []error
	for _, task := range tasks {
		if err := m.MessagingClient.SendPersistent(ctx, shared.SubjectTasksExpiring, shared.TaskExpiringMsg{Task: task}); err != nil {
			errs = append(errs, err)
			continue
		}

		task.SetExpiringInfoAt()
		err := m.TaskRepository.Update(ctx, task)
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (m *Module) checkExpiredTasks(ctx context.Context) error {
	m.Logger.Info("check expired tasks")

	tasks, err := m.TaskRepository.GetExpired(ctx)
	if err != nil {
		return err
	}

	count := len(tasks)
	if 0 < count {
		m.Logger.Info("found expired tasks", slog.Int("count", count))
	}

	var errs []error
	for _, task := range tasks {
		if err := m.MessagingClient.SendPersistent(ctx, shared.SubjectTasksExpired, shared.TaskExpiredMsg{Task: task}); err != nil {
			errs = append(errs, err)
			continue
		}

		task.SetExpiredInfoAt()
		err := m.TaskRepository.Update(ctx, task)
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
