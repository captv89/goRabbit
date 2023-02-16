package send

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// QueuePublisher is a wrapper around the AMQP connection and channel
type QueuePublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

// NewQueuePublisher creates a new QueuePublisher instance
func NewQueuePublisher(url string, queueName string) *QueuePublisher {
	conn, err := amqp.Dial(url)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}

	return &QueuePublisher{
		connection: conn,
		channel:    ch,
		queue:      q,
	}
}

// Publish sends a message to the queue
func (p *QueuePublisher) Publish(ctx context.Context, msg []byte) error {
	err := p.channel.PublishWithContext(
		ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		return err
	}

	log.Printf("[->] Sent message: %s", msg)
	return nil
}

// Close closes the AMQP connection
func (p *QueuePublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return err
		}
	}

	if p.connection != nil {
		if err := p.connection.Close(); err != nil {
			return err
		}
	}

	return nil
}
