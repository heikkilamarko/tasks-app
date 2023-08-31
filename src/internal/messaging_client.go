package internal

import "context"

type Message interface {
	Subject() string
	Data() []byte
	Ack() error
}

type MessagingClient interface {
	SendMsg(ctx context.Context, subject string, data any) error
	SendPersistentMsg(ctx context.Context, subject string, data any) error
	PullPersistentMsgs(ctx context.Context, stream string, consumer string, batchSize int) ([]Message, error)
}
