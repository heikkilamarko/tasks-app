package internal

import (
	"context"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Task, error)
	GetActive(ctx context.Context) ([]*Task, error)
	GetCompleted(ctx context.Context) ([]*Task, error)
	GetExpiring(ctx context.Context, d time.Duration) ([]*Task, error)
	GetExpired(ctx context.Context) ([]*Task, error)
}
