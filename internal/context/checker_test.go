package context

import (
	"os"
	"path/filepath"
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
