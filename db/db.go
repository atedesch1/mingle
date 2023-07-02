package db

import postgres "github.com/atedesch1/mingle/db/postgres"

func NewStorage() (Storage, error) {
	return postgres.NewPostgresStorage("postgresql://root:pass@0.0.0.0:5432/postgres?sslmode=disable")
}
