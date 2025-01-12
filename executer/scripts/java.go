package scripts

import (
	"bytes"
	"fmt"
	"main/models"
	"os"
	"os/exec"
	"sync"
)

func JavaExecuter(req models.Req, cases int) ([]models.ResStatus, int, error) {
	sourceFileName := "Main.java"
	var res []models.ResStatus
	solved := 0

	if err := os.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		solved = -1
		return nil, solved, err
	}
	defer os.Remove(sourceFileName)

	compileCmd := exec.Command("javac", sourceFileName)
	Cout, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = -1
		return nil, solved, fmt.Errorf("compilation error: %s", string(Cout))
	}
	defer os.Remove("Main.class")
	os.Remove(sourceFileName)

	results := make(chan models.ResStatus, cases)

	var wg sync.WaitGroup

	runTestCase := func(index int) {
		defer wg.Done()

		input := req.Testcases[index].Input
		output := req.Testcases[index].Output

		var out models.ResStatus

		runCmd := exec.Command("java", "Main")

		if input != "" {
			runCmd.Stdin = bytes.NewBufferString(input)
		}

		runOutput, err := runCmd.CombinedOutput()
		if err != nil {
			out.Id = 2
			out.Description = fmt.Sprintf("runtime error: %s", string(runOutput))
			results <- out
			return
		}

		if string(runOutput) == output {
			out.Id = 1
			out.Description = "Accepted"
			solved++
		} else {
			out.Id = 2
			out.Description = "Rejected"
		}

		results <- out
	}

	for i := 0; i < cases; i++ {
		wg.Add(1)
		go runTestCase(i)
	}

	wg.Wait()
	close(results)

	defer os.Remove("Main.class")

	for result := range results {
		res = append(res, result)
	}

	return res, solved, nil
}
