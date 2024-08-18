package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetUserInfo(db *sql.DB,accessToken string) (string, error) {
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

	user,_:= GetUser(db,userInfo.Email)
	if(user == nil){
		CreateUser(db,userInfo)
	}


	token,err := SignJWT(&userInfo);
	if(err != nil) {
		return "",err;
	}
	return token,nil;
}

func SignJWT(user *models.User)(string,error){
	claims := jwt.MapClaims{
		"sub": user.Id,
		"name": user.Name,
		"email": user.Email,
		"iss": "oauth-app-golang",
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signed,err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "",err;
	}
	return signed,nil
}

func ParseJWT(token string)(*models.User,error){
	JWT,err := jwt.Parse(token,func(tok *jwt.Token)(interface{}, error){
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := JWT.Claims.(jwt.MapClaims); ok && JWT.Valid {

		user := &models.User{
			Id:    claims["sub"].(string),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
		}
		return user, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}