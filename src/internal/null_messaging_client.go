package internal

import (
	"context"
	"log/slog"
)

type NullMessagingClient struct {
	Logger *slog.Logger
}

func (c *NullMessagingClient) SendMsg(ctx context.Context, subject string, data any) error {
	c.Logger.Info("send msg",
		slog.Group("msg",
			slog.String("subject", subject),
			slog.Any("data", data),
		),
	)

	return nil
}

func (c *NullMessagingClient) SendPersistentMsg(ctx context.Context, subject string, data any) error {
	c.Logger.Info("send persistent msg",
		slog.Group("msg",
			slog.String("subject", subject),
			slog.Any("data", data),
		),
	)

	return nil
}

func (c *NullMessagingClient) PullPersistentMsgs(ctx context.Context, stream string, consumer string, batchSize int) ([]Message, error) {
	c.Logger.Info("pull persistent msgs",
		slog.Group("params",
			slog.String("stream", stream),
			slog.String("consumer", consumer),
			slog.Int("batchSize", batchSize),
		),
	)

	return nil, nil
}
