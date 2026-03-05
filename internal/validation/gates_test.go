package validation

import (
	"os"
	"path/filepath"
	"testing"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
)

func TestEnsureTargetWorkspace(t *testing.T) {
	d := t.TempDir()
	res := EnsureTargetWorkspace(d)
	if res.Valid {
		t.Fatal("expected invalid without workspace")
	}
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	res = EnsureTargetWorkspace(d)
	if !res.Valid {
		t.Fatalf("expected valid, got %+v", res)
	}
	_ = os.Remove(filepath.Join(d, ".multipowers", "CLAUDE.md"))
	res = EnsureTargetWorkspace(d)
	if res.Valid {
		t.Fatal("expected invalid when required file missing")
	}
}
