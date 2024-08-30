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



func newSubmission(db *sql.DB,id string,code models.Code,email string) (error) {
	Pid ,_ := strconv.Atoi(id)
	encodedCode := base64.StdEncoding.EncodeToString([]byte(code.Code))
	query := `INSERT INTO Solutions (Pid, Email, Code, Language) VALUES ($1,$2,$3,$4)`
	_,err := db.Exec(query,Pid,email,encodedCode,code.Language);
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
	    err = rows.Scan(&submission.Sid,&submission.Pid,&submission.Email,&submission.Code,&submission.Language)
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
	var solved []int
	getQuery := `SELECT problems FROM Users WHERE email = $1`
	err := db.QueryRow(getQuery, email).Scan(pq.Array(&solved))
	if err != nil {
		fmt.Println("ERR : ",err)
		return err;
	}
	fmt.Println("Problems: ",solved)
	for i := range solved {
		if(solved[i] == Pid){
			fmt.Println("Already Exists")
			return nil;
		}
	}

	addQuery := `UPDATE Users SET problems = array_append(problems, $1) WHERE email = $2`
	_,err = db.Exec(addQuery,Pid,email);
	if(err != nil){
		fmt.Println("ERR2 : ",err)
		return err
	}
	return nil
}