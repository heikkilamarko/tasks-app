package internal

import "time"

const (
	SubjectTasksExpiring   = "tasks.expiring"
	SubjectTasksExpired    = "tasks.expired"
	SubjectTasksUIExpiring = "tasks.ui.expiring"
	SubjectTasksUIExpired  = "tasks.ui.expired"
)

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type TaskExpiringMsg struct {
	Task *Task `json:"task"`
}

type TaskExpiredMsg struct {
	Task *Task `json:"task"`
}

func NewTask(name string, expiresAt *time.Time) *Task {
	return &Task{
		Name:      name,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

func (t *Task) Update(name string, expiresAt *time.Time) {
	t.Name = name
	t.ExpiresAt = expiresAt

	now := time.Now()
	t.UpdatedAt = &now
}
