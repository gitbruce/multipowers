package workflows

// Design executes the design workflow using the orchestration engine.
func Design(prompt string) map[string]any {
	return runWorkflowHelper("design", prompt)
}
