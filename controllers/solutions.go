package controllers

import (
	"database/sql"
	"main/models"
)

func newSubmission(db *sql.DB,submission models.Solutions) (sql.Result,error) {
	query := `INSERT INTO Solutions VALUES ($1,$2,$3,$4)`
	res,err := db.Exec(query,submission.Pid,submission.Username,submission.Code,submission.Language);
	if(err != nil){
		return nil,err
	}
	return res,nil
}

func getSubmission(db *sql.DB,pid int,username string,language string) (*models.Solutions,error) {
	query := `SELECT * FROM Solutions WHERE pid = $1 AND username = $2 AND language = $3`
	rows,err := db.Query(query,pid,username,language);
	if(err != nil){
		return nil,err
	}
	var submission models.Solutions
	err = rows.Scan(&submission.Pid,&submission.Username,&submission.Code,&submission.Language)
	if err != nil {
		return nil,err;
	}
	return &submission,nil
}