package models

import (
	"database/sql"

	"github.com/lib/pq"
)

type Problem struct {
	pid         int
	title       string
	description string
	examples    []io
	testcases   []io
}

type io struct {
	input  string
	output string
}

func ViewProblem(db *sql.DB,pid int) (*Problem,error) {
	query := `SELECT * FROM problems WHERE pid = $1`;
	row := db.QueryRow(query,pid);
	var problem Problem

	err := row.Scan(&problem.pid,&problem.title,&problem.description,pq.Array(&problem.examples),pq.Array(&problem.testcases));
	if err != nil{
		return nil,err;
	}

	return &problem,nil;
}

func FetchProblems(db *sql.DB) ([]Problem,error) {
	query := `SELECT pid,title FROM problems`;
	rows,err := db.Query(query);
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var problem Problem
		err := rows.Scan(&problem.pid,&problem.title);
		if err != nil{
			return nil,err;
		}
		problems = append(problems,problem)
	}
	

	return problems,nil;
}