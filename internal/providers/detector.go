package providers

import "os/exec"

type Status struct {
	Codex  bool `json:"codex"`
	Gemini bool `json:"gemini"`
	Claude bool `json:"claude"`
}

func DetectAll() Status {
	_, cErr := exec.LookPath("codex")
	_, gErr := exec.LookPath("gemini")
	_, clErr := exec.LookPath("claude")
	return Status{Codex: cErr == nil, Gemini: gErr == nil, Claude: clErr == nil}
}
