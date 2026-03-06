package validation

import (
	"os"
	"path/filepath"
	"strings"
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

func TestEnsureTargetWorkspaceRequiresCanonicalTracksRegistryPath(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	canonical := filepath.Join(d, ".multipowers", "tracks", "tracks.md")
	if err := os.Remove(canonical); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, ".multipowers", "tracks.md"), []byte(strings.Repeat("legacy track\n", 10)), 0o644); err != nil {
		t.Fatal(err)
	}

	res := EnsureTargetWorkspace(d)
	if res.Valid {
		t.Fatal("expected legacy-only tracks registry to be rejected")
	}
	if res.Reason != "required context incomplete" {
		t.Fatalf("reason=%q want required context incomplete", res.Reason)
	}
}
