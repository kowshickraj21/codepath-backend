package models

type Problem struct {
	Pid         int      `json:"pid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Examples    []IO     `json:"examples"`
	Testcases   []IO     `json:"testcases"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
}

type IO struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}
