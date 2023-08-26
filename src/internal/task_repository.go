package internal

import "context"

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*Task, error)
	GetAll(ctx context.Context) ([]*Task, error)
}
