package fuzz

import (
	"context"
	"goRabbit/receive"
	"goRabbit/send"
	"testing"
	"time"
)

// FuzzMain is the entry point for the fuzzer.
func FuzzMain(f *testing.F) {
	// create a publisher
	publisher := send.NewQueuePublisher("amqp://guest:guest@localhost:8080/", "fuzz")

	// create a receiver
	receiver := receive.NewQueueReceiver("amqp://guest:guest@localhost:8080/", receive.Queue{Name: "fuzz", Durable: false, AutoAck: true})

	// Start the receiver
	go func() {
		receiver.Receive()
	}()

	// Send a message every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f.Fuzz(func(t *testing.T, data []byte) {
					err := publisher.Publish(context.Background(), data)
					if err != nil {
						t.Fatalf("Failed to publish message: %s", err)
					}
				})
			}
		}
	}()

	// Wait for messages to be received for 30 seconds
	time.Sleep(30 * time.Second)
}
