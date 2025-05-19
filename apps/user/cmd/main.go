package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
)

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, gCtx := errgroup.WithContext(mainCtx)

	// Main goroutine
	g.Go(func() error {
		fmt.Println("Hello World!")
		return nil
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
		return
	}

	fmt.Println("All done!")
}
