package db

import "github.com/atedesch1/mingle/models"

type UserStore interface {
	GetUser(id uint64) (models.User, error)
	GetUsers() ([]models.User, error)
	CreateUser(params models.UserCreateParams) (models.User, error)
	DeleteUser(id uint64) error
}

type MessageStore interface {
	GetMessage(id uint64) (models.Message, error)
	GetMessages() ([]models.Message, error)
	GetLatestMessages(quantity uint) ([]models.Message, error)
	GetMessagesRange(begin uint, quantity uint) ([]models.Message, error)
	CreateMessage(params models.MessageCreateParams) (models.Message, error)
	DeleteMessage(id uint64) error
}

type Storage interface {
	UserStore
	MessageStore
}
