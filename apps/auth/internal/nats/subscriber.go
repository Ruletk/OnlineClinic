package nats

import (
	"context"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
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
	if s.nc == nil {
		logging.Logger.Error("NATS connection is nil, cannot subscribe to user created events")
		return fmt.Errorf("NATS connection is nil")
	}

	_, err := s.nc.Subscribe("user.created", func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}
