package scripts

import (
	"bytes"
	"context"
	"fmt"
	"main/models"
	"os"
	"os/exec"
	"sync"
	"time"
)

func CppExecuter(req models.Req, cases int) ([]models.ResStatus, int, error) {
	var res []models.ResStatus
	solved := 0

	compileCmd := exec.Command("g++", "-x", "c++", "-", "-o", "main")
	compileCmd.Stdin = bytes.NewBufferString(req.Code)
	compileOut, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = -1
		return nil, solved, fmt.Errorf("compilation error: %s", string(compileOut))
	}
	defer os.Remove("main")

	results := make(chan models.ResStatus, cases)

	var wg sync.WaitGroup

	runTestCase := func(index int) {
		defer wg.Done()

		input := req.Testcases[index].Input
		output := req.Testcases[index].Output

		var out models.ResStatus

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		runCmd := exec.CommandContext(ctx, "./main")

		if input != "" {
			runCmd.Stdin = bytes.NewBufferString(input)
		}

		runOutput, err := runCmd.CombinedOutput()

		if ctx.Err() == context.DeadlineExceeded {
			out.Id = 3
			out.Description = fmt.Sprintf("runtime error: %s", string(context.DeadlineExceeded.Error()))
			results <- out
			return
		}

		if err != nil {
			out.Id = 3
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
	os.Remove("main")

	for result := range results {
		if result.Id == 3 {
			return res,0,fmt.Errorf("%s",result.Description)
		}
		res = append(res, result)
	}

	return res, solved, nil
}
