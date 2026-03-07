package workflows

// Debug executes the debugging workflow using the orchestration engine.
func Debug(prompt string) map[string]any {
	return runWorkflowHelper("debug", prompt)
}
