package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongoClient *mongo.Client) Models {
	client = mongoClient
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpatedAt:  time.Now(),
	})
	if err != nil {
		log.Println("Failed to insert log entry", err)
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{
		{Key: "created_at", Value: -1},
	})
	cursor, err := collection.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		log.Println("Failed to find log entries", err)
		return nil, err
	}
	var logs []LogEntry
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var logEntry LogEntry
		err := cursor.Decode(&logEntry)
		if err != nil {
			log.Println("Failed to decode log entry", err)
			return nil, err
		}
		logs = append(logs, logEntry)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to create object id from hex", err)
		return nil, err
	}
	var logEntry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&logEntry)
	if err != nil {
		log.Println("Failed to find log entry", err)
		return nil, err
	}
	return &logEntry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Failed to drop collection",
			err)
		return err
	}
	return nil
}

func (l *LogEntry) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Failed to create object id from hex", err)
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{
		"$set": bson.M{
			"name":       l.Name,
			"data":       l.Data,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		log.Println("Failed to update log entry", err)
		return err
	}
	return nil
}
