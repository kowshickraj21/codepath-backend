package models

type Judge0Response struct {
	SourceCode     string      `json:"source_code"`
	LanguageID     int         `json:"language_id"`
	Stdin          string      `json:"stdin"`
	ExpectedOutput interface{} `json:"expected_output"`
	Stdout         string      `json:"stdout"`
	StatusID       int         `json:"status_id"`
	CreatedAt      string      `json:"created_at"`
	FinishedAt     string      `json:"finished_at"`
	Time           string      `json:"time"`
	Memory         int         `json:"memory"`
	Stderr         interface{} `json:"stderr"`
	Token          string      `json:"token"`
	NumberOfRuns   int         `json:"number_of_runs"`
	CpuTimeLimit   string      `json:"cpu_time_limit"`
}

type RequestToken struct {
	Token string `json:"token"`
}

type Judge0Request struct {
	SourceCode     string `json:"source_code"`
	LanguageID     int    `json:"language_id"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output,omitempty"`
}

type Code struct {
	Code string `json:"code"`
}