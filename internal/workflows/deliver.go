package workflows

// Deliver executes the delivery workflow using the orchestration engine
func Deliver(prompt string) map[string]any {
	return runWorkflowHelper("deliver", prompt)
}
