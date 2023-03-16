package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bugrakocabay/dummy-bank-microservice/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Printf("Starting Logging service on port: %s", webPort)

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error with loading env: ", err)
	}

	mongoClient, err := connectToMongo(config.MongoURI, config.Username, config.Password)
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo(mongoURI, username, password string) (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting mongo:", err)
		return nil, err
	}

	return c, nil
}
