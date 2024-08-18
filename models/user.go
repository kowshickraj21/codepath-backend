package models

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Name     string `json:"name"`
	Problems []int64
}

type JWT struct {
	Token string `json:"token"`
}