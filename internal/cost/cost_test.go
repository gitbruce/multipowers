package cost

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEstimateFromPromptReturnsPositiveValues(t *testing.T) {
	rep := EstimateFromPrompt("implement migration plan")
	if rep.EstimatedInputTokens <= 0 {
		t.Fatalf("expected input tokens > 0, got %d", rep.EstimatedInputTokens)
	}
	if rep.EstimatedCostUSD <= 0 {
		t.Fatalf("expected cost > 0, got %f", rep.EstimatedCostUSD)
	}
}

func TestBuildReportAggregatesModelOutputs(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "model_outputs.2026-03-05.jsonl")
	payload := "{\"model\":\"gpt-5.3-codex\",\"tokens_input\":120,\"tokens_output\":80}\n" +
		"{\"model\":\"gpt-5.3-codex\",\"tokens_input\":20,\"tokens_output\":30}\n"
	if err := os.WriteFile(f, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}

	rep, err := BuildReport(d)
	if err != nil {
		t.Fatalf("expected report build success: %v", err)
	}
	if rep.TotalInputTokens != 140 {
		t.Fatalf("expected total input 140, got %d", rep.TotalInputTokens)
	}
	if rep.TotalOutputTokens != 110 {
		t.Fatalf("expected total output 110, got %d", rep.TotalOutputTokens)
	}
}
