package shared

import (
	"context"
	"mime/multipart"
)

type TaskAttachmentsRepository interface {
	Close() error
	GetAttachment(ctx context.Context, taskID int, name string) ([]byte, error)
	SaveAttachments(ctx context.Context, taskID int, fileHeaders []*multipart.FileHeader) error
	DeleteAttachments(ctx context.Context, taskID int, deleted map[int]string) error
	DeleteTask(ctx context.Context, taskID int) error
}
