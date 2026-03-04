package isolation

import "testing"

func TestSelectTopCandidate(t *testing.T) {
	t.Run("select highest weighted score", func(t *testing.T) {
		candidates := []CandidateScore{
			{Model: "gpt-4o", WeightedScore: 8.2, ExecutionFailures: 0, DurationMs: 1400},
			{Model: "claude-sonnet", WeightedScore: 9.1, ExecutionFailures: 1, DurationMs: 1800},
			{Model: "deepseek-v3", WeightedScore: 8.9, ExecutionFailures: 0, DurationMs: 1300},
		}

		top, err := SelectTopCandidate(candidates)
		if err != nil {
			t.Fatalf("SelectTopCandidate error = %v", err)
		}
		if top.Model != "claude-sonnet" {
			t.Fatalf("top model = %q, want claude-sonnet", top.Model)
		}
	})

	t.Run("tie breakers are deterministic", func(t *testing.T) {
		candidates := []CandidateScore{
			{Model: "model-c", WeightedScore: 9.0, ExecutionFailures: 1, DurationMs: 1000},
			{Model: "model-b", WeightedScore: 9.0, ExecutionFailures: 0, DurationMs: 1200},
			{Model: "model-a", WeightedScore: 9.0, ExecutionFailures: 0, DurationMs: 1200},
		}

		top, err := SelectTopCandidate(candidates)
		if err != nil {
			t.Fatalf("SelectTopCandidate error = %v", err)
		}
		if top.Model != "model-a" {
			t.Fatalf("top model = %q, want model-a", top.Model)
		}
	})
}

func TestSelectTopCandidate_EmptyCandidates(t *testing.T) {
	if _, err := SelectTopCandidate(nil); err == nil {
		t.Fatal("expected error for empty candidates")
	}
}
