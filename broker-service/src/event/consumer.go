package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (*consumer, error) {
	consumer := &consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (c *consumer) setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare exchange
	return declareExchange(ch, "logs_topic")
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, e := declareRandomQueue(ch)
	if e != nil {
		return e
	}
	for _, topic := range topics {
		if err := ch.QueueBind(
			q.Name,       // queue name
			topic,        // routing key
			"logs_topic", // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}
	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range messages {
			payload := Payload{}
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()
	fmt.Printf(" [*] Waiting for message [Exchange, Queue] [logs_topic,%s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// do something
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// do something

	default:
		// do something
	}
}

func logEvent(data Payload) error {
	// create json we'll send to logging service
	jsonData, _ := json.MarshalIndent(data, "", " \t")
	// call the service
	request, err := http.NewRequest(
		"POST",
		"http://logger-service/log",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
