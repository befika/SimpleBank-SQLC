package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *pgxpool.Pool

const (
	// dbDriver = "pgx"
	// dbSource = "postgresql://postgres:postgres@localhost:5432/sample_bank?sslmode=disable&pool_max_conns=10000"
	dbSource = "postgresql://root:@localhost:26257/sample_bank?sslmode=disable&pool_max_conns=1000"
)

func TestMain(m *testing.M) {
	var err error
	testDb, err = pgxpool.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}
