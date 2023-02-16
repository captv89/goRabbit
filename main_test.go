package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"goRabbit/receive"
	"goRabbit/send"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	url := "amqp://guest:guest@localhost:8080/"
	// Create a RabbitMQ connection
	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			t.Fatalf("Failed to close connection: %s", err)
		}
	}(conn)

	// Set up a queue and bind it to the exchange
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open a channel: %s", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			t.Fatalf("Failed to close channel: %s", err)
		}
	}(ch)

	_, err = ch.QueueDeclare(
		"test_queue", // queue name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		t.Fatalf("Failed to declare a queue: %s", err)
	}

	// Set up a subscriber to listen on the test queue
	recv := receive.NewQueueReceiver(url, receive.Queue{Name: "test_queue", AutoAck: true})
	if err != nil {
		t.Fatalf("Failed to create receiver: %s", err)
	}

	go func() {
		recv.Receive()
	}()

	// Set up a publisher to publish messages to the test queue
	pub := send.NewQueuePublisher(url, "test_queue")

	// Publish a test message every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := pub.Publish(context.Background(), []byte("Hello World!"))
				if err != nil {
					t.Errorf("Failed to publish message: %s", err)
					return
				}
			}
		}
	}()

	// Wait for messages to be received for 30 seconds
	time.Sleep(30 * time.Second)
}
