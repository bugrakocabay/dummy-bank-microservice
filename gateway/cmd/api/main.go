package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connectRabbitMQ()
	if err != nil {
		log.Println("rabbitconn err: ", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("listening for and consuming RabbitMQ messages...")

	app := Config{
		rabbit: rabbitConn,
	}

	log.Printf("Starting Gateway service on port: %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("amqp not ready yet...")
			counts++
		} else {
			log.Println("connected to rabbitmq!")
			connection = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}
