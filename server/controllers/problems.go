package controllers

import (
	"database/sql"
	"encoding/json"
	"main/models"

	"github.com/lib/pq"
)


func ViewProblem(db *sql.DB,pid int) (*models.Problem,error) {
	query := `SELECT * FROM problems WHERE pid = $1`;
	row := db.QueryRow(query,pid);

	var problem models.Problem
	var exampleStr [] string
	var testcaseStr [] string
	err := row.Scan(&problem.Pid,&problem.Title,&problem.Description,pq.Array(&exampleStr),pq.Array(&testcaseStr));
	

	for i := range exampleStr{
		var example models.IO
	json.Unmarshal([]byte(exampleStr[i]),&example);
	problem.Examples = append(problem.Examples,example)
	if err != nil{
		return nil,err;
	}
	}

	for i := range testcaseStr{
		var testcase models.IO
	json.Unmarshal([]byte(testcaseStr[i]),&testcase);
	problem.Testcases = append(problem.Testcases,testcase)
	if err != nil{
		return nil,err;
	}
	}

	return &problem,nil;
}


func FetchProblems(db *sql.DB) ([]models.Problem, error) {
	query := `SELECT pid, title FROM problems`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []models.Problem

	for rows.Next() {

		var problem models.Problem
		err := rows.Scan(&problem.Pid, &problem.Title)
		problems = append(problems, problem)
		
	if err != nil {
		return nil, err
	}
}

	return problems, nil
}