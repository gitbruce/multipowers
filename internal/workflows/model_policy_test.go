package workflows

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/gitbruce/multipowers/internal/providers"
)

func TestDebate_UsesAllConfiguredProviders(t *testing.T) {
	selection, err := providers.ConfiguredProvidersForWorkflow(filepath.Join("..", ".."), "debate")
	if err != nil {
		t.Fatalf("configured providers for debate: %v", err)
	}
	got := append([]string(nil), selection.ProviderProfiles...)
	sort.Strings(got)
	want := []string{"claude_code", "codex_cli", "gemini_cli"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("debate providers mismatch\nwant: %v\n got: %v", want, got)
	}
}

func TestBrainstorm_UsesConfiguredParallelProviders(t *testing.T) {
	selection, err := providers.ConfiguredProvidersForWorkflow(filepath.Join("..", ".."), "brainstorm")
	if err != nil {
		t.Fatalf("configured providers for brainstorm: %v", err)
	}
	if len(selection.Models) < 2 {
		t.Fatalf("expected brainstorm to configure parallel models, got %v", selection.Models)
	}
	if len(selection.ProviderProfiles) < 2 {
		t.Fatalf("expected brainstorm to map to multiple providers, got %v", selection.ProviderProfiles)
	}
}
