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
var testStore *Store

func TestMain(m *testing.M) {
	ctx := context.Background()
	testDBPool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer testDBPool.Close()
	testQueries = New(testDBPool)
	testStore = NewStore(testDBPool)
	m.Run()
}
