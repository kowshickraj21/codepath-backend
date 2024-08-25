package controllers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lib/pq"
)


func CreateReq(db *sql.DB,code models.Code,id string) ([]models.Judge0Response,error) {

	sourceCode := readFile(code,id)
	Id,_ := strconv.Atoi(id)
	testcases,err := readCases(db,Id)
	var tokens []models.RequestToken
	if err != nil{
		return nil,err
	}

	for i := range testcases{

	input := testcases[i].Input
    expectedOutput := testcases[i].Output

    requestPayload := models.Judge0Request{
        SourceCode:     sourceCode,
        LanguageID:     91,  // HardCoded Java
        Stdin:          encodeBase64(input),
        ExpectedOutput: encodeBase64(expectedOutput),
    }

	jsonData, err := json.Marshal(requestPayload)
    if err != nil {
        return nil,err
    }

	ApiKey := os.Getenv("JUDGE0_API_KEY")
	ApiHost := os.Getenv("JUDGE0_API_HOST")

	url := `https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=false&fields=*`


	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

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
	tokens = append(tokens, token)
}
	res,err := GetReq(tokens)
	if err != nil {
	   return nil,err;
	}
	return res,nil;
}

func GetReq(tokens []models.RequestToken) ([]models.Judge0Response,error){

	ApiKey := os.Getenv("JUDGE0_API_KEY")
	ApiHost := os.Getenv("JUDGE0_API_HOST")
	var responses []models.Judge0Response

	baseUrl := `https://judge0-ce.p.rapidapi.com/submissions/$?base64_encoded=true&fields=*`

	for i := range tokens{

	url := strings.Replace(baseUrl,"$",tokens[i].Token,1)

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
	json.Unmarshal(body,&response)
	responses = append(responses, response)
	fmt.Println(response.StatusID);
}
	return responses,nil;
}

func encodeBase64(data string) string {
    return base64.StdEncoding.EncodeToString([]byte(data))
}

func readFile(code models.Code,id string) (string){
	fileurl := "problems/$1/Main.$2.txt"
	fileurl = strings.Replace(fileurl,"$1",id,1)
	fileurl = strings.Replace(fileurl,"$2",code.Language,1)
	
	file,err := os.ReadFile(fileurl)
	if err != nil {
		return "";
	}
	boilerplate := string(file)
	sourceCode := strings.Replace(boilerplate,"$",code.Code,1) 
	return encodeBase64(sourceCode)
}

func readCases(db *sql.DB, id int) ([]models.IO, error) {
	query := `SELECT testcases FROM problems WHERE pid = $1`
	row := db.QueryRow(query, id)

	var testcaseStr [] string
	var testcases []models.IO

	err := row.Scan(pq.Array(&testcaseStr))
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	for i := range testcaseStr{
		var testcase models.IO
	json.Unmarshal([]byte(testcaseStr[i]),&testcase);
	testcase.Input = strings.ReplaceAll(testcase.Input,"n","\n")
	testcases = append(testcases,testcase)
	}


	return testcases, nil
}
