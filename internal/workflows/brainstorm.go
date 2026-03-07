package workflows

// Brainstorm executes the brainstorming workflow using the orchestration engine.
func Brainstorm(prompt string) map[string]any {
	return runWorkflowHelper("brainstorm", prompt)
}
