package app

import (
	"os"
	"path/filepath"
	"testing"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

func TestRunSpecPipelineHighComplexityMissingPlanReturnsStructuredBlock(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	called := false
	r := RunSpecPipeline(d, true, []string{"develop", "all"}, "refactor the entire authentication flow", func() api.Response {
		called = true
		return api.Response{Status: "ok"}
	})
	if called {
		t.Fatal("expected pipeline to block before execFn")
	}
	if r.Status != "blocked" {
		t.Fatalf("expected blocked response, got %+v", r)
	}
	if r.Action != "ask_user_questions" {
		t.Fatalf("action=%q want ask_user_questions", r.Action)
	}
	if got, _ := r.Data["requires_planning"].(bool); !got {
		t.Fatalf("expected requires_planning=true, got %+v", r.Data)
	}
	if got, _ := r.Data["requires_worktree"].(bool); !got {
		t.Fatalf("expected requires_worktree=true, got %+v", r.Data)
	}
}

func TestRunSpecPipelineHighComplexityMissingWorktreeReturnsStructuredBlock(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
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

	r := RunSpecPipeline(d, true, []string{"develop", "all"}, "refactor the entire authentication flow", func() api.Response {
		return api.Response{Status: "ok"}
	})
	if r.Status != "blocked" {
		t.Fatalf("expected blocked response, got %+v", r)
	}
	if got, _ := r.Data["requires_planning"].(bool); got {
		t.Fatalf("expected requires_planning=false, got %+v", r.Data)
	}
	if got, _ := r.Data["requires_worktree"].(bool); !got {
		t.Fatalf("expected requires_worktree=true, got %+v", r.Data)
	}
}

func TestRunSpecPipelineLowComplexityContinues(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	called := false
	r := RunSpecPipeline(d, true, []string{"develop", "all"}, "update readme typo", func() api.Response {
		called = true
		return api.Response{Status: "ok"}
	})
	if !called {
		t.Fatal("expected low-complexity pipeline to continue")
	}
	if r.Status != "ok" {
		t.Fatalf("expected ok response, got %+v", r)
	}
}
