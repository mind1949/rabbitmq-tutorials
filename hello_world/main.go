package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	sender()
	receiver()

	select {}
}

func sender() {
	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

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

		select {}
	}()
}

func receiver() {
	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

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
			for d := range msgs {
				log.Printf("queue(%s) -> %s\n", q.Name, d.Body)
			}
		}()
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
