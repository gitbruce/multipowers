package cli

import (
	"path/filepath"
	"testing"
)

func TestCLIRawStream_Persists(t *testing.T) {
	d := t.TempDir()
	code := Run([]string{"status", "--dir", d, "--json"})
	if code != 0 {
		t.Fatalf("status code=%d", code)
	}
	matches, err := filepath.Glob(filepath.Join(d, ".multipowers", "policy", "autosync", "events.raw.*.jsonl"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected raw stream file")
	}
}
