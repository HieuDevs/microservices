package api

import (
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	WebPort int
	Rabbit  *amqp.Connection
}

func (app *Config) Serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.WebPort),
		Handler: app.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}

}
