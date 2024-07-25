package controllers

import (
	"database/sql"
	"main/models"
)

type problemController struct{
	DB *sql.DB
}

func newProblemController(db *sql.DB) *problemController {
	return &problemController{DB:db}
}

func (pc *problemController) ViewProblem(pid int) (*models.Problem,error){
	problem,err := models.ViewProblem(pc.DB,pid)
	if err != nil {
		return nil,err;
	}
	return problem,nil;
}

func (pc *problemController) FetchProblems() ([]models.Problem,error){
	problem,err := models.FetchProblems(pc.DB)
	if err != nil {
		return nil,err;
	}
	return problem,nil;
}