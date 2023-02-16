package main

import (
	"context"
	"fmt"
	"goRabbit/receive"
	"goRabbit/send"
	"log"
	"time"
)

func main() {

	// create a publisher
	publisher := send.NewQueuePublisher("amqp://guest:guest@localhost:8080/", "hello")

	// create a receiver
	receiver, err := receive.NewQueueReceiver("amqp://guest:guest@localhost:8080/", receive.Queue{Name: "hello", Durable: false, AutoAck: true})
	if err != nil {
		log.Fatalf("Failed to create receiver: %v", err)
	}

	// Start the receiver
	go func() {
		if err := receiver.Receive(); err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}
	}()

	// Send a message every 5 seconds
	for {
		msg := fmt.Sprintf("Hello World! @ %v", time.Now())
		if err := publisher.Publish(context.Background(), []byte(msg)); err != nil {
			log.Fatalf("Failed to publish: %v", err)
		}
		//log.Printf("Sent: %s", msg)
		time.Sleep(10 * time.Second)
	}
}
