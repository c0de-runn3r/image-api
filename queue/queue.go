package queue

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewRabbitMQ(addr string) *RabbitMQ {
	conn, err := amqp.Dial(addr)
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	queue, err := channel.QueueDeclare(
		"imageOptimization",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	return &RabbitMQ{
		Channel: channel,
		Queue:   queue,
	}
}

func (r *RabbitMQ) SendID(body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.Channel.PublishWithContext(ctx,
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [RabbitMQ] Sent %s\n", body)
}

func (r *RabbitMQ) RecieveID() (bool, string) {
	msg, ok, err := r.Channel.Get(
		r.Queue.Name,
		true,
	)
	failOnError(err, "Failed to register a consumer")
	if ok {
		log.Printf("Received an ID: %s", msg.Body)
		return true, string(msg.Body)
	}
	return false, ""
}
