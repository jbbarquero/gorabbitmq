package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

const queueName = "work_queue"

func main() {
	fmt.Println("Go send tasks to Rabbit")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, //name
		true,      //durable
		false,     //delete when unused
		false,     //exclusive
		false,     //no wait
		nil,       //arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := bodyFrom(os.Args)

	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  //mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	fmt.Println("Go send tasks to Rabbit ;)")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || args[1] == "" {
		s = "Task from go"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}
