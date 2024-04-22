package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
}

func TestMain(m *testing.M) {
	var err error

	testDB, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
