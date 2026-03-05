package fingerprint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFingerprint_ProbesRequiredDocs(t *testing.T) {
	d := t.TempDir()
	for _, name := range []string{"README.md", "CLAUDE.md", "AGENTS.md", "PRODUCT.md"} {
		if err := os.WriteFile(filepath.Join(d, name), []byte("# test\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	res, err := Scan(d)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	for _, dim := range []string{"docs.required"} {
		if len(res.EvidenceMap[dim]) == 0 {
			t.Fatalf("expected evidence for %s", dim)
		}
	}
}

func TestFingerprint_OutputIncludesEvidenceMapAndConfidence(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "go.mod"), []byte("module x\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	res, err := Scan(d)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(res.EvidenceMap) == 0 {
		t.Fatal("expected evidence_map")
	}
	if len(res.Confidence) == 0 {
		t.Fatal("expected confidence")
	}
}
