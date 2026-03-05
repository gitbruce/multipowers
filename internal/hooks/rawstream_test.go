package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestHookRawStream_Persists(t *testing.T) {
	d := t.TempDir()
	_ = Handle(d, api.HookEvent{Event: "SessionStart"})
	matches, err := filepath.Glob(filepath.Join(d, ".multipowers", "policy", "autosync", "events.raw.*.jsonl"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("expected raw stream file")
	}
	if _, err := os.Stat(matches[0]); err != nil {
		t.Fatalf("expected file exists: %v", err)
	}
}
