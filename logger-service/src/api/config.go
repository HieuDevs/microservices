package api

import (
	"log"
	"logger-service/src/data"
	"net/http"
)

type Config struct {
	WebPort  string
	RPCPort  string
	MongoURL string
	GRPCPort string
	Models   data.Models
}

func (app *Config) Serve() {
	srv := &http.Server{
		Addr:    ":" + app.WebPort,
		Handler: app.Routes(),
	}
	log.Println("Starting logger-service on port", app.WebPort)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
