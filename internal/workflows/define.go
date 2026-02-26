package workflows

func Define(prompt string) map[string]any {
	return map[string]any{"workflow": "define", "prompt": prompt}
}
