package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func Connect(url string) (*amqp091.Connection, *amqp091.Channel) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		"otp_email_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	return conn, ch
}
