package workflows

import (
	"github.com/gitbruce/claude-octopus/internal/providers"
	"github.com/gitbruce/claude-octopus/internal/validation"
)

// DiscoverResult represents the result of a discover workflow
type DiscoverResult struct {
	Workflow       string                   `json:"workflow"`
	Prompt         string                   `json:"prompt"`
	Validation     validation.TypedResult   `json:"validation,omitempty"`
	ProviderRoute  providers.IntentRoutingResult `json:"provider_route,omitempty"`
	Providers      int                      `json:"providers"`
	Status         string                   `json:"status"`
}

// Discover executes the discovery workflow using atomic operations
// This is a compatibility facade that orchestrates atomic commands
func Discover(prompt string) map[string]any {
	avail := providers.AvailableProviders()

	result := DiscoverResult{
		Workflow:  "discover",
		Prompt:    prompt,
		Providers: len(avail),
		Status:    "initialized",
	}

	// Return structured result
	return map[string]any{
		"workflow":  result.Workflow,
		"prompt":    result.Prompt,
		"providers": result.Providers,
		"status":    result.Status,
	}
}

// DiscoverWithValidation executes discover with full validation
func DiscoverWithValidation(projectDir, prompt string) DiscoverResult {
	avail := providers.AvailableProviders()

	result := DiscoverResult{
		Workflow:  "discover",
		Prompt:    prompt,
		Providers: len(avail),
		Status:    "initialized",
	}

	// Step 1: Validate workspace
	result.Validation = validation.ValidateByType(projectDir, validation.TypeWorkspace)
	if !result.Validation.Valid {
		result.Status = "blocked"
		return result
	}

	// Step 2: Route providers
	result.ProviderRoute = providers.RouteIntent("discover", "")
	if result.ProviderRoute.Error != "" {
		result.Status = "error"
		return result
	}

	result.Status = "ready"
	return result
}
