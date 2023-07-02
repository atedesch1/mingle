package db

import "github.com/atedesch1/mingle/models"

type UserStore interface {
	GetUser(id uint64) (models.User, error)
	GetUsers(id uint64) ([]models.User, error)
	CreateUser(params models.UserCreateParams) (models.User, error)
	DeleteUser(id uint64) error
}

type MessageStore interface {
	GetMessage(id uint64) (models.Message, error)
	GetMessages() ([]models.Message, error)
	GetMessagesRange(begin uint, end uint) ([]models.Message, error)
	CreateMessage(params models.MessageCreateParams) (models.Message, error)
	DeleteMessage(id uint64) error
}

type Storage interface {
	UserStore
	MessageStore
}
