package models

import "time"

type Message struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type MessageCreateParams struct {
	UserID  uint64
	Content string
}
