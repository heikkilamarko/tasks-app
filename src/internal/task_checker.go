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

	tc.MessagingClient.SendMsg(ctx, "tasks.expiring", TaskExpiringMsg{
		Task{
			ID:        1,
			Name:      "expiring task",
			CreatedAt: time.Now(),
		},
	})

	return nil
}

func (tc *TaskChecker) CheckExpiredTasks(ctx context.Context) error {
	tc.Logger.Info("check expired tasks")

	tc.MessagingClient.SendMsg(ctx, "tasks.expired", TaskExpiredMsg{
		Task{
			ID:        1,
			Name:      "expired task",
			CreatedAt: time.Now(),
		},
	})

	return nil
}
