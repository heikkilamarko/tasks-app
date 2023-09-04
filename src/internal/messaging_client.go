package internal

import "context"

type Message interface {
	Subject() string
	Data() []byte
	Ack() error
}

type MessagingClient interface {
	Close() error
	Send(ctx context.Context, subject string, data any) error
	SendPersistent(ctx context.Context, subject string, data any) error
	Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, msg Message) error) error
	SubscribePersistent(ctx context.Context, stream string, consumer string, handler func(ctx context.Context, msg Message) error) error
}
