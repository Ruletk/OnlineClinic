package main

import (
	"fmt"
	config2 "github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/nats-io/nats.go"
	"notification/internal/repositories/email"
	"notification/internal/subscribers"
)

func main() {
	config, err := config2.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		panic(err)
	}
	logging.InitLogger(*config)
	logging.Logger.Infof("Starting Notification Service")

	emailSender := email.NewMockEmailSender()
	logging.Logger.Debugf("Creating NATS connection: %s", config.Nats.Url)
	conn, err := nats.Connect(config.Nats.Url, nats.Name("Notification Service"))
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to connect to NATS")
		panic(err)
	}
	defer conn.Close()

	natsService := subscribers.NewNatsService(conn, emailSender)
	natsService.InitNatsSubscriber()

	logging.Logger.Info("NATS subscriber initialized and running")

	// Просто ждем сигнала завершения
	wait := make(chan struct{})
	<-wait

	logging.Logger.Info("Service stopped")
}
