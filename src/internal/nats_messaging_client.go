package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSMessagingClientOptions struct {
	NATSURL   string
	NATSToken string
	Logger    *slog.Logger
}

type NATSMessagingClient struct {
	Options NATSMessagingClientOptions
	Conn    *nats.Conn
}

func NewNATSMessagingClient(options NATSMessagingClientOptions) (*NATSMessagingClient, error) {
	conn, err := nats.Connect(
		options.NATSURL,
		nats.Token(options.NATSToken),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(
			func(_ *nats.Conn, err error) {
				options.Logger.Info("nats disconnected", "reason", err)
			}),
		nats.ReconnectHandler(
			func(c *nats.Conn) {
				options.Logger.Info("nats reconnected", "address", c.ConnectedUrl())
			}),
		nats.ErrorHandler(
			func(_ *nats.Conn, _ *nats.Subscription, err error) {
				options.Logger.Error("nats error", "err", err)
				os.Exit(1)
			}),
	)
	if err != nil {
		return nil, err
	}

	return &NATSMessagingClient{options, conn}, nil
}

func (c *NATSMessagingClient) SendMsg(ctx context.Context, subject string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.Conn.Publish(subject, payload)
}

func (c *NATSMessagingClient) SendPersistentMsg(ctx context.Context, subject string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js, err := jetstream.New(c.Conn)
	if err != nil {
		return err
	}

	_, err = js.Publish(ctx, subject, payload)
	return err
}
