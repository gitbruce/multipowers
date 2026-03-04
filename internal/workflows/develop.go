package workflows

// Develop executes the development workflow using the orchestration engine
func Develop(prompt string) map[string]any {
	return runWorkflowHelper("develop", prompt)
}
