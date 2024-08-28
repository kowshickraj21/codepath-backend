package controllers

import (
	"database/sql"
	"fmt"
	"main/models"
	"strconv"

	"github.com/lib/pq"
)

func newSubmission(db *sql.DB,submission models.Solutions) (error) {
	query := `INSERT INTO Solutions VALUES ($1,$2,$3,$4)`
	_,err := db.Exec(query,submission.Pid,submission.Email,submission.Code,submission.Language);
	if(err != nil){
		return err
	}
	return nil
}

func getSubmission(db *sql.DB,pid int,username string,language string) (*models.Solutions,error) {
	query := `SELECT * FROM Solutions WHERE pid = $1 AND username = $2 AND language = $3`
	rows,err := db.Query(query,pid,username,language);
	if(err != nil){
		return nil,err
	}
	var submission models.Solutions
	err = rows.Scan(&submission.Pid,&submission.Email,&submission.Code,&submission.Language)
	if err != nil {
		return nil,err;
	}
	return &submission,nil
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

	addQuery := `UPDATE Users SET problems = problems || $1 WHERE email = $2`
	_,err = db.Exec(addQuery,pq.Array([]int{Pid}),email);
	if(err != nil){
		fmt.Println("ERR2 : ",err)
		return err
	}
	return nil
}