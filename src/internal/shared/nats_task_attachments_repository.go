package shared

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSTaskAttachmentsRepository struct {
	js     jetstream.JetStream
	conn   *nats.Conn
	logger *slog.Logger
}

var _ TaskAttachmentsRepository = (*NATSTaskAttachmentsRepository)(nil)

func NewNATSTaskAttachmentsRepository(conn *nats.Conn, logger *slog.Logger) (*NATSTaskAttachmentsRepository, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	return &NATSTaskAttachmentsRepository{js, conn, logger}, nil
}

func (repo *NATSTaskAttachmentsRepository) GetAttachment(ctx context.Context, taskID int, name string) ([]byte, error) {
	obs, err := repo.js.ObjectStore(ctx, repo.getBucketName(taskID))
	if err == jetstream.ErrBucketNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	data, err := obs.GetBytes(ctx, name)
	if err == jetstream.ErrObjectNotFound || err == jetstream.ErrBucketNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *NATSTaskAttachmentsRepository) SaveAttachments(ctx context.Context, taskID int, fileHeaders []*multipart.FileHeader) error {
	if len(fileHeaders) == 0 {
		return nil
	}

	obs, err := repo.js.CreateOrUpdateObjectStore(ctx, jetstream.ObjectStoreConfig{
		Bucket:   repo.getBucketName(taskID),
		Replicas: 3,
	})
	if err != nil {
		return err
	}

	for _, fileHeader := range fileHeaders {
		srcFile, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if _, err := obs.Put(ctx, jetstream.ObjectMeta{Name: fileHeader.Filename}, srcFile); err != nil {
			return err
		}
	}

	return nil
}

func (repo *NATSTaskAttachmentsRepository) DeleteAttachments(ctx context.Context, taskID int, deleted map[int]string) error {
	if len(deleted) == 0 {
		return nil
	}

	obs, err := repo.js.ObjectStore(ctx, repo.getBucketName(taskID))
	if err == jetstream.ErrBucketNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	for _, name := range deleted {
		if err := obs.Delete(ctx, name); err != nil && err != jetstream.ErrObjectNotFound {
			return err
		}
	}

	return nil
}

func (repo *NATSTaskAttachmentsRepository) DeleteTask(ctx context.Context, taskID int) error {
	if err := repo.js.DeleteObjectStore(ctx, repo.getBucketName(taskID)); err != nil && err != jetstream.ErrStreamNotFound {
		return err
	}
	return nil
}

func (repo *NATSTaskAttachmentsRepository) getBucketName(taskID int) string {
	return fmt.Sprintf("task_attachments_%d", taskID)
}
