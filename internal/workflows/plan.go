package workflows

// Plan executes the planning workflow using the orchestration engine.
func Plan(prompt string) map[string]any {
	return runWorkflowHelper("plan", prompt)
}
