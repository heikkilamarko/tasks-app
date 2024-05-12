package shared

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Task struct {
	ID             int         `json:"id"`
	UserID         string      `json:"user_id"`
	Name           string      `json:"name"`
	ExpiresAt      *time.Time  `json:"expires_at"`
	ExpiringInfoAt *time.Time  `json:"expiring_info_at"`
	ExpiredInfoAt  *time.Time  `json:"expired_info_at"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at"`
	CompletedAt    *time.Time  `json:"completed_at"`
	Attachments    Attachments `json:"attachments"`
}

type Attachments []*Attachment

type Attachment struct {
	ID        int        `json:"id"`
	TaskID    int        `json:"task_id"`
	FileName  string     `json:"file_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type TaskExpiringMsg struct {
	Task *Task `json:"task"`
}

type TaskExpiredMsg struct {
	Task *Task `json:"task"`
}

func NewTask(name string, expiresAt *time.Time) *Task {
	now := time.Now().UTC()

	return &Task{
		Name:      name,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}
}

func (t *Task) Update(name string, expiresAt *time.Time) {
	now := time.Now().UTC()

	t.Name = name
	t.ExpiresAt = expiresAt
	t.ExpiringInfoAt = nil
	t.ExpiredInfoAt = nil
	t.UpdatedAt = &now
}

func (t *Task) SetExpiringInfoAt() {
	now := time.Now().UTC()

	t.ExpiringInfoAt = &now
	t.UpdatedAt = &now
}

func (t *Task) SetExpiredInfoAt() {
	now := time.Now().UTC()

	t.ExpiredInfoAt = &now
	t.UpdatedAt = &now
}

func (t *Task) SetCompleted() {
	now := time.Now().UTC()

	t.CompletedAt = &now
	t.UpdatedAt = &now
}

func (a *Attachments) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attachments) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte")
	}

	return json.Unmarshal(b, a)
}
