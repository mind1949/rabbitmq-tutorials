package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	sender()
	select {}
}

func sender() {
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

	go func() {
		defer conn.Close()
		defer ch.Close()

		for count := 0; ; count++ {
			body := fmt.Sprintf("%d(Hello World!)", count)
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			log.Printf("queue(%s) <- %s\n", q.Name, body)

			time.Sleep(2 * time.Second)
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
