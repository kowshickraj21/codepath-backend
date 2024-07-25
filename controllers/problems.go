package controllers

import (
	"database/sql"
	"main/models"

	"github.com/lib/pq"
)


func ViewProblem(db *sql.DB,pid int) (*models.Problem,error) {
	query := `SELECT * FROM problems WHERE pid = $1`;
	row := db.QueryRow(query,pid);

	var problem models.Problem
	err := row.Scan(&problem.Pid,&problem.Title,&problem.Description,pq.Array(&problem.Examples),pq.Array(&problem.Testcases));
	if err != nil{
		return nil,err;
	}

	return &problem,nil;
}


func FetchProblems(db *sql.DB) ([]models.Problem,error) {
	query := `SELECT pid,title FROM problems`;
	rows,err := db.Query(query);
	if err != nil {
        return nil, err
    }
    defer rows.Close()

	var problems []models.Problem
	for rows.Next() {
		var problem models.Problem
		err := rows.Scan(&problem.Pid,&problem.Title);
		if err != nil{
			return nil,err;
		}
		problems = append(problems,problem)
	}
	

	return problems,nil;
}