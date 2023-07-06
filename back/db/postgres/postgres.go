package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	m.Up()

	return &PostgresStorage{
		UserStore: &UserStore{
			DB:  db,
			dsn: dsn,
		},
		MessageStore: &MessageStore{
			DB:  db,
			dsn: dsn,
		},
	}, nil
}

type PostgresStorage struct {
	*UserStore
	*MessageStore
}
