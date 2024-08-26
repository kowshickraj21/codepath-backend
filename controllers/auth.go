package controllers

import (
	"database/sql"
	"fmt"
	"main/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func registerUser(db *sql.DB,userInfo models.User,provider string) {

	user,_:= GetUser(db,userInfo.Email)
	if(user == nil){
		CreateUser(db,userInfo,provider)
	}else{
		fmt.Println("Already Exists")
	}

}

func SignJWT(user *models.User)(string,error){
	claims := jwt.MapClaims{
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
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
		}
		return user, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}