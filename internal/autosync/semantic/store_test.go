package semantic

import (
	"os"
	"strings"
	"testing"
)

func TestSemanticStore_StoresPreferencePatternsOnly(t *testing.T) {
	d := t.TempDir()
	s := NewStore(d)
	if err := s.Upsert(Pattern{Key: "Language/Go", Value: "prefer_small_diffs", Score: 0.9, PrivateRaw: "/very/private/path"}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	b, err := os.ReadFile(s.Path())
	if err != nil {
		t.Fatalf("read store: %v", err)
	}
	if strings.Contains(string(b), "/very/private/path") {
		t.Fatalf("private raw leaked: %s", string(b))
	}
}

func TestSemanticMigration_RequiresSimilarityThreshold(t *testing.T) {
	d := DecisionInput{Similarity: 0.62, SimilarityThreshold: 0.7, ConflictRate: 0.01}
	got := EvaluateMigration(d)
	if got.Apply {
		t.Fatalf("expected no apply when below threshold: %+v", got)
	}
	if !got.ShadowOnly {
		t.Fatalf("expected shadow_only when below threshold: %+v", got)
	}
}

func TestSemanticMigration_CanRollbackOnConflictSpike(t *testing.T) {
	d := DecisionInput{Similarity: 0.9, SimilarityThreshold: 0.7, ConflictRate: 0.55}
	got := EvaluateMigration(d)
	if !got.Rollback {
		t.Fatalf("expected rollback on conflict spike: %+v", got)
	}
}
