package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/multipowers/internal/tracks"
)

func TestAdmissionHighComplexityWithoutTrackRequiresPlan(t *testing.T) {
	d := t.TempDir()
	res := EnsureSpecAdmission(d, "develop", "refactor the entire authentication flow")
	if res.Valid {
		t.Fatalf("expected blocked admission, got %+v", res)
	}
	if !res.RequiresPlanning {
		t.Fatal("expected planning to be required")
	}
	if !res.RequiresWorktree {
		t.Fatal("expected worktree to be required for high complexity")
	}
}

func TestAdmissionHighComplexityWithMissingArtifactsRequiresPlan(t *testing.T) {
	d := t.TempDir()
	trackID := "track-1"
	if err := tracks.SetActiveTrack(d, trackID); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteMetadata(d, trackID, tracks.Metadata{ID: trackID, ComplexityScore: 9, WorktreeRequired: true}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"intent.md", "index.md"} {
		path := filepath.Join(tracks.Dir(d, trackID), name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	res := EnsureSpecAdmission(d, "develop", "refactor the entire authentication flow")
	if res.Valid {
		t.Fatalf("expected blocked admission, got %+v", res)
	}
	if !res.RequiresPlanning {
		t.Fatal("expected planning to remain required when artifacts are missing")
	}
	if len(res.MissingArtifacts) != 2 {
		t.Fatalf("missing_artifacts=%v want 2 entries", res.MissingArtifacts)
	}
}

func TestAdmissionHighComplexityPlannedButOutsideWorktreeRequiresIsolation(t *testing.T) {
	d := t.TempDir()
	trackID := "track-1"
	if err := tracks.SetActiveTrack(d, trackID); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteMetadata(d, trackID, tracks.Metadata{ID: trackID, ComplexityScore: 9, WorktreeRequired: true}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"intent.md", "design.md", "implementation-plan.md", "index.md"} {
		path := filepath.Join(tracks.Dir(d, trackID), name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(d, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	res := EnsureSpecAdmission(d, "develop", "refactor the entire authentication flow")
	if res.Valid {
		t.Fatalf("expected blocked admission, got %+v", res)
	}
	if res.RequiresPlanning {
		t.Fatalf("expected planning complete case to only require isolation, got %+v", res)
	}
	if !res.RequiresWorktree {
		t.Fatal("expected worktree requirement to remain true")
	}
}

func TestAdmissionLowComplexityAllowsExecution(t *testing.T) {
	d := t.TempDir()
	res := EnsureSpecAdmission(d, "develop", "update readme typo")
	if !res.Valid {
		t.Fatalf("expected low-complexity prompt to pass, got %+v", res)
	}
}
