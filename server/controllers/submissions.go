package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/executers"
	"main/models"
	"os"
	"strconv"
	"strings"

	"github.com/lib/pq"
)


func HandleSubmissions(db *sql.DB,code models.Code,id string,jwt string) ([]models.ResStatus,error) {
	authUser,err := GetAuthUser(db, jwt)
	if(authUser == nil){
		return nil,err
	}
	outputs,passed,err := CreateReq(db,code,id,6)
	if err != nil {
		return nil,err;
	}
	if  passed != -1 {
		addSolved(db,authUser.Email,id)
		newSubmission(db,id,code,authUser.Email,string(passed))
	}

	return outputs,nil
}

func HandleRun(db *sql.DB,code models.Code,id string,jwt string) ([]models.ResStatus,error) {
	authUser,err := GetAuthUser(db, jwt)
	if(authUser == nil){
		return nil,err
	}
	outputs,_,err := CreateReq(db,code,id,2)
	if err != nil {
		return nil,err;
	}

	return outputs,nil
}

func CreateReq(db *sql.DB,code models.Code,id string,cases int) ([]models.ResStatus,int,error) {
	
	sourceCode := readFile(code,id)
	Id,_ := strconv.Atoi(id)
	testcases,err := readCases(db,Id)
	if err != nil{
		return nil,-1,err
	}

	payload := models.Req{
        Code:     sourceCode,
		Testcases : testcases,
    }

	res,solved,err := executers.JavaExecuter(payload,cases)
	if err != nil {
		return nil,-1,err
	}

	return res,solved,nil;
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
	return sourceCode
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
