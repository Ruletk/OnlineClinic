package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(mainCtx, 5*time.Second)
	defer shutdownCancel()

	g, gCtx := errgroup.WithContext(shutdownCtx)

	// Main goroutine
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				fmt.Println("Main goroutine shutting down...")
				return nil
			default:
				// Simulate some work
				fmt.Println("Main goroutine doing work...")
				time.Sleep(1 * time.Second)
			}
		}
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
		fmt.Println("Error:", err)
	}
	fmt.Println("All done!")
}
