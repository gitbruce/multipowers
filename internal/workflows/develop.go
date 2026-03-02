package workflows

import (
	"github.com/gitbruce/claude-octopus/internal/providers"
	"github.com/gitbruce/claude-octopus/internal/validation"
)

// DevelopResult represents the result of a develop workflow
type DevelopResult struct {
	Workflow       string                        `json:"workflow"`
	Prompt         string                        `json:"prompt"`
	Validation     validation.TypedResult        `json:"validation,omitempty"`
	TDDValidation  validation.TypedResult        `json:"tdd_validation,omitempty"`
	ProviderRoute  providers.IntentRoutingResult `json:"provider_route,omitempty"`
	TestResult     TestRunResult                 `json:"test_result,omitempty"`
	Status         string                        `json:"status"`
}

// Develop executes the development workflow using atomic operations
// This is a compatibility facade that orchestrates atomic commands
func Develop(prompt string) map[string]any {
	result := DevelopResult{
		Workflow: "develop",
		Prompt:   prompt,
		Status:   "initialized",
	}

	return map[string]any{
		"workflow": result.Workflow,
		"prompt":   result.Prompt,
		"status":   result.Status,
	}
}

// DevelopWithValidation executes develop with full validation
func DevelopWithValidation(projectDir, prompt string) DevelopResult {
	result := DevelopResult{
		Workflow: "develop",
		Prompt:   prompt,
		Status:   "initialized",
	}

	// Step 1: Validate workspace
	result.Validation = validation.ValidateByType(projectDir, validation.TypeWorkspace)
	if !result.Validation.Valid {
		result.Status = "blocked"
		return result
	}

	// Step 2: Validate TDD environment
	result.TDDValidation = validation.ValidateByType(projectDir, validation.TypeTDDEnv)
	if !result.TDDValidation.Valid {
		result.Status = "blocked"
		return result
	}

	// Step 3: Route providers
	result.ProviderRoute = providers.RouteIntent("develop", "")
	if result.ProviderRoute.Error != "" {
		result.Status = "error"
		return result
	}

	// Step 4: Run tests (optional, to check initial state)
	result.TestResult = TestRun(projectDir)

	result.Status = "ready"
	return result
}
