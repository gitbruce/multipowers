package autosync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRawStream_EmitWritesFile(t *testing.T) {
	d := t.TempDir()
	path, err := EmitRawEvent(d, "test", "emit", map[string]any{"k": "v"})
	if err != nil {
		t.Fatalf("emit error: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(b), "\"source\":\"test\"") {
		t.Fatalf("missing source in jsonl: %s", string(b))
	}
	if !strings.HasSuffix(path, filepath.Join(".multipowers", "policy", "autosync", filepath.Base(path))) {
		t.Fatalf("unexpected path: %s", path)
	}
}
