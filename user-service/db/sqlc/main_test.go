package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5433/users_test?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testQueries = New(conn)
	cleanDB(testQueries)

	os.Exit(m.Run())
}

func cleanDB(queries *Queries) {
	query1 := "DELETE FROM users;"
	_, err := queries.db.QueryContext(context.Background(), query1)
	if err != nil {
		log.Printf("error cleaning accounts table: %v", err)
	}
}
