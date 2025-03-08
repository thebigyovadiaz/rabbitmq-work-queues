package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thebigyovadiaz/rabbitmq-work-queues/src/util"
	"os"
	"strings"
	"time"
)

func bodyForm(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.LogFailOnError(err, "Failed to connect to RabbitMQ")
	util.LogSuccessful("RabbitMQ connect successfully")
	defer conn.Close()

	ch, err := conn.Channel()
	util.LogFailOnError(err, "Failed to open a channel")
	util.LogSuccessful("Channel open successfully")
	defer ch.Close()

	qD, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	util.LogFailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyForm(os.Args)
	err = ch.PublishWithContext(
		ctx,
		"",
		qD.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)

	util.LogFailOnError(err, "Failed to publish a message")
	util.LogSuccessful(fmt.Sprintf(" [x] Sent %s", body))
}
