package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMissingAndComplete(t *testing.T) {
	d := t.TempDir()
	if Complete(d) {
		t.Fatal("should be incomplete")
	}
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	if !Complete(d) {
		t.Fatalf("expected complete, missing=%v", Missing(d))
	}
	p := filepath.Join(d, ".multipowers", "CLAUDE.md")
	_ = os.Remove(p)
	if Complete(d) {
		t.Fatal("should be incomplete when CLAUDE.md missing")
	}
}

func TestCompleteRequiresCanonicalTracksRegistryPath(t *testing.T) {
	d := t.TempDir()
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	canonical := filepath.Join(d, ".multipowers", "tracks", "tracks.md")
	if err := os.Remove(canonical); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(d, ".multipowers", "tracks.md")
	if err := os.WriteFile(legacy, []byte(strings.Repeat("legacy track\n", 10)), 0o644); err != nil {
		t.Fatal(err)
	}

	if Complete(d) {
		t.Fatal("legacy tracks.md should not satisfy context completeness")
	}
	missing := Missing(d)
	found := false
	for _, name := range missing {
		if name == "tracks/tracks.md" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected canonical tracks registry to be missing, got %v", missing)
	}
}
