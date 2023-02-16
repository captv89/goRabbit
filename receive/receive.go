package receive

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

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

func NewQueueReceiver(url string, q Queue) *QueueReceiver {
	conn, err := amqp.Dial(url)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
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
		failOnError(err, "Failed to declare a queue")
	}

	return &QueueReceiver{
		conn: conn,
		ch:   ch,
		q:    &amqp.Queue{Name: q.Name},
	}
}

func (r *QueueReceiver) Receive() {
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
		failOnError(err, "Failed to register a consumer")
	}

	go func() {
		for msg := range msgs {
			log.Printf("[<-] Received message: %s", string(msg.Body))
		}
	}()
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
