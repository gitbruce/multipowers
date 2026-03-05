package providers

import (
	"testing"
)

func TestRouteIntent_Discover(t *testing.T) {
	result := RouteIntent("discover", "")
	if result.Intent != "discover" {
		t.Errorf("expected intent discover, got %s", result.Intent)
	}
	if result.Error != "" {
		// May fail if no providers available, that's OK for this test
		t.Logf("discover routing returned error (no providers?): %s", result.Error)
	}
	if result.Reason == "" {
		t.Error("expected non-empty reason for routing decision")
	}
}

func TestRouteIntent_Develop(t *testing.T) {
	result := RouteIntent("develop", "")
	if result.Intent != "develop" {
		t.Errorf("expected intent develop, got %s", result.Intent)
	}
	if result.Reason == "" {
		t.Error("expected non-empty reason for routing decision")
	}
}

func TestRouteIntent_Deliver(t *testing.T) {
	result := RouteIntent("deliver", "")
	if result.Intent != "deliver" {
		t.Errorf("expected intent deliver, got %s", result.Intent)
	}
	if result.Reason == "" {
		t.Error("expected non-empty reason for routing decision")
	}
}

func TestRouteIntent_Debate(t *testing.T) {
	result := RouteIntent("debate", "")
	if result.Intent != "debate" {
		t.Errorf("expected intent debate, got %s", result.Intent)
	}
	// Debate requires 2+ providers
	if len(result.AvailableProviders) < 2 && result.Error == "" {
		t.Error("expected error when insufficient providers for debate")
	}
	if result.MinimumForSuccess != 2 {
		t.Errorf("expected minimum 2 for debate, got %d", result.MinimumForSuccess)
	}
}

func TestRouteIntent_Fallback(t *testing.T) {
	result := RouteIntent("develop", "")
	// Check fallback flag is set correctly
	if len(result.AvailableProviders) > len(result.SelectedProviders) {
		if !result.FallbackEnabled {
			t.Error("expected fallback enabled when more providers available than selected")
		}
	}
}

func TestRouteIntent_ProviderPolicy(t *testing.T) {
	result := RouteIntent("develop", "claude-first")
	if result.ProviderPolicy != "claude-first" {
		t.Errorf("expected provider policy claude-first, got %s", result.ProviderPolicy)
	}
}

func TestFormatProviders(t *testing.T) {
	tests := []struct {
		providers []string
		expected  string
	}{
		{[]string{}, "none"},
		{[]string{"claude"}, "claude"},
		{[]string{"claude", "codex"}, "claude and codex"},
		{[]string{"claude", "codex", "gemini"}, "claude, codex and gemini"},
	}

	for _, tc := range tests {
		result := formatProviders(tc.providers)
		if result != tc.expected {
			t.Errorf("formatProviders(%v) = %s, expected %s", tc.providers, result, tc.expected)
		}
	}
}

func TestBuildRoutingReason(t *testing.T) {
	tests := []struct {
		intent    string
		hasError  bool
		wantEmpty bool
	}{
		{"discover", false, false},
		{"develop", false, false},
		{"deliver", false, false},
		{"debate", false, false},
		{"unknown", false, false},
		{"develop", true, false}, // Error case
	}

	for _, tc := range tests {
		st := Strategy{
			Mode:     tc.intent,
			Selected: []string{"claude"},
		}
		if tc.hasError {
			st.Error = "test error"
		}
		reason := buildRoutingReason(tc.intent, st)
		if (reason == "") != tc.wantEmpty {
			t.Errorf("intent=%s, hasError=%v: reason=%s, wantEmpty=%v", tc.intent, tc.hasError, reason, tc.wantEmpty)
		}
	}
}
