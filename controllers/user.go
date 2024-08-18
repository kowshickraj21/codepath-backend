package controllers

import (
	"database/sql"
	"main/models"

	"github.com/lib/pq"
)


func CreateUser(db *sql.DB, user models.User) (sql.Result,error) {
	user.Problems = []int64{}
	query := `INSERT INTO Users (id,name, email, picture, problems) VALUES ($1, $2, $3, $4, $5)`
	res,err := db.Exec(query,user.Id, user.Name, user.Email, user.Picture, pq.Array(user.Problems))
	if err != nil {
		return nil,err
	}
	return res,nil
}

func GetUser(db *sql.DB, token string) (*models.User, error) {
	authUser,err := ParseJWT(token);
	if(err != nil){
		return nil,err
	}
	query := `SELECT name, email, picture, problems FROM Users WHERE email = $1`
	row := db.QueryRow(query, authUser.Email)

	var user models.User
	err = row.Scan(&user.Name, &user.Email, &user.Picture, pq.Array(&user.Problems))

	if err != nil {
		return nil, err
	}

	return &user, nil
}
