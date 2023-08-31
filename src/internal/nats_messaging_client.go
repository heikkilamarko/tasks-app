package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

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

func (c *NATSMessagingClient) PullPersistentMsgs(ctx context.Context, stream string, consumer string, batchSize int) ([]Message, error) {
	js, err := jetstream.New(c.Conn)
	if err != nil {
		return nil, err
	}

	con, err := js.Consumer(ctx, stream, consumer)
	if err != nil {
		return nil, err
	}

	batch, err := con.Fetch(batchSize, jetstream.FetchMaxWait(5*time.Second))
	if err != nil {
		return nil, err
	}

	if err := batch.Error(); err != nil {
		return nil, err
	}

	var messages []Message
	for msg := range batch.Messages() {
		messages = append(messages, msg)
	}

	return messages, nil
}
