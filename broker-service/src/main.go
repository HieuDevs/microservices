package main

import (
	"log"
	"my-broker/src/api"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//try to connect to RabbitMQ
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	app := api.Config{
		WebPort: 80,
		Rabbit:  conn,
	}
	log.Println("Starting broker-service on port", app.WebPort)
	app.Serve()
}

func connectToRabbitMQ() (*amqp.Connection, error) {
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
