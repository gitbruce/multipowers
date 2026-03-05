package autosync

import "testing"

func TestDefaultPaths_AreStable(t *testing.T) {
	got := DefaultPaths("/tmp/p")
	if got.ProjectRoot != "/tmp/p/.multipowers/policy/autosync" {
		t.Fatalf("project root mismatch: %s", got.ProjectRoot)
	}
	if got.EventsRawDir != "/tmp/p/.multipowers/policy/autosync" {
		t.Fatalf("events dir mismatch: %s", got.EventsRawDir)
	}
	if got.GlobalSemanticFile == "" {
		t.Fatal("global semantic file should not be empty")
	}
}
