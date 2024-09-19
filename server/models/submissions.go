package models

type ResStatus struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
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
}