package main

import (
	"bytes"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/thebigyovadiaz/rabbitmq-work-queues/src/util"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	util.LogFailOnError(err, "Failed to connect to RabbitMQ")
	util.LogSuccessful("Connected to RabbitMQ")
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

	err = ch.Qos(
		1,
		0,
		false,
	)
	util.LogFailOnError(err, "Failed to set QoS")

	messages, err := ch.Consume(
		qD.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	util.LogFailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range messages {
			util.LogSuccessful(fmt.Sprintf("Received a message: %s", d.Body))
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			util.LogSuccessful("Done")
			d.Ack(false)
		}
	}()

	util.LogSuccessful(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
