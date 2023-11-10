package shared

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/nats-io/nats.go"
)

type NATSTaskAttachmentsRepository struct {
	Config *Config
	Logger *slog.Logger
	conn   *nats.Conn
	js     nats.JetStreamContext
}

func NewNATSTaskAttachmentsRepository(config *Config, logger *slog.Logger) (*NATSTaskAttachmentsRepository, error) {
	conn, err := nats.Connect(
		config.Shared.NATSURL,
		nats.UserInfo(config.Shared.NATSUser, config.Shared.NATSPassword),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(
			func(_ *nats.Conn, err error) {
				logger.Info("nats disconnected", "reason", err)
			}),
		nats.ReconnectHandler(
			func(c *nats.Conn) {
				logger.Info("nats reconnected", "address", c.ConnectedUrl())
			}),
		nats.ErrorHandler(
			func(_ *nats.Conn, _ *nats.Subscription, err error) {
				logger.Error("nats error", "error", err)
				os.Exit(1)
			}),
	)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	return &NATSTaskAttachmentsRepository{config, logger, conn, js}, nil
}

func (repo *NATSTaskAttachmentsRepository) Close() error {
	return repo.conn.Drain()
}

func (repo *NATSTaskAttachmentsRepository) GetAttachment(ctx context.Context, taskID int, name string) ([]byte, error) {
	os, err := repo.js.ObjectStore(fmt.Sprintf("attachments_%d", taskID))
	if err != nil {
		return nil, err
	}

	data, err := os.GetBytes(name)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *NATSTaskAttachmentsRepository) SaveAttachments(ctx context.Context, taskID int, fileHeaders []*multipart.FileHeader) error {
	os, err := repo.js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: fmt.Sprintf("attachments_%d", taskID),
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

		_, err = os.Put(&nats.ObjectMeta{Name: fileHeader.Filename}, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *NATSTaskAttachmentsRepository) DeleteAttachments(ctx context.Context, taskID int, deleted map[int]string) error {
	os, err := repo.js.ObjectStore(fmt.Sprintf("attachments_%d", taskID))
	if err != nil {
		return err
	}

	for _, name := range deleted {
		if err := os.Delete(name); err != nil {
			return err
		}
	}

	return nil
}

func (repo *NATSTaskAttachmentsRepository) DeleteTask(ctx context.Context, taskID int) error {
	return repo.js.DeleteObjectStore(fmt.Sprintf("attachments_%d", taskID))
}
