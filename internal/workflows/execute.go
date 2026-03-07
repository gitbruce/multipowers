package workflows

// Execute executes the implementation workflow using the orchestration engine.
func Execute(prompt string) map[string]any {
	return runWorkflowHelper("execute", prompt)
}
