package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type TaskChecker struct {
	Config         *Config
	Logger         *slog.Logger
	TaskRepository TaskRepository
}

func (tc *TaskChecker) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				tc.Logger.Info("exit task checker")
				return
			case <-time.After(10 * time.Second):
				if err := tc.CheckTasks(); err != nil {
					tc.Logger.Error("run checks", "error", err)
				}
			}
		}
	}()
}

func (tc *TaskChecker) CheckTasks() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return errors.Join(
		tc.CheckExpiringTasks(),
		tc.CheckExpiredTasks(),
	)
}

func (tc *TaskChecker) CheckExpiringTasks() error {
	tc.Logger.Info("check expiring tasks")
	return nil
}

func (tc *TaskChecker) CheckExpiredTasks() error {
	tc.Logger.Info("check expired tasks")
	return nil
}
