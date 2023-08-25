package internal

import "time"

type Task struct {
	ID          int
	Name        string
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	CompletedAt *time.Time
}
