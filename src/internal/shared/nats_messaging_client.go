package shared

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSMsg struct {
	msg *nats.Msg
}

func (m *NATSMsg) Subject() string                        { return m.msg.Subject }
func (m *NATSMsg) Data() []byte                           { return m.msg.Data }
func (m *NATSMsg) Ack() error                             { return m.msg.Ack() }
func (m *NATSMsg) Nak() error                             { return m.msg.Nak() }
func (m *NATSMsg) NakWithDelay(delay time.Duration) error { return m.msg.NakWithDelay(delay) }

type NATSMessagingClient struct {
	js     jetstream.JetStream
	conn   *nats.Conn
	logger *slog.Logger
}

var _ MessagingClient = (*NATSMessagingClient)(nil)

func NewNATSMessagingClient(conn *nats.Conn, logger *slog.Logger) (*NATSMessagingClient, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	return &NATSMessagingClient{js, conn, logger}, nil
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
			return nil
		default:
			msg, err := sub.NextMsg(5 * time.Second)
			if err != nil {
				if !errors.Is(err, nats.ErrTimeout) {
					c.logger.Error("get next message", "error", err)
				}
				continue
			}

			if err := handler(ctx, &NATSMsg{msg}); err != nil {
				c.logger.Error("handle message", "error", err)
			}
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
			return nil
		default:
			msg, err := con.Next(jetstream.FetchMaxWait(5 * time.Second))
			if err != nil {
				if !errors.Is(err, nats.ErrTimeout) {
					c.logger.Error("get next persistent message", "error", err)
				}
				continue
			}

			if err := handler(ctx, msg); err != nil {
				c.logger.Error("handle persistent message", "error", err)
			}
		}
	}
}
