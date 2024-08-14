package ui

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
)

type NATSKVSessions[T authentication.Ctx] struct {
	js   jetstream.JetStream
	conn *nats.Conn
}

func NewNATSKVSessions[T authentication.Ctx](conn *nats.Conn) (*NATSKVSessions[T], error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	return &NATSKVSessions[T]{js, conn}, nil
}

func (s *NATSKVSessions[T]) Get(id string) (T, error) {
	ctx := context.Background()

	var session T

	kv, err := s.js.KeyValue(ctx, "sessions")
	if err == jetstream.ErrBucketNotFound {
		return session, errors.New("session bucket not found")
	}
	if err != nil {
		return session, err
	}

	entry, err := kv.Get(ctx, id)
	if err == jetstream.ErrKeyNotFound {
		return session, errors.New("session not found")
	}
	if err != nil {
		return session, err
	}

	if err := json.Unmarshal(entry.Value(), &session); err != nil {
		return session, err
	}

	return session, nil
}

func (s *NATSKVSessions[T]) Set(id string, session T) error {
	ctx := context.Background()

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	kv, err := s.js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:   "sessions",
		Replicas: 3,
		TTL:      30 * time.Minute,
	})
	if err != nil {
		return err
	}

	if _, err := kv.Put(ctx, id, data); err != nil {
		return err
	}

	return nil
}
