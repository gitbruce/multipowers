package workflows

// Embrace executes the full workflow using the orchestration engine
func Embrace(prompt string) map[string]any {
	return runWorkflowHelper("embrace", prompt)
}
