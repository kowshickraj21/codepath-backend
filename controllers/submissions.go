package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
	"os"
	"strings"
)


func CreateReq() (*models.RequestToken,error) {

	ApiKey := os.Getenv("JUDGE0_API_KEY")
	ApiHost := os.Getenv("JUDGE0_API_HOST")


	url := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=false&fields=*"

	payload := strings.NewReader("{\"language_id\":52,\"source_code\":\"I2luY2x1ZGUgPHN0ZGlvLmg+CgppbnQgbWFpbih2b2lkKSB7CiAgY2hhciBuYW1lWzEwXTsKICBzY2FuZigiJXMiLCBuYW1lKTsKICBwcmludGYoImhlbGxvLCAlc1xuIiwgbmFtZSk7CiAgcmV0dXJuIDA7Cn0=\",\"stdin\":\"SnVkZ2Uw\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", ApiHost)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil,err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var token models.RequestToken
	json.Unmarshal(body,&token)
	return &token,nil;
}

func GetReq() error{

	ApiKey := os.Getenv("JUDGE0_API_KEY")
	ApiHost := os.Getenv("JUDGE0_API_HOST")

	url := "https://judge0-ce.p.rapidapi.com/submissions/447b0a7b-0c27-4ca0-bf66-4f1506647785?base64_encoded=true&fields=*"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", ApiHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err;
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	return nil;
}
