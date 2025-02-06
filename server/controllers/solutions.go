package controllers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"main/models"
	"strconv"

	"github.com/lib/pq"
)

func HandleSolutions(db *sql.DB, id int, jwt string) ([]models.Solutions,error){
	authUser,err := GetAuthUser(db, jwt)
	if err != nil {
		return nil,err
	}
	solutions,err := getSubmission(db,id,authUser.Email)
	if err != nil {
		return nil,err
	}
	return solutions,nil
}



func newSubmission(db *sql.DB,id string,code models.Code,email string,passed int) (error) {
	Pid ,_ := strconv.Atoi(id)
	var status string
	if passed == 6 {
		status = "Accepted"
	}else{
		status = "Rejected"
	}
	encodedCode := base64.StdEncoding.EncodeToString([]byte(code.Code))
	query := `INSERT INTO Solutions (Pid, Email, Code, Language,Status) VALUES ($1,$2,$3,$4,$5)`
	_,err := db.Exec(query,Pid,email,encodedCode,code.Language,status);
	if(err != nil){
		return err
	}
	return nil
}



func getSubmission(db *sql.DB,pid int,email string) ([]models.Solutions,error) {
	query := `SELECT * FROM Solutions WHERE pid = $1 AND email = $2`
	rows,err := db.Query(query,pid,email);
	if(err != nil){
		return nil,err
	}
	defer rows.Close()

	var submissions []models.Solutions
	for rows.Next() {
	    var submission models.Solutions
	    err = rows.Scan(&submission.Sid,&submission.Pid,&submission.Email,&submission.Code,&submission.Language,&submission.Status)
		if err != nil {
			return nil,err
		}
		decodedCode,_ := base64.StdEncoding.DecodeString(submission.Code)

		submission.Code = string(decodedCode)
	    submissions = append(submissions, submission)
    }
	return submissions,nil
}



func addSolved(db *sql.DB, email string,pid string) (error) {
	Pid,_ := strconv.Atoi(pid)
	var solved pq.Int64Array
	getQuery := `SELECT problems FROM Users WHERE email = $1`
	err := db.QueryRow(getQuery, email).Scan(&solved)
	if err != nil {
		fmt.Println("ERR : ",err)
		return err;
	}

	for i := range solved {
		if(solved[i] == int64(Pid)){
			fmt.Println("Already Exists")
			return nil;
		}
	}

	addQuery := `UPDATE Users SET problems = array_append(problems, $1) WHERE email = $2`
	_,err = db.Exec(addQuery,int64(Pid),email);
	if(err != nil){
		fmt.Println("ERR2 : ",err)
		return err
	}
	return nil
}