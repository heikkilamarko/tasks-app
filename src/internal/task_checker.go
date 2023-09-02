package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type TaskChecker struct {
	Config          *Config
	Logger          *slog.Logger
	TaskRepository  TaskRepository
	MessagingClient MessagingClient
}

func (tc *TaskChecker) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				tc.Logger.Info("exit task checker")
				return
			case <-time.After(tc.Config.TaskCheckInterval):
				if err := tc.CheckTasks(ctx); err != nil {
					tc.Logger.Error("run checks", "error", err)
				}
			}
		}
	}()
}

func (tc *TaskChecker) CheckTasks(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return errors.Join(
		tc.CheckExpiringTasks(ctx),
		tc.CheckExpiredTasks(ctx),
	)
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
		if err := tc.MessagingClient.SendPersistentMsg(ctx, SubjectTasksExpiring, TaskExpiringMsg{task}); err != nil {
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
		if err := tc.MessagingClient.SendPersistentMsg(ctx, SubjectTasksExpired, TaskExpiredMsg{task}); err != nil {
			errs = append(errs, err)
			continue
		}

		task.SetExpiredInfoAt()
		err := tc.TaskRepository.Update(ctx, task)
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
