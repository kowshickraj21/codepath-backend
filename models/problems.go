package models

type Problem struct {
	Pid         int
	Title       string
	Description string
	Examples    []io
	Testcases   []io
}

type io struct {
	input  string
	output string
}
