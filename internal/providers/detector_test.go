package providers

import "testing"

func TestExternalProvidersAvailabilityIsNotBinaryGated(t *testing.T) {
	if !(Codex{}).Available() {
		t.Fatalf("expected codex provider to be considered available without LookPath gating")
	}
	if !(Gemini{}).Available() {
		t.Fatalf("expected gemini provider to be considered available without LookPath gating")
	}
}
