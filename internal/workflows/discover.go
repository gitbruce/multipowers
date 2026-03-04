package workflows

// Discover executes the discovery workflow using the orchestration engine
func Discover(prompt string) map[string]any {
	return runWorkflowHelper("discover", prompt)
}
