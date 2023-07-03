package db

import (
	"fmt"

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

	return &PostgresStorage{
		UserStore: &UserStore{
			DB: db,
            dsn: dsn,
		},
		MessageStore: &MessageStore{
			DB: db,
            dsn: dsn,
		},
	}, nil
}

type PostgresStorage struct {
    *UserStore
    *MessageStore
}
