package models

import "time"

type User struct {
	ID        uint64    `json:"id"        db:"id"`
	Name      string    `json:"name"      db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type UserCreateParams struct {
	Name string
}
