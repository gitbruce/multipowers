package workflows

func Deliver(prompt string) map[string]any {
	return map[string]any{"workflow": "deliver", "prompt": prompt}
}
