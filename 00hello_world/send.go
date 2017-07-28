package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Go send to RabbitMQ")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello_from_java",
		false, //durable
		false, //delete when unused
		false, //exclusive
		false, //no wait
		nil,   //arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello from go"
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  //mandatory
		false,  //inmediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	fmt.Println("Go sent to RabbitMQ ;) ")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}
