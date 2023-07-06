package db

import (
	"os"

	postgres "github.com/atedesch1/mingle/db/postgres"
)

func NewStorage() (Storage, error) {
	dsn := os.Getenv("DATABASE_URL")
	return postgres.NewPostgresStorage(dsn)
}
