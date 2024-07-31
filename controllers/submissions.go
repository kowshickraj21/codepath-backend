package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
	"os"
	"strings"
)


func CreateReq(code string) (*models.RequestToken,error) {

	file,err := os.ReadFile("problems/1/Main.java.txt")
	if err != nil {
		return nil,err;
	}
	fmt.Println(code)
	boilerplate := string(file)
	sourceCode := strings.Replace(boilerplate,"$",code,1) 
    input := "15"
    expectedOutput := "25"
	fmt.Println(sourceCode)
    encodedSourceCode := encodeBase64(sourceCode)

    requestPayload := models.Judge0Request{
        SourceCode:     encodedSourceCode,
        LanguageID:     91,
        Stdin:          encodeBase64(input),
        ExpectedOutput: encodeBase64(expectedOutput),
    }

	jsonData, err := json.Marshal(requestPayload)
    if err != nil {
        return nil,err
    }

	ApiKey := os.Getenv("JUDGE0_API_KEY")
	ApiHost := os.Getenv("JUDGE0_API_HOST")


	url := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=false&fields=*"

	// payload := strings.NewReader("{\"language_id\":52,\"source_code\":\"I2luY2x1ZGUgPHN0ZGlvLmg+CgppbnQgbWFpbih2b2lkKSB7CiAgY2hhciBuYW1lWzEwXTsKICBzY2FuZigiJXMiLCBuYW1lKTsKICBwcmludGYoImhlbGxvLCAlc1xuIiwgbmFtZSk7CiAgcmV0dXJuIDA7Cn0=\",\"stdin\":\"SnVkZ2Uw\"}")

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", ApiHost)
	req.Header.Add("Content-Type", "application/json")

	// fmt.Println(req)
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

func GetReq(token *models.RequestToken) (*models.Judge0Response,error){

	ApiKey := os.Getenv("JUDGE0_API_KEY")
	ApiHost := os.Getenv("JUDGE0_API_HOST")

	baseUrl := `https://judge0-ce.p.rapidapi.com/submissions/$?base64_encoded=true&fields=*`

	url := strings.Replace(baseUrl,"$",token.Token,1)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", ApiKey)
	req.Header.Add("x-rapidapi-host", ApiHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil,err;
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response models.Judge0Response
	fmt.Println(string(body));
	json.Unmarshal(body,&response)
	return &response,nil;
}

func encodeBase64(data string) string {
    return base64.StdEncoding.EncodeToString([]byte(data))
}

// func decodeBase64(encoded string) (string, error) {
//     decoded, err := base64.StdEncoding.DecodeString(encoded)
//     if err != nil {
//         return "", err
//     }
//     return string(decoded), nil
// }
