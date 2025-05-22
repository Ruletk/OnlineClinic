package nats

import (
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/Ruletk/OnlineClinic/pkg/proto/utils/gen/email"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Publisher interface {
	PublishEmailMessage(to, subject, message string) error
}

type NatsPublisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) Publisher {
	logging.Logger.Debug("Creating NATS publisher. I am not nil")
	return &NatsPublisher{nc: nc}
}

func (p *NatsPublisher) PublishEmailMessage(to, subject, message string) error {
	logging.Logger.Debug("Publishing email message to NATS")
	data, err := proto.Marshal(&email.SendEmailRequest{
		To:      to,
		Subject: subject,
		Message: message,
	})
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to marshal email message")
		return fmt.Errorf("failed to marshal email message: %w", err)
	}
	logging.Logger.Debugf("Publishing email message to NATS, email: %s", to[0:5])
	return p.nc.Publish("email.message", data)
}

func (p *NatsPublisher) publish(data []byte) error {
	if p.nc == nil {
		logging.Logger.Error("NATS connection is nil, cannot publish email message")
		return fmt.Errorf("NATS connection is nil")
	}

	logging.Logger.Debugf("Publishing data to NATS, length: %d", len(data))
	return p.nc.Publish("email.message", data)
}
