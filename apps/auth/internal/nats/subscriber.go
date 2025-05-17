package nats

import (
	"context"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	nc *nats.Conn
}

func NewSubscriber(nc *nats.Conn) *Subscriber {
	return &Subscriber{nc: nc}
}

// SubscribeUserCreated SAMPLE CODE FOR SUBSCRIBING TO USER CREATED EVENTS
func (s *Subscriber) SubscribeUserCreated(ctx context.Context, handler func([]byte)) error {
	_, err := s.nc.Subscribe("user.created", func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}
