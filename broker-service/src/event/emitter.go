package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type emitter struct {
	conn *amqp.Connection
}

func (e *emitter) setup() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare exchange
	return declareExchange(ch, "logs_topic")
}

func (e *emitter) Push(event string, severity string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	log.Println("Pushing to channel", event, severity)
	return ch.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
}

func NewEmitter(conn *amqp.Connection) (*emitter, error) {
	emitter := &emitter{
		conn: conn,
	}
	err := emitter.setup()
	if err != nil {
		return nil, err
	}
	return emitter, nil
}
