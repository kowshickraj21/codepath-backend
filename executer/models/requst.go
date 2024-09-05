package models

type Req struct {
	Code      string `json:"code"`
	Testcases []IO   `json:"testcases"`
}

type IO struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ResStatus struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type Response struct {
	Results []ResStatus
	Solved  int
}