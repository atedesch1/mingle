package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"

	. "github.com/atedesch1/mingle/errors"
	"github.com/atedesch1/mingle/models"
)

type MessageStore struct {
	*sqlx.DB
	dsn string
}

const getMessageQuery = `
SELECT * 
FROM messages 
WHERE id = $1`

func (s *MessageStore) GetMessage(id uint64) (models.Message, error) {
	message := models.Message{}
	if err := s.Get(&message, getMessageQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return message, &NotFoundError{Ty: "getting message", Err: err}
		}

		return message, &InternalError{Ty: "getting message", Err: err}
	}
	return message, nil
}

const getMessagesQuery = `
SELECT * 
FROM messages`

func (s *MessageStore) GetMessages() ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getMessagesQuery); err != nil {
		return messages, &InternalError{Ty: "getting messages", Err: err}
	}
	return messages, nil
}

const getMessagesExpandedUserQuery = `
SELECT 
    m.id, m.content, m.created_at,
    u.id "user.id", u.name "user.name", u.created_at "user.created_at"
FROM messages as m
JOIN users as u
ON m.user_id = u.id`

func (s *MessageStore) GetMessagesExpandedUser() ([]models.MessageExpandedUser, error) {
	messages := []models.MessageExpandedUser{}
	if err := s.Select(&messages, getMessagesExpandedUserQuery); err != nil {
		return messages, &InternalError{Ty: "getting messages", Err: err}
	}
	return messages, nil
}

const getLatestMessagesQuery = `
SELECT * 
FROM messages 
ORDER BY id DESC LIMIT $1`

func (s *MessageStore) GetLatestMessages(quantity uint) ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getLatestMessagesQuery, quantity); err != nil {
		return messages, &InternalError{Ty: "getting latest messages", Err: err}
	}
	return messages, nil
}

const getLatestMessagesExpandedUserQuery = `
SELECT 
    m.id, m.content, m.created_at,
    u.id "user.id", u.name "user.name", u.created_at "user.created_at"
FROM messages as m
JOIN users as u
ON m.user_id = u.id
ORDER BY m.id DESC LIMIT $1`

func (s *MessageStore) GetLatestMessagesExpandedUser(quantity uint) ([]models.MessageExpandedUser, error) {
	messages := []models.MessageExpandedUser{}
	if err := s.Select(&messages, getLatestMessagesExpandedUserQuery, quantity); err != nil {
		return messages, &InternalError{Ty: "getting latest messages", Err: err}
	}
	return messages, nil
}

const getMessagesRangeQuery = `
SELECT * 
FROM messages 
WHERE id < $1 
ORDER BY id DESC LIMIT $2`

func (s *MessageStore) GetMessagesRange(fromID uint64, quantity uint) ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getMessagesRangeQuery, fromID, quantity); err != nil {
		return messages, &InternalError{Ty: "getting messages range", Err: err}
	}
	return messages, nil
}

const getMessagesRangeExpandedUserQuery = `
SELECT 
    m.id, m.content, m.created_at,
    u.id "user.id", u.name "user.name", u.created_at "user.created_at"
FROM messages as m
JOIN users as u
ON m.user_id = u.id
WHERE m.id < $1 ORDER BY m.id DESC LIMIT $2`

func (s *MessageStore) GetMessagesRangeExpandedUser(fromID uint64, quantity uint) ([]models.MessageExpandedUser, error) {
	messages := []models.MessageExpandedUser{}
	if err := s.Select(&messages, getMessagesRangeExpandedUserQuery, fromID, quantity); err != nil {
		return messages, &InternalError{Ty: "getting messages range", Err: err}
	}
	return messages, nil
}

const messageChannelName = `message_channel`

func (s *MessageStore) SubscribeToMessages(messageChannel chan<- []byte, unsubscribe chan struct{}) {
	notifier, err := newNotifier(s.dsn, messageChannelName)
	if err != nil {
		log.Println("notifier error:", err.Error())
	}

	if err := notifier.fetch(messageChannel, unsubscribe); err != nil {
		log.Println("notifier error:", err.Error())
		unsubscribe <- struct{}{}
	}
}

const createMessageQuery = `
INSERT INTO messages (user_id, content) 
VALUES ($1, $2) 
RETURNING *`

func (s *MessageStore) CreateMessage(params models.MessageCreateParams) (models.Message, error) {
	message := models.Message{}
	if err := s.Get(&message, createMessageQuery, params.UserID, params.Content); err != nil {
		return message, &InternalError{Ty: "creating message", Err: err}
	}
	return message, nil
}

const deleteMessageQuery = `
DELETE FROM messages 
WHERE id = $1`

func (s *MessageStore) DeleteMessage(id uint64) error {
	if _, err := s.Exec(deleteMessageQuery, id); err != nil {
		return &InternalError{Ty: "deleting message", Err: err}
	}
	return nil
}
