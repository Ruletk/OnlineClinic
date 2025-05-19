package subscribers

import (
    "context"
    "github.com/Ruletk/OnlineClinic/pkg/logging"
    natspb "github.com/Ruletk/OnlineClinic/pkg/proto/nats/gen"
    "github.com/nats-io/nats.go"
    "google.golang.org/protobuf/proto"
    "notification/internal/models"
    "notification/internal/repositories/email"
)

// NatsService is a struct that represents the NATS service.
type NatsService struct {
    emailSender email.Sender
    client      *nats.Conn
}

// NewNatsService creates a new NatsService instance.
func NewNatsService(client *nats.Conn, emailSender email.Sender) *NatsService {
    return &NatsService{
        client:      client,
        emailSender: emailSender,
    }
}

func (n *NatsService) InitNatsSubscriber(ctx context.Context) error {
    sub, err := n.client.Subscribe("email.message", n.emailSubscriber)
    if err != nil {
        logging.Logger.Error("Failed to subscribe to NATS topic: %v", err)
        return err
    }
    defer func(sub *nats.Subscription) {
        _ = sub.Unsubscribe()
    }(sub)

    logging.Logger.Info("NATS subscriber initialized successfully")

    <-ctx.Done()
    logging.Logger.Debug("NATS subscriber received shutdown signal")
    return nil
}

func (n *NatsService) emailSubscriber(msg *nats.Msg) {
    logging.Logger.Debug("Received email message request")
    var natsEmailMessage natspb.EmailMessage
    if err := proto.Unmarshal(msg.Data, &natsEmailMessage); err != nil {
        logging.Logger.Error("Failed to unmarshal email message: %v", err)
        return
    }
    logging.Logger.Debugf("Unmarshalled email message, to: %s", natsEmailMessage.To[0:10])
    emailMsg := models.EmailMessage{
        To:      natsEmailMessage.To,
        Subject: natsEmailMessage.Subject,
        Message: natsEmailMessage.Message,
    }
    if err := n.emailSender.Send(emailMsg); err != nil {
        logging.Logger.Error("Failed to send email: %v", err)
        return
    }
    logging.Logger.Infof("Email sent successfully to: %s", natsEmailMessage.To[0:10])
}
