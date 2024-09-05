package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func HandleGithubUser(db *sql.DB, code string) (string, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	tokenURL := "https://github.com/login/oauth/access_token"

	data := make(url.Values)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching access token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling token response: %v", err)
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("invalid access token response")
	}

	userInfo, err := fetchGithubUserInfo(accessToken)
	if err != nil {
		return "", fmt.Errorf("error fetching user info: %v", err)
	}

	registerUser(db, userInfo,"github");

	token, err := SignJWT(&userInfo)
	if err != nil {
		return "", fmt.Errorf("error signing JWT: %v", err)
	}

	return token, nil
}

func fetchGithubUserInfo(accessToken string) (models.User, error) {
	userInfoURL := "https://api.github.com/user"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return models.User{}, fmt.Errorf("error fetching user info: %v", err)
	}
	defer res.Body.Close()

	var userInfo models.User
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return models.User{}, fmt.Errorf("error decoding user info: %v", err)
	}
	fmt.Println("picture",userInfo.Picture)

	return userInfo, nil
}
