package api

import (
	"authencation-service/src/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	WebPort int
	DB      *sql.DB
	Models  data.Models
}

func (app *Config) Serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.WebPort),
		Handler: app.Router(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("authencation service failed to start: %v", err)
	}
}
