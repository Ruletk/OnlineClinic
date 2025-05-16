package main

import (
	"context"
	"fmt"
	config2 "github.com/Ruletk/OnlineClinic/pkg/config"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/nats-io/nats.go"
	"notification/internal/repositories/email"
	"notification/internal/subscribers"
	"os/signal"
	"syscall"
)

func main() {
	config, err := config2.GetDefaultConfiguration()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		panic(err)
	}
	logging.InitLogger(*config)
	logging.Logger.Infof("Starting Notification Service")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := startApp(ctx, config); err != nil {
		logging.Logger.Errorf("Application error: %v", err)
	}

	logging.Logger.Warn("Shutting down gracefully...")
	logging.Logger.Info("Service stopped")
}

func startApp(ctx context.Context, config *config2.Config) error {
	emailSender := email.NewMockEmailSender()

	conn, err := nats.Connect(config.Nats.Url, nats.Name("Notification Service"))
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer conn.Close()

	natsService := subscribers.NewNatsService(conn, emailSender)
	natsService.InitNatsSubscriber()

	logging.Logger.Info("NATS subscriber initialized and running")

	<-ctx.Done()

	logging.Logger.Info("Context cancelled, shutting down NATS service")
	return nil
}
