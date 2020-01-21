package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	receiver()
	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func receiver() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")

	go func() {
		defer conn.Close()
		defer ch.Close()

		for d := range msgs {
			log.Printf("queue(%s) -> %s\n", q.Name, d.Body)
		}
	}()
}
