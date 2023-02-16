package receive

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Queue struct {
	Name    string
	Durable bool
	AutoAck bool
}

type QueueReceiver struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    *amqp.Queue
}

func NewQueueReceiver(url string, q Queue) (*QueueReceiver, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		q.Name,
		q.Durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &QueueReceiver{
		conn: conn,
		ch:   ch,
		q:    &amqp.Queue{Name: q.Name},
	}, nil
}

func (r *QueueReceiver) Receive() error {
	msgs, err := r.ch.Consume(
		r.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", string(msg.Body))
		}
	}()

	return nil
}

func (r *QueueReceiver) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}
	if err := r.conn.Close(); err != nil {
		return err
	}
	return nil
}
