package controllers

import (
	"database/sql"
	"main/models"
)

// UserController handles user-related operations
type UserController struct {
	DB *sql.DB
}

// NewUserController creates a new UserController
func NewUserController(db *sql.DB) *UserController {
	return &UserController{DB: db}
}

// CreateUser creates a new user
func (uc *UserController) CreateUser(username, name, email, picture string, problems []int) (*models.User, error) {
	user := &models.User{
		Username: username,
		Name:     name,
		Email:    email,
		Picture:  picture,
		Problems: problems,
	}
	err := models.CreateUser(uc.DB, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser retrieves a user by username
func (uc *UserController) GetUser(username string) (*models.User, error) {
	user, err := models.GetUser(uc.DB, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
