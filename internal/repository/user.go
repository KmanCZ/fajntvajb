package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	db *sqlx.DB
}

type User struct {
	ID          int    `db:"id"`
	Username    string `db:"username"`
	DisplayName string `db:"display_name"`
	Password    string `db:"password"`
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{db: db}
}

func (users *Users) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := users.db.Get(user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, nil
}

func (users *Users) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := users.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, nil
}

func (users *Users) CreateUser(username, displayName, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username:    username,
		DisplayName: displayName,
		Password:    string(hashedPassword),
	}
	_, err = users.db.NamedExec("INSERT INTO users (username, display_name, password) VALUES (:username, :display_name, :password)", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (users *Users) UpdateDisplayName(id int, newName string) error {
	_, err := users.db.Exec("UPDATE users SET display_name = $1 WHERE id = $2", newName, id)
	if err != nil {
		return err
	}
	return nil
}

func (users *Users) UpdatePassword(id int, newPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	_, err = users.db.Exec("UPDATE users SET password = $1 WHERE id = $2", string(hashedPassword), id)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (users *Users) DeleteUser(id int) error {
	_, err := users.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
