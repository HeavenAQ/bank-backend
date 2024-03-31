package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/simple-bank-db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	testQueries = New(dbpool)
	m.Run()
}
