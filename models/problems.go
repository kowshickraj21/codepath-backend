package models

type Problem struct {
	Pid         int
	Title       string
	Description string
	Examples    []IO
	Testcases   []IO
}

type IO struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}
