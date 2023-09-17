package shared

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSMsg struct {
	msg *nats.Msg
}

func (m *NATSMsg) Subject() string { return m.msg.Subject }
func (m *NATSMsg) Data() []byte    { return m.msg.Data }
func (m *NATSMsg) Ack() error      { return m.msg.Ack() }

type NATSMessagingClientOptions struct {
	NATSURL   string
	NATSToken string
	Logger    *slog.Logger
}

type NATSMessagingClient struct {
	Options NATSMessagingClientOptions
	conn    *nats.Conn
	js      jetstream.JetStream
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
				options.Logger.Error("nats error", "error", err)
				os.Exit(1)
			}),
	)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	return &NATSMessagingClient{options, conn, js}, nil
}

func (c *NATSMessagingClient) Close() error {
	return c.conn.Drain()
}

func (c *NATSMessagingClient) Send(ctx context.Context, subject string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.conn.Publish(subject, payload)
}

func (c *NATSMessagingClient) SendPersistent(ctx context.Context, subject string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = c.js.Publish(ctx, subject, payload)
	return err
}

func (c *NATSMessagingClient) Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, msg Message) error) error {
	sub, err := c.conn.SubscribeSync(subject)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			c.Options.Logger.Info("exit subscriber")
			return nil
		default:
			msg, err := sub.NextMsg(5 * time.Second)
			if err != nil {
				if !errors.Is(err, nats.ErrTimeout) {
					c.Options.Logger.Error("get next message", "error", err)
				}
				continue
			}

			if err := handler(ctx, &NATSMsg{msg}); err != nil {
				c.Options.Logger.Error("handle message", "error", err)
				continue
			}

			msg.Ack()
		}
	}
}

func (c *NATSMessagingClient) SubscribePersistent(ctx context.Context, stream string, consumer string, handler func(ctx context.Context, msg Message) error) error {
	con, err := c.js.Consumer(ctx, stream, consumer)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			c.Options.Logger.Info("exit persistent subscriber")
			return nil
		default:
			msg, err := con.Next(jetstream.FetchMaxWait(5 * time.Second))
			if err != nil {
				if !errors.Is(err, nats.ErrTimeout) {
					c.Options.Logger.Error("get next persistent message", "error", err)
				}
				continue
			}

			if err := handler(ctx, msg); err != nil {
				c.Options.Logger.Error("handle persistent message", "error", err)
				continue
			}

			msg.Ack()
		}
	}
}
