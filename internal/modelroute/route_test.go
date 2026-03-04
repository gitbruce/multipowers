package modelroute

import (
	"testing"

	"github.com/gitbruce/claude-octopus/internal/benchmark"
)

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

func TestResolveBestModelFromHistory(t *testing.T) {
	signature := benchmark.BuildSimilaritySignature("develop", []string{"api"}, "gin", "go")
	records := make([]benchmark.HistoryJudgeRecord, 0, 20)
	for i := 0; i < 10; i++ {
		records = append(records, benchmark.HistoryJudgeRecord{
			Model:         "claude-opus",
			Signature:     signature,
			WeightedScore: 4.7,
		})
	}

	model, samples, ok := ResolveBestModelFromHistory(SmartRoutingRequest{
		Enabled:            true,
		Signature:          signature,
		MinSamplesPerModel: 10,
		Records:            records,
	})
	if !ok {
		t.Fatal("expected history override")
	}
	if model != "claude-opus" {
		t.Fatalf("model = %q, want claude-opus", model)
	}
	if samples != 10 {
		t.Fatalf("samples = %d, want 10", samples)
	}
}
