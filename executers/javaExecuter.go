package executers

import (
	"io/ioutil"
	"main/models"
	"os"
	"os/exec"
)

func JavaExecuter (req models.Req) (models.ResStatus,error) {

	sourceFileName := "Main.java"
	inputFileName := "input.txt"
	var res models.ResStatus

	
	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		return res,err
	}
	defer os.Remove(sourceFileName) 

	
	if req.Input != "" {
		if err := ioutil.WriteFile(inputFileName, []byte(req.Input), 0644); err != nil {
			res.Id = 5
			res.Description = "Something went Wrong"
			return res,err
		}
	}
	defer os.Remove(inputFileName) 

	
	compileCmd := exec.Command("javac", sourceFileName)
	compileOutput, err := compileCmd.CombinedOutput()
	if err != nil {
		res.Id = 3
		res.Description = string(compileOutput)
		return res,nil
	}
	
	var runCmd *exec.Cmd
	if req.Input != "" {
		runCmd = exec.Command("java", "Main")
		runCmd.Stdin, _ = os.Open(inputFileName) 
	} else {
		runCmd = exec.Command("java", "Main")
	}

	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		res.Id = 4
		res.Description = string(runOutput)
		return res,nil
	}
	defer os.Remove("Main.class") 

	if(string(runOutput) == req.Output){
		res.Id = 1
		res.Description = "Accepted"
	}else{
		res.Id = 2
		res.Description = "Rejected"
	}
	return res,nil
}

