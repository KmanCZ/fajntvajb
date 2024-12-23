package repository

import (
	"github.com/jmoiron/sqlx"
)

type Users struct {
	db *sqlx.DB
}

type User struct {
	ID          int
	Username    string
	DisplayName string
	Password    string
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{db: db}
}

func (users *Users) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := users.db.Get(user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (users *Users) CreateUser(username, displayName, password string) (*User, error) {
	user := &User{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
	}
	_, err := users.db.NamedExec("INSERT INTO users (username, display_name, password) VALUES (:username, :display_name, :password)", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
