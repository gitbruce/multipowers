package workflows

import (
	"github.com/gitbruce/claude-octopus/internal/providers"
	"github.com/gitbruce/claude-octopus/internal/validation"
)

// DeliverResult represents the result of a deliver workflow
type DeliverResult struct {
	Workflow        string                        `json:"workflow"`
	Prompt          string                        `json:"prompt"`
	Validation      validation.TypedResult        `json:"validation,omitempty"`
	NoShellResult   validation.TypedResult        `json:"no_shell_validation,omitempty"`
	ProviderRoute   providers.IntentRoutingResult `json:"provider_route,omitempty"`
	TestResult      TestRunResult                 `json:"test_result,omitempty"`
	CoverageResult  CoverageResult                `json:"coverage_result,omitempty"`
	Status          string                        `json:"status"`
}

// Deliver executes the delivery workflow using atomic operations
// This is a compatibility facade that orchestrates atomic commands
func Deliver(prompt string) map[string]any {
	result := DeliverResult{
		Workflow: "deliver",
		Prompt:   prompt,
		Status:   "initialized",
	}

	return map[string]any{
		"workflow": result.Workflow,
		"prompt":   result.Prompt,
		"status":   result.Status,
	}
}

// DeliverWithValidation executes deliver with full validation
func DeliverWithValidation(projectDir, prompt string, coverageThreshold float64) DeliverResult {
	result := DeliverResult{
		Workflow: "deliver",
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
	result.ProviderRoute = providers.RouteIntent("deliver", "")
	if result.ProviderRoute.Error != "" {
		result.Status = "error"
		return result
	}

	// Step 3: Run tests
	result.TestResult = TestRun(projectDir)
	if result.TestResult.Status == "error" {
		result.Status = "error"
		return result
	}
	if result.TestResult.Status == "failed" {
		result.Status = "blocked"
		return result
	}

	// Step 4: Check coverage
	result.CoverageResult = CoverageCheck(projectDir, coverageThreshold)
	if result.CoverageResult.Status == "failed" {
		result.Status = "blocked"
		return result
	}

	// Step 5: Validate no-shell runtime
	result.NoShellResult = validation.ValidateByType(projectDir, validation.TypeNoShell)
	if !result.NoShellResult.Valid {
		result.Status = "blocked"
		return result
	}

	result.Status = "ready"
	return result
}
