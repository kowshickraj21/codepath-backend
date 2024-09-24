package models

type ResStatus struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type Response struct {
	Results []ResStatus `json:"results"`
	Solved  int         `json:"solved"`
}

type ExecStatus struct {
	Id             int    `json:"id"`
	Output         string `json:"output"`
	ExpectedOutput string `json:"ExpectedOutput"`
}

type Code struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type Req struct {
	Code      string `json:"code"`
	Testcases []IO   `json:"testcases"`
	Language  string `json:"language"`
}