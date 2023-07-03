package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/atedesch1/mingle/models"
)

type MessageStore struct {
	*sqlx.DB
    dsn string
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

const getLatestMessagesQuery = `SELECT * FROM messages ORDER BY id DESC LIMIT $1`

func (s *MessageStore) GetLatestMessages(quantity uint) ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getLatestMessagesQuery, quantity); err != nil {
		return messages, fmt.Errorf("error getting messages: %w", err)
	}
	return messages, nil
}

const getMessagesRangeQuery = `SELECT * FROM messages WHERE id < $1 ORDER BY id DESC LIMIT $2`

func (s *MessageStore) GetMessagesRange(fromID uint64, quantity uint) ([]models.Message, error) {
	messages := []models.Message{}
	if err := s.Select(&messages, getMessagesRangeQuery, fromID, quantity); err != nil {
		return messages, fmt.Errorf("error getting messages: %w", err)
	}
	return messages, nil
}

const messageChannelName = `message_channel`

func (s *MessageStore) SubscribeToMessages(messageChannel chan<- []byte) {
    notifier, err := newNotifier(s.dsn, messageChannelName)
    if err != nil {
        log.Println(err.Error())
    }

    if err := notifier.fetch(messageChannel); err != nil {
        log.Println(err.Error())
    }
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
