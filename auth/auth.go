package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GetUserInfo(accesstoken string) (map[string]interface{}, error) {
	infoEndpoint := "https://googleapis.com/oauth2/v2/userinfo"
	res,err := http.Get(fmt.Sprintf("%s?access_token=%s",infoEndpoint,accesstoken))
	
	if err != nil {
		return nil,err
	}
	defer res.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	
	return userInfo,nil
}

func SignJWT(user map[string]interface{})(string,error){
	claims := jwt.MapClaims{
		"sub": user["id"],
		"name": user["name"],
		"email": user["email"],
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