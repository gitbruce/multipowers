package workflows

import (
	"github.com/gitbruce/claude-octopus/internal/providers"
)

// Debate executes the debate workflow using the orchestration engine
func Debate(prompt string) (map[string]any, bool) {
	result := runWorkflowHelper("debate", prompt)
	
	// Determine quorum from providers
	avail := providers.AvailableProviders()
	ok := providers.HasQuorum(len(avail))
	
	return result, ok
}
