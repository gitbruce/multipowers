package modelroute

import "testing"

func TestResolveForPromptDefault(t *testing.T) {
	r := ResolveForPrompt(t.TempDir(), "/mp:discover oauth patterns")
	if r.Provider != "gemini" {
		t.Fatalf("expected gemini for discover, got %+v", r)
	}
	if r.Model == "" {
		t.Fatalf("expected model to be resolved, got %+v", r)
	}
}

func TestResolveForPromptDevelop(t *testing.T) {
	r := ResolveForPrompt(t.TempDir(), "/mp:develop auth system")
	if r.Role != "heavy_coding" {
		t.Fatalf("expected heavy_coding role, got %+v", r)
	}
}
