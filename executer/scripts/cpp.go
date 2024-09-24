package scripts

import (
	"fmt"
	"io/ioutil"
	"main/models"
	"os"
	"os/exec"
)

func CppExecuter(req models.Req, cases int) ([]models.ResStatus, int, error) {

	sourceFileName := "main.cpp"
	inputFileName := "input.txt"
	var res []models.ResStatus
	solved := 0

	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		solved = -1
		return nil, solved, err
	}
	defer os.Remove(sourceFileName) 
	compileCmd := exec.Command("g++", "-o", "main", sourceFileName)
	compileOut, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = -1
		return nil, solved, fmt.Errorf("compilation error: %s",string(compileOut))
	}

	for i := 0; i < cases; i++ {
		input := req.Testcases[i].Input
		output := req.Testcases[i].Output

		var out models.ResStatus

			if input != "" {	
			if err := ioutil.WriteFile(inputFileName, []byte(input), 0644); err != nil {
				solved = -1
				return nil, solved, err
			}
		}

		var runCmd *exec.Cmd
		if input != "" {
			runCmd = exec.Command("./main")
			runCmd.Stdin, _ = os.Open(inputFileName)
		} else {
			runCmd = exec.Command("./main")
		}

		runOutput, err := runCmd.CombinedOutput()
		if err != nil {
			return nil, solved, fmt.Errorf("runtime error: %s",string(runOutput))
		}
		defer os.Remove("main") 

		if string(runOutput) == output {
			solved += 1
			out.Id = 1
			out.Description = "Accepted"
		} else {
			out.Id = 2
			out.Description = "Rejected"
		}
		res = append(res, out)
	}

	return res, solved, nil
}
