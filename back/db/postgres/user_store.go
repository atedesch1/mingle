package db

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	. "github.com/atedesch1/mingle/errors"
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
		if errors.Is(err, sql.ErrNoRows) {
			return user, &NotFoundError{Ty: "getting user", Err: err}
		}

		return user, &InternalError{Ty: "getting user", Err: err}
	}
	return user, nil
}

const getUsersQuery = `SELECT * FROM users`

func (s *UserStore) GetUsers() ([]models.User, error) {
	users := []models.User{}
	if err := s.Select(&users, getUsersQuery); err != nil {
		return users, &InternalError{Ty: "getting users", Err: err}
	}
	return users, nil
}

const createUserQuery = `INSERT INTO users (name) VALUES ($1) RETURNING *`

func (s *UserStore) CreateUser(params models.UserCreateParams) (models.User, error) {
	user := models.User{}
	if err := s.Get(&user, createUserQuery, params.Name); err != nil {
		return user, &InternalError{Ty: "creating user", Err: err}
	}
	return user, nil
}

const deleteUserQuery = `DELETE FROM users WHERE id = $1`

func (s *UserStore) DeleteUser(id uint64) error {
	if _, err := s.Exec(deleteUserQuery, id); err != nil {
		return &InternalError{Ty: "deleting user", Err: err}
	}
	return nil
}
