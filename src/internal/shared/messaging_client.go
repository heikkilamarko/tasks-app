package shared

import (
	"context"
	"time"
)

type Message interface {
	Subject() string
	Data() []byte
	Ack() error
	Nak() error
	NakWithDelay(delay time.Duration) error
}

type MessagingClient interface {
	Send(ctx context.Context, subject string, data any) error
	SendPersistent(ctx context.Context, subject string, data any) error
	Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, msg Message) error) error
	SubscribePersistent(ctx context.Context, stream string, consumer string, handler func(ctx context.Context, msg Message) error) error
}
