package shared

import (
	"context"
	"time"
)

type TaskRepository interface {
	Close() error
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	UpdateAttachments(ctx context.Context, taskID int, attachmentNames []string) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Task, error)
	GetActive(ctx context.Context) ([]*Task, error)
	GetCompleted(ctx context.Context) ([]*Task, error)
	GetExpiring(ctx context.Context, d time.Duration) ([]*Task, error)
	GetExpired(ctx context.Context) ([]*Task, error)
	DeleteCompleted(ctx context.Context, d time.Duration) (int64, error)
}
