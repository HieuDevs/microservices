package main

import (
	"authencation-service/src/api"
	"authencation-service/src/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.Println("Starting authencation-service")
	// Open a database connection
	db := connectToDB()
	// Set up the server
	app := api.Config{
		WebPort: 80,
		DB:      db,
		Models:  data.New(db),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.WebPort),
		Handler: app.Router(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("authencation service failed to start: %v", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil
}

const maxConnectionAttempts = 10
const retryInterval = 1

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	retrCount := 0
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("error connecting to database: %v", err)
			retrCount++
		} else {
			log.Println("Connected to database")
			return connection
		}
		if retrCount >= maxConnectionAttempts {
			log.Panicf("failed to connect to database after %d attempts", maxConnectionAttempts)
			return nil
		}
		time.Sleep(retryInterval * time.Second)
		continue
	}

}
