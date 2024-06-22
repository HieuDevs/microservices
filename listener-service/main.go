package main

import (
	"listener-service/event"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")
	// create consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Panicf("Failed to create consumer: %s", err)
	}
	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Printf("Failed to consume messages: %s", err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Printf("RabbitMQ not yet ready...: %s", err)
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = conn
			break
		}
		if counts > 5 {
			log.Printf("RabbitMQ is ready after %d retries", counts)
			return nil, err
		}
		backOff = backOff * 2
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
