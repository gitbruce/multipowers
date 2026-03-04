package workflows

// Define executes the definition workflow using the orchestration engine
func Define(prompt string) map[string]any {
	return runWorkflowHelper("define", prompt)
}
