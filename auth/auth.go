package auth

import (
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetUserInfo(accessToken string) (*models.GoogleUser, error) {
	infoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	res, err := http.Get(fmt.Sprintf("%s?access_token=%s", infoEndpoint, accessToken))
	if err != nil {
		return nil,err
	}
	defer res.Body.Close()
	
	var userInfo models.GoogleUser
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo,nil
}

func SignJWT(user *models.GoogleUser)(string,error){
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