package scripts

import (
	"fmt"
	"io/ioutil"
	"main/models"
	"os"
	"os/exec"
	"strings"
)

func JavaExecuter (req models.Req,cases int) ([]models.ResStatus,int,error) {

	sourceFileName := "Main.java"
	inputFileName := "input.txt"
	var res []models.ResStatus
	solved := 0
	
	if err := ioutil.WriteFile(sourceFileName, []byte(req.Code), 0644); err != nil {
		solved = -1
		fmt.Println("Error",err)
		return nil,solved,err
	}
	defer os.Remove(sourceFileName) 


	compileCmd := exec.Command("javac", sourceFileName)
	_, err := compileCmd.CombinedOutput()
	if err != nil {
		solved = -1
		fmt.Println("Error",err)
		return nil,solved,err
	}

	for i := 0;i < cases;i++{
	 
		input := strings.ReplaceAll(req.Testcases[i].Input,"n","\n")
		output := req.Testcases[i].Output
		fmt.Println("/",output,"/")	
		var out models.ResStatus

		if input != "" {
			if err := ioutil.WriteFile(inputFileName, []byte(input), 0644); err != nil {
				solved = -1
				fmt.Println("Error",err)
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
			solved = -1;
			fmt.Println("Error",err)
			out.Description = string(runOutput)
			return nil,solved,err
		}
		defer os.Remove("Main.class") 

		fmt.Println("/",string(runOutput),"/")
		if(strings.Trim(string(runOutput)," ") == (strings.Trim(output," "))){

			solved += 1;
			out.Id = 1
			out.Description = "Accepted"
		}else{
			out.Id = 2
			out.Description = "Rejected"
		}
		res = append(res, out)
	}
	
	fmt.Println("Res",res)
	return res,solved,nil
}

