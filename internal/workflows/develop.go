package workflows

func Develop(prompt string) map[string]any {
	return map[string]any{"workflow": "develop", "prompt": prompt}
}
