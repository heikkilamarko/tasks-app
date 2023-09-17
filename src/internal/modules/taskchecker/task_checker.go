package taskchecker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tasks-app/internal/shared"
	"time"
)

type TaskChecker struct {
	Config          *shared.Config
	Logger          *slog.Logger
	TaskRepository  shared.TaskRepository
	MessagingClient shared.MessagingClient
}

func (tc *TaskChecker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			tc.Logger.Info("exit task checker")
			return nil
		case <-time.After(tc.Config.TaskCheckInterval):
			if err := tc.CheckTasks(ctx); err != nil {
				tc.Logger.Error("run checks", "error", err)
			}
		}
	}
}

func (tc *TaskChecker) Close() error {
	return nil
}

func (tc *TaskChecker) CheckTasks(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return errors.Join(
		tc.CheckCompletedTasks(ctx),
		tc.CheckExpiringTasks(ctx),
		tc.CheckExpiredTasks(ctx),
	)
}

func (tc *TaskChecker) CheckCompletedTasks(ctx context.Context) error {
	tc.Logger.Info("check completed tasks")

	count, err := tc.TaskRepository.DeleteCompleted(ctx, tc.Config.TaskCheckDeleteWindow)
	if err != nil {
		return err
	}

	if 0 < count {
		tc.Logger.Info("found completed tasks", slog.Int64("count", count))
	}

	return nil
}

func (tc *TaskChecker) CheckExpiringTasks(ctx context.Context) error {
	tc.Logger.Info("check expiring tasks")

	tasks, err := tc.TaskRepository.GetExpiring(ctx, tc.Config.TaskCheckExpiringWindow)
	if err != nil {
		return err
	}

	count := len(tasks)
	if 0 < count {
		tc.Logger.Info("found expiring tasks", slog.Int("count", count))
	}

	var errs []error
	for _, task := range tasks {
		if err := tc.MessagingClient.SendPersistent(ctx, shared.SubjectTasksExpiring, shared.TaskExpiringMsg{Task: task}); err != nil {
			errs = append(errs, err)
			continue
		}

		task.SetExpiringInfoAt()
		err := tc.TaskRepository.Update(ctx, task)
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (tc *TaskChecker) CheckExpiredTasks(ctx context.Context) error {
	tc.Logger.Info("check expired tasks")

	tasks, err := tc.TaskRepository.GetExpired(ctx)
	if err != nil {
		return err
	}

	count := len(tasks)
	if 0 < count {
		tc.Logger.Info("found expired tasks", slog.Int("count", count))
	}

	var errs []error
	for _, task := range tasks {
		if err := tc.MessagingClient.SendPersistent(ctx, shared.SubjectTasksExpired, shared.TaskExpiredMsg{Task: task}); err != nil {
			errs = append(errs, err)
			continue
		}

		task.SetExpiredInfoAt()
		err := tc.TaskRepository.Update(ctx, task)
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
