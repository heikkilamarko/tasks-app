package shared

import (
	"context"
	"mime/multipart"
)

type TaskAttachmentsRepository interface {
	SaveAttachments(ctx context.Context, taskID int, fileHeaders []*multipart.FileHeader) error
	DeleteTask(ctx context.Context, taskID int) error
}
