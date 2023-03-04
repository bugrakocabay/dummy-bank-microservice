package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {
	log.Printf("Starting Account service on port: %s", webPort)

	// TODO Connect to db

	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
