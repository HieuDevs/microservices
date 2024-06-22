package main

import (
	"context"
	"log"
	"logger-service/src/api"
	"logger-service/src/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	app := api.Config{
		WebPort:  "80",
		RPCPort:  "5001",
		MongoURL: "mongodb://mongo:27017",
		GRPCPort: "50001",
	}
	mongoClient, err := connectToMongo(app.MongoURL)
	if err != nil {
		log.Panic(err)
	}
	if mongoClient == nil {
		log.Panic("Crashed while connecting to MongoDB")
	}
	client = mongoClient
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Panic(err)
		}
	}()
	app.Models = data.New(client)
	app.Serve()
}

const maxRetries = 5
const delay = 1

func connectToMongo(mongoUrl string) (*mongo.Client, error) {
	var retryCount int = 0
	for {
		client, err := mongo.Connect(
			context.TODO(),
			options.Client().
				ApplyURI(mongoUrl).
				SetAuth(options.Credential{
					Username: "admin",
					Password: "password",
				}),
		)
		if err != nil {
			if retryCount < maxRetries {
				retryCount++
				log.Println("Failed to connect to MongoDB", err)
			} else {
				log.Println("Failed to connect to MongoDB after 5 retries")
				return nil, err
			}
		} else {
			log.Println("Connected to MongoDB")
			return client, nil
		}
		time.Sleep(delay * time.Second)
	}
}
