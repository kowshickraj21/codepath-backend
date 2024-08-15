package models

type User struct {
	Name     string
	Email    string
	Picture  string
	Problems []int64
}

type GoogleUser struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
	Name    string `json:"name"`
}
