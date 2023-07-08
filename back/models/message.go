package models

import "time"

type Message struct {
	ID        uint64    `json:"id"        db:"id"`
	UserID    uint64    `json:"userId"    db:"user_id"`
	Content   string    `json:"content"   db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type MessageExpandedUser struct {
	ID        uint64    `json:"id"        db:"id"`
	Content   string    `json:"content"   db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	User `json:"user" db:"user"`
}

type MessageCreateParams struct {
	UserID  uint64
	Content string
}
