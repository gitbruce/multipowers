package cli

import (
	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/providers"
	"github.com/gitbruce/multipowers/internal/validation"
)

// RuntimeStatus represents the comprehensive runtime health status
type RuntimeStatus struct {
	// Context status
	ContextComplete   bool     `json:"context_complete"`
	ContextMissing    []string `json:"context_missing,omitempty"`
	ContextPath       string   `json:"context_path,omitempty"`

	// Provider status
	ProvidersAvailable []string `json:"providers_available"`
	ProvidersCount     int      `json:"providers_count"`

	// Validation status
	LastValidation     string `json:"last_validation,omitempty"`
	ValidationStatus   string `json:"validation_status,omitempty"`

	// Hook status
	HookReady          bool   `json:"hook_ready"`
	HookEvents         []string `json:"hook_events,omitempty"`

	// Overall status
	Status             string `json:"status"`
	Ready              bool   `json:"ready"`
}

// GetRuntimeStatus aggregates runtime health information
func GetRuntimeStatus(projectDir string) RuntimeStatus {
	status := RuntimeStatus{
		Status:     "unknown",
		HookEvents: []string{"SessionStart", "UserPromptSubmit", "PreToolUse", "PostToolUse", "Stop"},
	}

	// Check context
	status.ContextComplete = ctxpkg.Complete(projectDir)
	status.ContextMissing = ctxpkg.Missing(projectDir)
	status.ContextPath = projectDir

	// Check providers
	available := providers.AvailableProviders()
	providerNames := make([]string, len(available))
	for i, p := range available {
		providerNames[i] = p.Name()
	}
	status.ProvidersAvailable = providerNames
	status.ProvidersCount = len(available)

	// Check validation (workspace)
	validationResult := validation.EnsureTargetWorkspace(projectDir)
	if validationResult.Valid {
		status.ValidationStatus = "passed"
		status.LastValidation = "workspace"
	} else {
		status.ValidationStatus = "failed: " + validationResult.Reason
	}

	// Hooks are always ready (they're part of the runtime)
	status.HookReady = true

	// Determine overall status
	if status.ContextComplete && len(available) > 0 {
		status.Status = "ready"
		status.Ready = true
	} else if !status.ContextComplete {
		status.Status = "context_incomplete"
		status.Ready = false
	} else if len(available) == 0 {
		status.Status = "no_providers"
		status.Ready = false
	} else {
		status.Status = "degraded"
		status.Ready = false
	}

	return status
}
