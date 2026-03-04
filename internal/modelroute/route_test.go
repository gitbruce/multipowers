package modelroute

import "testing"

func TestResolveForPromptDefault(t *testing.T) {
	r := ResolveForPrompt(t.TempDir(), "/mp:discover oauth patterns")
	if r.Command != "discover" {
		t.Fatalf("expected discover command, got %+v", r)
	}
	if r.Model != "" {
		t.Fatalf("expected no default model in legacy shim, got %+v", r)
	}
}

func TestResolveForPromptDevelop(t *testing.T) {
	r := ResolveForPrompt(t.TempDir(), "/mp:develop auth system")
	if r.Command != "develop" {
		t.Fatalf("expected develop command, got %+v", r)
	}
}
