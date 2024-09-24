package executers

import (
	"fmt"
	"main/models"
	"os"
	"os/exec"
)

func CppExecuter(req models.Req, cases int) ([]models.ResStatus, int, error) {
	sourceFileName := "main.cpp"
	inputFileName := "input.txt"
	var res []models.ResStatus
	solved := 0

	if err := os.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		return nil, -1, err
	}
	defer os.Remove(sourceFileName)

	compileCmd := exec.Command("g++", "-o", "main", sourceFileName)
	compileOut, err := compileCmd.CombinedOutput()
	if err != nil {
		return nil, -1, fmt.Errorf("compilation error: %s", string(compileOut))
	}
	defer os.Remove("main") 

	for i := 0; i < cases; i++ {
		input := req.Testcases[i].Input
		output := req.Testcases[i].Output
		var out models.ResStatus

		if input != "" {
			if err := os.WriteFile(inputFileName, []byte(input), 0644); err != nil {
				return nil, -1, err
			}
			defer os.Remove(inputFileName) 
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
			return nil, solved, fmt.Errorf("runtime error: %s", string(runOutput))
		}

		if string(runOutput) == output {
			solved++
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
