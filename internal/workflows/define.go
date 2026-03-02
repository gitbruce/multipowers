package workflows

import (
	"github.com/gitbruce/claude-octopus/internal/providers"
	"github.com/gitbruce/claude-octopus/internal/validation"
)

// DefineResult represents the result of a define workflow
type DefineResult struct {
	Workflow       string                        `json:"workflow"`
	Prompt         string                        `json:"prompt"`
	Validation     validation.TypedResult        `json:"validation,omitempty"`
	ProviderRoute  providers.IntentRoutingResult `json:"provider_route,omitempty"`
	Status         string                        `json:"status"`
}

// Define executes the definition workflow using atomic operations
// This is a compatibility facade that orchestrates atomic commands
func Define(prompt string) map[string]any {
	result := DefineResult{
		Workflow: "define",
		Prompt:   prompt,
		Status:   "initialized",
	}

	return map[string]any{
		"workflow": result.Workflow,
		"prompt":   result.Prompt,
		"status":   result.Status,
	}
}

// DefineWithValidation executes define with full validation
func DefineWithValidation(projectDir, prompt string) DefineResult {
	result := DefineResult{
		Workflow: "define",
		Prompt:   prompt,
		Status:   "initialized",
	}

	// Step 1: Validate workspace
	result.Validation = validation.ValidateByType(projectDir, validation.TypeWorkspace)
	if !result.Validation.Valid {
		result.Status = "blocked"
		return result
	}

	// Step 2: Route providers
	result.ProviderRoute = providers.RouteIntent("define", "")
	if result.ProviderRoute.Error != "" {
		result.Status = "error"
		return result
	}

	result.Status = "ready"
	return result
}
