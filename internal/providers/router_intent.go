package providers

import "fmt"

// IntentRoutingResult represents the result of intent-based provider routing
type IntentRoutingResult struct {
	Intent              string   `json:"intent"`
	ProviderPolicy      string   `json:"provider_policy,omitempty"`
	Mode                string   `json:"mode"`
	AvailableProviders  []string `json:"available_providers"`
	SelectedProviders   []string `json:"selected_providers"`
	MinimumForSuccess   int      `json:"minimum_for_success"`
	Warnings            []string `json:"warnings,omitempty"`
	Error               string   `json:"error,omitempty"`
	Reason              string   `json:"reason,omitempty"`
	FallbackEnabled     bool     `json:"fallback_enabled"`
	SingleProviderMode  bool     `json:"single_provider_mode"`
}

// RouteIntent routes to appropriate providers based on intent
func RouteIntent(intent string, providerPolicy string) IntentRoutingResult {
	available := AvailableProviders()
	st := Degrade(intent, available)

	result := IntentRoutingResult{
		Intent:             intent,
		ProviderPolicy:     providerPolicy,
		Mode:               st.Mode,
		AvailableProviders: st.Available,
		SelectedProviders:  st.Selected,
		MinimumForSuccess:  st.MinimumForSuccess,
		Warnings:           st.Warnings,
		Error:              st.Error,
		SingleProviderMode: len(st.Selected) == 1,
		FallbackEnabled:    len(st.Available) > len(st.Selected),
	}

	// Add explainability reason
	result.Reason = buildRoutingReason(intent, st)

	return result
}

// buildRoutingReason creates a human-readable explanation of the routing decision
func buildRoutingReason(intent string, st Strategy) string {
	if st.Error != "" {
		return fmt.Sprintf("Routing failed: %s", st.Error)
	}

	switch intent {
	case "discover", "research":
		return fmt.Sprintf("Research/discovery mode: using %s for broad exploration", formatProviders(st.Selected))
	case "define", "plan":
		return fmt.Sprintf("Definition mode: using %s for structured planning", formatProviders(st.Selected))
	case "develop", "build":
		return fmt.Sprintf("Development mode: using %s for code generation", formatProviders(st.Selected))
	case "deliver", "review":
		return fmt.Sprintf("Delivery/review mode: using %s for quality validation", formatProviders(st.Selected))
	case "debate":
		return fmt.Sprintf("Debate mode: using %s for multi-perspective analysis", formatProviders(st.Selected))
	case "embrace", "multi":
		return fmt.Sprintf("Full workflow mode: using %s for comprehensive orchestration", formatProviders(st.Selected))
	default:
		if len(st.Selected) == 1 {
			return fmt.Sprintf("Single-provider mode: using %s", st.Selected[0])
		}
		return fmt.Sprintf("Multi-provider mode: using %s", formatProviders(st.Selected))
	}
}

// formatProviders formats a list of provider names for display
func formatProviders(providers []string) string {
	if len(providers) == 0 {
		return "none"
	}
	if len(providers) == 1 {
		return providers[0]
	}
	result := ""
	for i, p := range providers {
		if i > 0 {
			if i == len(providers)-1 {
				result += " and "
			} else {
				result += ", "
			}
		}
		result += p
	}
	return result
}

// IsValidIntent checks if an intent is recognized
func IsValidIntent(intent string) bool {
	validIntents := map[string]bool{
		"discover": true,
		"research": true,
		"define":   true,
		"plan":     true,
		"develop":  true,
		"build":    true,
		"deliver":  true,
		"review":   true,
		"debate":   true,
		"embrace":  true,
		"multi":    true,
	}
	return validIntents[intent]
}

// AllValidIntents returns all recognized intents
func AllValidIntents() []string {
	return []string{
		"discover", "research",
		"define", "plan",
		"develop", "build",
		"deliver", "review",
		"debate",
		"embrace", "multi",
	}
}
