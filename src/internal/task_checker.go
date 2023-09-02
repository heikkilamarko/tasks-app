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
	seconds := time.Duration(tc.Config.TaskCheckIntervalSeconds) * time.Second

	go func() {
		for {
			select {
			case <-ctx.Done():
				tc.Logger.Info("exit task checker")
				return
			case <-time.After(seconds):
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

	tasks, err := tc.TaskRepository.GetExpiring(ctx, 24*time.Hour)
	if err != nil {
		return err
	}

	var errs []error
	for _, task := range tasks {
		err := tc.MessagingClient.SendPersistentMsg(ctx, SubjectTasksExpiring, TaskExpiringMsg{task})
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

	var errs []error
	for _, task := range tasks {
		err := tc.MessagingClient.SendPersistentMsg(ctx, SubjectTasksExpired, TaskExpiredMsg{task})
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
