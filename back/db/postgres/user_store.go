package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/atedesch1/mingle/models"
)

type UserStore struct {
	*sqlx.DB
    dsn string
}

const getUserQuery = `SELECT * FROM users WHERE id = $1`

func (s *UserStore) GetUser(id uint64) (models.User, error) {
	user := models.User{}
	if err := s.Get(&user, getUserQuery, id); err != nil {
		return user, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

const getUsersQuery = `SELECT * FROM users`

func (s *UserStore) GetUsers() ([]models.User, error) {
	users := []models.User{}
	if err := s.Select(&users, getUsersQuery); err != nil {
		return users, fmt.Errorf("error getting users: %w", err)
	}
	return users, nil
}

const createUserQuery = `INSERT INTO users (name) VALUES ($1) RETURNING *`

func (s *UserStore) CreateUser(params models.UserCreateParams) (models.User, error) {
	user := models.User{}
	if err := s.Get(&user, createUserQuery, params.Name); err != nil {
		return user, fmt.Errorf("error creating user: %w", err)
	}
	return user, nil
}

const deleteUserQuery = `DELETE FROM users WHERE id = $1`

func (s *UserStore) DeleteUser(id uint64) error {
    user := models.User{}
	if err := s.Get(&user, getUserQuery, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	if _, err := s.Exec(deleteUserQuery, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
