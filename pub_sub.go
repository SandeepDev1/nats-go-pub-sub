package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func consumer(ctx context.Context) {

	nc, err := nats.Connect("nats://127.0.0.1:15000")
	if err != nil {
		log.Fatal("Failed to connect to NATS server:", err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS server on port 15000")

	subject := "logs"
	messages := make(chan *nats.Msg, 1000)

	subscription, err := nc.ChanSubscribe(subject, messages)
	if err != nil {
		log.Fatal("Failed to subscribe to subject:", err)
	}

	defer func() {
		subscription.Unsubscribe()
		close(messages)
	}()
	log.Println("Subscribed to", subject)

	for {
		select {
		case <-ctx.Done():
			log.Println("exiting from consumer")
			return
		case msg := <-messages:
			log.Println("received", string(msg.Data))
		}
	}
}

func producer(ctx context.Context) {
	nc, err := nats.Connect("nats://127.0.0.1:15000")
	if err != nil {
		log.Fatal("Failed to connect to NATS server:", err)
	}
	defer nc.Close()

	subject := "logs"

	i := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("exiting from producer")
			return
		default:
			i += 1
			message := fmt.Sprintf("message %v", i)
			time.Sleep(10 * time.Millisecond)
			err = nc.Publish(subject, []byte(message))
			if err != nil {
				log.Println("Failed to publish message:", err)
			}
		}
	}
}
