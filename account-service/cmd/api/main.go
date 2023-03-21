package main

import (
	"database/sql"
	"fmt"
	"log"

	db "github.com/bugrakocabay/dummy-bank-microservice/account-service/db/sqlc"
	_ "github.com/lib/pq"
)

const webPort = "80"

func main() {
	log.Printf("Starting Account service on port: %s", webPort)

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error with loading env: ", err)
	}

	conn, err := sql.Open("postgres", config.AccountDbConnString)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := NewServer(store)

	address := fmt.Sprintf(":%s", webPort)
	if err = server.Start(address); err != nil {
		log.Fatal("Server start failed:", err)
	}
}
