package providers

import (
	"path/filepath"
	"testing"
)

func TestQuorum(t *testing.T) {
	if HasQuorum(1) {
		t.Fatal("1 should not satisfy quorum")
	}
	if !HasQuorum(2) {
		t.Fatal("2 should satisfy quorum")
	}
}

func TestConfiguredProvidersForWorkflow_DebateProfiles(t *testing.T) {
	selection, err := ConfiguredProvidersForWorkflow(filepath.Join("..", ".."), "debate")
	if err != nil {
		t.Fatalf("configured providers: %v", err)
	}
	if len(selection.ProviderProfiles) < 2 {
		t.Fatalf("expected debate provider quorum, got %v", selection.ProviderProfiles)
	}
}
