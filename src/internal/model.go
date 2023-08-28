package internal

import "time"

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type TaskExpiringMsg struct {
	Task Task `json:"task"`
}

type TaskExpiredMsg struct {
	Task Task `json:"task"`
}
