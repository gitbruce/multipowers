package providers

func Registry() []Provider {
	return []Provider{Codex{}, Gemini{}, Claude{}}
}
