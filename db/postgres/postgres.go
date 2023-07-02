package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresStorage(dataSourceName string) (*PostgresStorage, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &PostgresStorage{
		UserStore: &UserStore{
			DB: db,
		},
		MessageStore: &MessageStore{
			DB: db,
		},
	}, nil
}

type PostgresStorage struct {
    *UserStore
    *MessageStore
}