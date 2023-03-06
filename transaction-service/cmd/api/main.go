package main

import (
	"database/sql"
	"fmt"
	db "github.com/bugrakocabay/dummy-bank-microservice/transaction-service/db/sqlc"
	"log"
)

const webPort = "80"

func main() {
	log.Printf("Starting Transaction service on port: %s", webPort)

	// TODO: Use env variables
	conn, err := sql.Open("postgres", "postgresql://postgres:postgres@account_db_postgres:5433/transactions?sslmode=disable")
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
