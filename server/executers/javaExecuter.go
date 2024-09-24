package executers

import (
	"fmt"
	"io/ioutil"
	"main/models"
	"os"
	"os/exec"
)

func JavaExecuter (req models.Req,cases int) ([]models.ResStatus,int,error) {

	sourceFileName := "Main.java"
	inputFileName := "input.txt"
	var res []models.ResStatus
	solved := 0

	
	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		solved = -1
		return nil,solved,err
	}
	defer os.Remove(sourceFileName) 

	compileCmd := exec.Command("javac", sourceFileName)
	Cout, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = -1
		return nil,solved, fmt.Errorf("compilation error: %s", string(Cout))
	}
	

	for i := 0;i < cases;i++{
		input := req.Testcases[i].Input
		output := req.Testcases[i].Output

		var out models.ResStatus

		if input != "" {
			if err := ioutil.WriteFile(inputFileName, []byte(input), 0644); err != nil {
				solved = -1
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
			return nil,solved,fmt.Errorf("runtime error: %s", string(runOutput))
		}
		defer os.Remove("Main.class") 


		if(string(runOutput) == output){
			solved += 1;
			out.Id = 1
			out.Description = "Accepted"
		}else{
			out.Id = 2
			out.Description = "Rejected"
		}
		res = append(res, out)
	}
	

	return res,solved,nil
}

