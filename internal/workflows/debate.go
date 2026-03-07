package workflows

import (
	"github.com/gitbruce/multipowers/internal/providers"
)

var configuredProvidersForWorkflow = providers.ConfiguredProvidersForWorkflow

// Debate executes the debate workflow using the orchestration engine.
func Debate(prompt string) (map[string]any, bool) {
	result := runWorkflowHelper("debate", prompt)
	selection, err := configuredProvidersForWorkflow(".", "debate")
	if err != nil {
		result["configured_provider_error"] = err.Error()
		return result, false
	}
	result["configured_models"] = selection.Models
	result["configured_provider_profiles"] = selection.ProviderProfiles
	ok := providers.HasQuorum(len(selection.ProviderProfiles))
	return result, ok
}
