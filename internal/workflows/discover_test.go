package workflows

import (
	"testing"
)

func TestDiscover(t *testing.T) {
	prompt := "Test prompt"
	result := Discover(prompt)

	if result["workflow"] != "discover" {
		t.Errorf("expected workflow discover, got %v", result["workflow"])
	}
	if result["prompt"] != prompt {
		t.Errorf("expected prompt %q, got %v", prompt, result["prompt"])
	}
	if result["status"] == "" {
		t.Error("expected status to be set")
	}
	if result["report"] == "" {
		t.Error("expected report to be generated")
	}
}
