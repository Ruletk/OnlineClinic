package main

import (
    "context"
    "fmt"
    "github.com/Ruletk/OnlineClinic/pkg/config"
    "github.com/Ruletk/OnlineClinic/pkg/logging"
    "github.com/nats-io/nats.go"
    "golang.org/x/sync/errgroup"
    "notification/internal/repositories/email"
    "notification/internal/subscribers"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    cfg, err := config.GetDefaultConfiguration()
    if err != nil {
        fmt.Printf("Failed to load configuration: %v\n", err)
        panic(err)
    }
    logging.InitLogger(*cfg)
    logging.Logger.Infof("Starting Notification Service")

    emailSender := email.NewMockEmailSender()

    logging.Logger.Debugf("Creating NATS connection: %s", cfg.Nats.Url)
    conn, err := nats.Connect(cfg.Nats.Url, nats.Name("Notification Service"))
    if err != nil {
        logging.Logger.WithError(err).Error("Failed to connect to NATS")
        panic(err)
    }
    defer conn.Close()

    natsService := subscribers.NewNatsService(conn, emailSender)

    logging.Logger.Info("NATS subscriber initialized and running")

    shutdownCtx, shutdownCancel := context.WithTimeout(mainCtx, 5*time.Second)
    defer shutdownCancel()

    g, gCtx := errgroup.WithContext(shutdownCtx)

    // Main goroutine
    g.Go(func() error {
        return natsService.InitNatsSubscriber(gCtx)
    })

    // Shutdown goroutines
    g.Go(func() error {
        // Simulate some work
        fmt.Println("Doing some work...")
        // Wait for the context to be done
        <-gCtx.Done()
        fmt.Println("Shutting down...")
        return nil
    })

    if err := g.Wait(); err != nil {
        logging.Logger.WithError(err).Error("Error during shutdown")
    } else {
        logging.Logger.Info("Service stopped successfully")
    }
    logging.Logger.Info("Service stopped")
}
