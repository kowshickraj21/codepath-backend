package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"main/aws"
	"main/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lib/pq"
)


func HandleSubmissions(db *sql.DB, client *s3.Client, code models.Code,id string,jwt string) ([]models.ResStatus,error) {
	var outputs models.Response
	authUser,err := GetAuthUser(db, jwt)
	if(authUser == nil){
		return nil,err
	}
	outputs,passed,err := CreateReq(db,client,code,id,6)
	if err != nil {
		return nil,err;
	}
	if passed == 6 {
		addSolved(db,authUser.Email,id)
	}
	if  passed != -1 {
		newSubmission(db,id,code,authUser.Email,passed)
	}

	return outputs.Results,nil
}

func HandleRun(db *sql.DB,client *s3.Client, code models.Code, id,jwt string) ([]models.ResStatus,error) {
	var outputs models.Response
	authUser,err := GetAuthUser(db, jwt)
	if(authUser == nil){
		return nil,err
	}
	outputs,_,err = CreateReq(db,client,code,id,2)
	if err != nil {
		return nil,err;
	}

	return outputs.Results,nil
}

func CreateReq(db *sql.DB, client *s3.Client, code models.Code,id string,cases int) (models.Response,int,error) {
	
	sourceCode := readFile(code,client,id)
	Id,_ := strconv.Atoi(id)
	testcases,err := readCases(db,Id,cases)

	var res models.Response

	if err != nil{
		return res,-1,err
	}

	payload := models.Req{
        Code:     sourceCode,
		Testcases : testcases,
		Language : code.Language,
    }

	var solved int 
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return res,0,err
	}

	resp, err := http.Post(fmt.Sprintf("%s/execute",os.Getenv("EXECUTER_ORIGIN")), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return res,0,err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200{
	errBody,_ := io.ReadAll(resp.Body)
	return res,0,errors.New(string(errBody))
	}
	
	json.NewDecoder(resp.Body).Decode(&res)

	return res,solved,nil;
}

func readFile(code models.Code, client *s3.Client, id string) (string){

	bucket := os.Getenv("AWS_BUCKET");
	fileurl := "$1/Main.$2.txt"
	fileurl = strings.Replace(fileurl,"$1",id,1)
	fileurl = strings.Replace(fileurl,"$2",code.Language,1)

	headerFileUrl := "imports.$1.txt"
	headerFileUrl = strings.Replace(headerFileUrl,"$1",code.Language,1)
	
	file := aws.ReadFile(context.TODO(),client,bucket,fileurl)
	headerFile := aws.ReadFile(context.TODO(),client,bucket,headerFileUrl)

	boilerplate := string(file)
	headers := string(headerFile)
	sourceCode := strings.Replace(boilerplate,"#",headers,1)
	sourceCode = strings.Replace(sourceCode,"$",code.Code,1) 
	return sourceCode
}

func readCases(db *sql.DB, id int, cases int) ([]models.IO, error) {
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
	var finalCases []models.IO

	for i:=0; i<cases; i++ {
		finalCases = append(finalCases,testcases[i]);
	}

	return finalCases, nil
}
