package main

import (
	"fmt"
	"log"
	"my-broker/src/api"
	"net/http"
)

func main() {
	app := api.Config{
		WebPort: 80,
	}
	log.Println("Starting broker-service on port", app.WebPort)

	// Start the web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.WebPort),
		Handler: app.Routes(),
	}

	// Start the server
	if err := srv.ListenAndServe(); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}

}
