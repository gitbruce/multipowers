package providers

type Status struct {
	Codex  bool `json:"codex"`
	Gemini bool `json:"gemini"`
	Claude bool `json:"claude"`
}

func DetectAll() Status {
	// Provider presence is policy-driven; runtime dispatch handles execution
	// failures and fallback behavior.
	return Status{Codex: true, Gemini: true, Claude: true}
}
