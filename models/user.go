package models

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// User represents a user in the database
type User struct {
	Username string
	Name     string
	Email    string
	Picture  string
	Problems []int
}

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, user *User) error {
	query := `INSERT INTO Users (username, name, email, picture, problems) VALUES ($1, $2, $3, $4, $5) RETURNING username`
	err := db.QueryRow(query, user.Username, user.Name, user.Email, user.Picture, pq.Array(user.Problems)).Scan(&user.Username)
	if err != nil {
		return err
	}
	return nil
}

// GetUser retrieves a user from the database by username
func GetUser(db *sql.DB, username string) (*User, error) {
	query := `SELECT * FROM Users WHERE username = $1`
	row := db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.Username, &user.Name, &user.Email, &user.Picture, pq.Array(&user.Problems))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
