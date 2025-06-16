package nats

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsBroker struct {
	Conn     *nats.Conn
	Js       nats.JetStreamContext
	NameMess string
}

func New(host, port, nameMess string) (*NatsBroker, error) {
	log.Println("broker: connection to Nats started")

	nc, err := nats.Connect(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, fmt.Errorf("failed to broker connect: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	log.Println("broker: connect to Nats successfully")
	return &NatsBroker{Conn: nc, Js: js, NameMess: nameMess}, nil
}

func (b *NatsBroker) Close() error {
	log.Println("broker: Nats stop started")

	if b.Conn == nil {
		return errors.New("broker connection is already closed")
	}

	b.Conn.Close()
	b.Conn = nil

	log.Println("broker: Nats stop successful")
	return nil
}

func (b *NatsBroker) CreateStream(streamName, subject string) error {
	cfg := &nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		MaxMsgs:  1000,
	}

	if _, err := b.Js.AddStream(cfg); err != nil {
		return fmt.Errorf("failed to create stream %q: %w", streamName, err)
	}

	return nil
}

func (b *NatsBroker) EnsureConsumer(streamName, subject, consumerName string) (jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cfg := jetstream.ConsumerConfig{
		Durable:      consumerName,
		AckPolicy:    jetstream.AckExplicitPolicy,
		MaxDeliver:   10,
		ReplayPolicy: jetstream.ReplayInstantPolicy,
	}

	jsm, err := jetstream.New(b.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream manager: %w", err)
	}

	consumer, err := jsm.CreateOrUpdateConsumer(ctx, b.NameMess, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to add consumer %q to stream %q: %w", consumerName, streamName, err)
	}

	return consumer, nil
}
