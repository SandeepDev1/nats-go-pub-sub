package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// listen for interrupts to exit gracefully
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
		<-sigChannel
		cancel()
	}()

	// create the local server
	server := createServer()
	defer server.Shutdown()

	// register the consumer
	go consumer(ctx)

	// register the producer
	go producer(ctx)
	<-ctx.Done()

	log.Println("server shutdown completed")
	log.Println("exiting gracefully")
}
