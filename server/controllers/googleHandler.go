package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
)
func HandleGoogleUser(db *sql.DB,accessToken string) (string, error) {
	infoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	res, err := http.Get(fmt.Sprintf("%s?access_token=%s", infoEndpoint, accessToken))
	if err != nil {
		return "",err
	}
	defer res.Body.Close()
	
	var userInfo models.User
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return "", err
	}
	
	registerUser(db,userInfo,"google")

	token,err := SignJWT(&userInfo);
	if(err != nil) {
		return "",err;
	}
	return token,nil;
}