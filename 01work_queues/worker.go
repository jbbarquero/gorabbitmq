package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const queueName = "work_queue"

func main() {
	fmt.Println("Go works: receiving tasks from Rabbit")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			log.Printf("Received a message: %s", m.Body)
			dotCount := bytes.Count(m.Body, []byte("."))
			for i := 0; i < dotCount; i++ {
				fmt.Print(".")
				time.Sleep(1 * time.Second)
			}
			fmt.Println()
			log.Printf("\nDone")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	fmt.Println("Go works: receiving tasks from Rabbit ;)")
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", err, msg)
	}

}
