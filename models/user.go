package models

import "time"

type User struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type UserCreateParams struct {
    Name string
}
