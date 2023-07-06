package db

import (
	postgres "github.com/atedesch1/mingle/db/postgres"
)

const dsn = "postgresql://root:pass@0.0.0.0:5432/mingle?sslmode=disable"

func NewStorage() (Storage, error) {
	return postgres.NewPostgresStorage(dsn)
}
