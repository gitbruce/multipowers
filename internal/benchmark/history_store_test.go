package benchmark

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadHistoryJudgeRecords(t *testing.T) {
	d := t.TempDir()
	path := filepath.Join(d, "judge_scores.2026-03-05.jsonl")
	content := "{\"run_id\":\"1\",\"judged_model\":\"model-a\",\"signature\":\"develop|code||\",\"weighted_score\":4.2}\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	recs, err := LoadHistoryJudgeRecords(d)
	if err != nil {
		t.Fatalf("LoadHistoryJudgeRecords error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if recs[0].Model != "model-a" || recs[0].Signature != "develop|code||" {
		t.Fatalf("unexpected record: %+v", recs[0])
	}
}
