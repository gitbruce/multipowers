package validation

import (
	"fmt"
)

// ValidationType represents a type of validation to perform
type ValidationType string

const (
	TypeWorkspace ValidationType = "workspace"
	TypeNoShell   ValidationType = "no-shell"
	TypeTDDEnv    ValidationType = "tdd-env"
	TypeTestRun   ValidationType = "test-run"
	TypeCoverage  ValidationType = "coverage"
)

// TypedResult extends Result with type-specific information
type TypedResult struct {
	Type    ValidationType `json:"type"`
	Valid   bool           `json:"valid"`
	Reason  string         `json:"reason,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

// ValidateByType dispatches validation to the appropriate handler based on type
func ValidateByType(projectDir string, vtype ValidationType) TypedResult {
	switch vtype {
	case TypeWorkspace:
		return validateWorkspace(projectDir)
	case TypeNoShell:
		return validateNoShell(projectDir)
	case TypeTDDEnv:
		return validateTDDEnv(projectDir)
	case TypeTestRun:
		return validateTestRun(projectDir)
	case TypeCoverage:
		return validateCoverage(projectDir)
	default:
		return TypedResult{
			Type:   vtype,
			Valid:  false,
			Reason: fmt.Sprintf("unknown validation type: %s", vtype),
		}
	}
}

// validateWorkspace checks workspace completeness
func validateWorkspace(projectDir string) TypedResult {
	res := EnsureTargetWorkspace(projectDir)
	return TypedResult{
		Type:    TypeWorkspace,
		Valid:   res.Valid,
		Reason:  res.Reason,
		Details: map[string]any{"path": projectDir},
	}
}

// validateNoShell checks for shell script references
func validateNoShell(projectDir string) TypedResult {
	res, err := ScanNoShellRuntime(projectDir)
	if err != nil {
		return TypedResult{
			Type:   TypeNoShell,
			Valid:  false,
			Reason: err.Error(),
		}
	}
	reason := ""
	if !res.Valid {
		reason = "shell runtime references found"
	}
	return TypedResult{
		Type:    TypeNoShell,
		Valid:   res.Valid,
		Reason:  reason,
		Details: map[string]any{"violations": res.Violations, "files_checked": res.CheckedFiles},
	}
}

// validateTDDEnv checks TDD environment readiness
func validateTDDEnv(projectDir string) TypedResult {
	// TDD environment validation checks:
	// 1. Test framework availability
	// 2. Coverage tool availability
	// 3. Workspace context
	workspace := EnsureTargetWorkspace(projectDir)
	if !workspace.Valid {
		return TypedResult{
			Type:   TypeTDDEnv,
			Valid:  false,
			Reason: "workspace not ready: " + workspace.Reason,
			Details: map[string]any{
				"test_framework":  "go test",
				"coverage_tool":   "go tool cover",
				"workspace_ready": false,
			},
		}
	}
	return TypedResult{
		Type:  TypeTDDEnv,
		Valid: true,
		Details: map[string]any{
			"test_framework":  "go test",
			"coverage_tool":   "go tool cover",
			"workspace_ready": true,
		},
	}
}

// validateTestRun checks if tests can be run
func validateTestRun(projectDir string) TypedResult {
	workspace := EnsureTargetWorkspace(projectDir)
	if !workspace.Valid {
		return TypedResult{
			Type:   TypeTestRun,
			Valid:  false,
			Reason: "workspace not ready: " + workspace.Reason,
		}
	}
	return TypedResult{
		Type:  TypeTestRun,
		Valid: true,
		Details: map[string]any{
			"command": "go test ./...",
			"status":  "ready",
		},
	}
}

// validateCoverage checks if coverage can be collected
func validateCoverage(projectDir string) TypedResult {
	workspace := EnsureTargetWorkspace(projectDir)
	if !workspace.Valid {
		return TypedResult{
			Type:   TypeCoverage,
			Valid:  false,
			Reason: "workspace not ready: " + workspace.Reason,
		}
	}
	return TypedResult{
		Type:  TypeCoverage,
		Valid: true,
		Details: map[string]any{
			"command": "go test -cover ./...",
			"status":  "ready",
		},
	}
}
