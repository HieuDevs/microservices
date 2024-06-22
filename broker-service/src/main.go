package main

import (
	"log"
	"my-broker/src/api"
)

func main() {
	app := api.Config{
		WebPort: 80,
	}
	log.Println("Starting broker-service on port", app.WebPort)
	app.Serve()
}
