package executers

import (
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

	// Write the C++ code to a file
	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		solved = -1
		return nil, solved, err
	}
	defer os.Remove(sourceFileName) // Clean up the source file afterward

	// Compile the C++ code
	compileCmd := exec.Command("g++", "-o", "main", sourceFileName)
	_, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = -1
		return nil, solved, err
	}

	for i := 0; i < cases; i++ {
		input := req.Testcases[i].Input
		output := req.Testcases[i].Output

		var out models.ResStatus

		// If input is provided, write it to the input file
		if input != "" {
			if err := ioutil.WriteFile(inputFileName, []byte(input), 0644); err != nil {
				solved = -1
				return nil, solved, err
			}
		}

		var runCmd *exec.Cmd
		// Execute the compiled program, pass input if available
		if input != "" {
			runCmd = exec.Command("./main")
			runCmd.Stdin, _ = os.Open(inputFileName)
		} else {
			runCmd = exec.Command("./main")
		}

		// Capture the program's output
		runOutput, err := runCmd.CombinedOutput()
		if err != nil {
			out.Id = 4
			solved = -1
			out.Description = string(runOutput)
			return nil, solved, err
		}
		defer os.Remove("main") // Clean up the compiled binary afterward

		// Check if the program output matches the expected output
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
