package nats

import (
	"errors"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Conn *nats.Conn
}

func New(host, port string) (*NatsBroker, error) {
	log.Println("broker: connection to Nats started")

	nc, err := nats.Connect(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, fmt.Errorf("failed to broker connect: %w", err)
	}

	log.Println("broker: connect to Nats successfully")
	return &NatsBroker{Conn: nc}, nil
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
