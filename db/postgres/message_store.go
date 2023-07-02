package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/atedesch1/mingle/models"
)

type MessageStore struct {
	*sqlx.DB
}

const getMessageQuery = `SELECT * FROM messages WHERE id = $1`

func (s *MessageStore) GetMessage(id uint64) (models.Message, error) {
	message := models.Message{}
	if err := s.Get(&message, getMessageQuery, id); err != nil {
		return message, fmt.Errorf("error getting message: %w", err)
	}
	return message, nil
}

const getMessagesQuery = `SELECT * FROM messages`

func (s *MessageStore) GetMessages() ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getMessagesQuery); err != nil {
		return messages, fmt.Errorf("error getting messages: %w", err)
	}
	return messages, nil
}

const getLatestMessagesQuery = `SELECT * FROM messages WHERE ORDER BY id DESC LIMIT $1`

func (s *MessageStore) GetLatestMessages(quantity uint) ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getLatestMessagesQuery, quantity); err != nil {
		return messages, fmt.Errorf("error getting messages: %w", err)
	}
	return messages, nil
}

const getMessagesRangeQuery = `SELECT * FROM messages WHERE id < $1 ORDER BY id DESC LIMIT $2`

func (s *MessageStore) GetMessagesRange(begin uint, quantity uint) ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getMessagesRangeQuery, begin, quantity); err != nil {
		return messages, fmt.Errorf("error getting messages: %w", err)
	}
	return messages, nil
}

const createMessageQuery = `INSERT INTO messages (user_id, content) VALUES ($1, $2) RETURNING *`

func (s *MessageStore) CreateMessage(params models.MessageCreateParams) (models.Message, error) {
	message := models.Message{}
	if err := s.Get(&message, createMessageQuery, params.UserID, params.Content); err != nil {
		return message, fmt.Errorf("error creating message: %w", err)
	}
	return message, nil
}

const deleteMessageQuery = `DELETE FROM messages WHERE id = $1`

func (s *MessageStore) DeleteMessage(id uint64) error {
    message := models.Message{}
	if err := s.Get(&message, getMessageQuery, id); err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}
	if _, err := s.Exec(deleteMessageQuery, id); err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}
	return nil
}
