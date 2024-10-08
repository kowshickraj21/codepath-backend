package controllers

import (
	"database/sql"
	"main/models"

	"github.com/lib/pq"
)


func CreateUser(db *sql.DB, user models.User,provider string) (sql.Result,error) {
	user.Problems = []int64{}
	query := `INSERT INTO Users (name, email, picture, provider, problems) VALUES ($1, $2, $3, $4, $5)`
	res,err := db.Exec(query, user.Name, user.Email, user.Picture, provider, pq.Array(user.Problems))
	if err != nil {
		return nil,err
	}
	return res,nil
}

func GetUser(db *sql.DB, email string) (*models.User, error) {

	query := `SELECT name, email, picture, problems FROM Users WHERE email = $1`
	row := db.QueryRow(query,email)

	var user models.User
	err := row.Scan(&user.Name, &user.Email, &user.Picture, pq.Array(&user.Problems))

	if err != nil {
		return nil, err
	}

	return &user, nil
}
