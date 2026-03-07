package app

import (
	"testing"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

func TestPipelineMissingContextBlocks(t *testing.T) {
	d := t.TempDir()
	r := RunSpecPipeline(d, true, []string{"all"}, func() api.Response {
		return api.Response{Status: "ok"}
	})
	if r.Status != "blocked" || r.Action != "run_init" {
		t.Fatalf("expected blocked/run_init, got %+v", r)
	}
}

func TestPipelineBlocksWhenGroupEnforcementIncomplete(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	if err := tracks.SetActiveTrack(d, "track-g5"); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteMetadata(d, "track-g5", tracks.Metadata{
		ID:              "track-g5",
		Status:          "in_progress",
		CurrentGroup:    "g5",
		GroupStatus:     tracks.GroupStatusInProgress,
		CompletedGroups: []string{"g4"},
	}); err != nil {
		t.Fatal(err)
	}

	r := RunSpecPipeline(d, true, []string{"develop", "all"}, func() api.Response {
		return api.Response{Status: "ok"}
	})
	if r.Status != "blocked" {
		t.Fatalf("expected blocked pipeline, got %+v", r)
	}
}
