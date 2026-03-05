package doctor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDoctor_AutoSyncDriftWarns(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, ".multipowers", "policy", "autosync")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p, "daily_stats.json"), []byte(`{"drift_rate":0.42}`), 0o644); err != nil {
		t.Fatal(err)
	}
	res := checkAutoSyncDrift(CheckContext{ProjectDir: d, Now: time.Now})
	if res.Status != StatusWarn {
		t.Fatalf("status=%s want warn", res.Status)
	}
}

func TestDoctor_UnhandledHighConfidenceProposalWarns(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, ".multipowers", "policy", "autosync")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	line := `{"rule_id":"r1","confidence":0.98,"status":"shadow"}` + "\n"
	if err := os.WriteFile(filepath.Join(p, "proposals.jsonl"), []byte(line), 0o644); err != nil {
		t.Fatal(err)
	}
	res := checkAutoSyncUnresolvedHighConfidence(CheckContext{ProjectDir: d, Now: time.Now})
	if res.Status != StatusWarn {
		t.Fatalf("status=%s want warn", res.Status)
	}
}
