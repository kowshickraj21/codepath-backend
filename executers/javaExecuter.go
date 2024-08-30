package executers

import (
	"io/ioutil"
	"main/models"
	"os"
	"os/exec"
)

func JavaExecuter (req models.Req) ([]models.ResStatus,bool,error) {

	sourceFileName := "Main.java"
	inputFileName := "input.txt"
	var res []models.ResStatus
	solved := true

	
	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		solved = false
		return nil,solved,err
	}
	defer os.Remove(sourceFileName) 

	compileCmd := exec.Command("javac", sourceFileName)
	_, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = false
		return nil,solved,err
	}

	for i := range req.Testcases{
		input := req.Testcases[i].Input
		output := req.Testcases[i].Output

		var out models.ResStatus

		if input != "" {
			if err := ioutil.WriteFile(inputFileName, []byte(input), 0644); err != nil {
				solved = false
				return nil,solved,err
			}
		}
		
		var runCmd *exec.Cmd
		if input != "" {
			runCmd = exec.Command("java", "Main")
			runCmd.Stdin, _ = os.Open(inputFileName) 
		} else {
			runCmd = exec.Command("java", "Main")
		}
	
		runOutput, err := runCmd.CombinedOutput()
		if err != nil {
			out.Id = 4
			solved = false
			out.Description = string(runOutput)
			return nil,solved,err
		}
		defer os.Remove("Main.class") 


		if(string(runOutput) == output){
			out.Id = 1
			out.Description = "Accepted"
		}else{
			solved = false
			out.Id = 2
			out.Description = "Rejected"
		}
		res = append(res, out)
	}
	

	return res,solved,nil
}

