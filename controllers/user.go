package controllers

import (
	"database/sql"
	"main/models"

	"github.com/lib/pq"
)


func CreateUser(db *sql.DB, user models.User) (error) {
	query := `INSERT INTO Users (username, name, email, picture, problems) VALUES ($1, $2, $3, $4, $5) RETURNING username`
	err := db.QueryRow(query, user.Username, user.Name, user.Email, user.Picture, pq.Array(user.Problems)).Scan(&user.Username)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(db *sql.DB, username string) (*models.User, error) {
	query := `SELECT * FROM Users WHERE username = $1`
	row := db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.Username, &user.Name, &user.Email, &user.Picture, pq.Array(&user.Problems))
	if err != nil {
		return nil, err
	}

	return &user, nil
}
