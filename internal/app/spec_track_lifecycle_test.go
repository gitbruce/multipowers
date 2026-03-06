package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

func TestSpecTrackLifecycle(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	planResp := RunSpecPipeline(d, true, []string{"plan", "all"}, func() api.Response {
		trackCtx, err := prepareLifecycleTrack(d, "plan", "design runtime lifecycle")
		if err != nil {
			return api.Response{Status: "error", Message: err.Error()}
		}
		return api.Response{Status: "ok", Data: map[string]any{"track_id": trackCtx.ID}}
	})
	if planResp.Status != "ok" {
		t.Fatalf("plan pipeline=%+v", planResp)
	}
	trackID, _ := planResp.Data["track_id"].(string)
	if trackID == "" {
		t.Fatalf("expected plan track_id, got %+v", planResp.Data)
	}

	developResp := RunSpecPipeline(d, true, []string{"develop", "all"}, func() api.Response {
		trackCtx, err := prepareLifecycleTrack(d, "develop", "implement runtime lifecycle")
		if err != nil {
			return api.Response{Status: "error", Message: err.Error()}
		}
		return api.Response{Status: "ok", Data: map[string]any{"track_id": trackCtx.ID}}
	})
	if developResp.Status != "ok" {
		t.Fatalf("develop pipeline=%+v", developResp)
	}
	gotTrackID, _ := developResp.Data["track_id"].(string)
	if gotTrackID != trackID {
		t.Fatalf("expected active track reuse %q, got %q", trackID, gotTrackID)
	}

	if _, err := os.Stat(filepath.Join(d, ".multipowers", "context", "runtime.json")); err != nil {
		t.Fatalf("runtime.json missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "tracks", "tracks.md")); err != nil {
		t.Fatalf("canonical tracks registry missing: %v", err)
	}
	for _, name := range []string{"intent.md", "design.md", "implementation-plan.md", "metadata.json", "index.md"} {
		path := filepath.Join(d, ".multipowers", "tracks", trackID, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("artifact %s missing: %v", name, err)
		}
	}
}

func prepareLifecycleTrack(projectDir, command, prompt string) (tracks.TrackContext, error) {
	coordinator := tracks.TrackCoordinator{}
	trackCtx, err := coordinator.ResolveTrack(projectDir, command)
	if err != nil {
		return tracks.TrackContext{}, err
	}
	meta, err := tracks.ReadMetadata(projectDir, trackCtx.ID)
	if err != nil {
		return tracks.TrackContext{}, err
	}
	values := tracks.DefaultArtifactValues(trackCtx, command, prompt)
	values["CompletedGroups"] = appendLifecycleGroup(meta.CompletedGroups, command)
	if strings.TrimSpace(meta.Title) != "" {
		values["TrackTitle"] = meta.Title
	}
	if err := coordinator.EnsureArtifacts(projectDir, trackCtx, values); err != nil {
		return tracks.TrackContext{}, err
	}
	if err := tracks.UpdateMetadata(projectDir, trackCtx.ID, func(current *tracks.Metadata) error {
		if strings.TrimSpace(current.Title) == "" {
			current.Title = fmt.Sprint(values["TrackTitle"])
		}
		current.Status = "in_progress"
		current.ExecutionMode = "cli"
		current.CurrentGroup = command
		current.CompletedGroups = appendLifecycleGroup(current.CompletedGroups, command)
		current.LastVerifiedAt = time.Now().UTC().Format(time.RFC3339)
		return nil
	}); err != nil {
		return tracks.TrackContext{}, err
	}
	if err := coordinator.UpdateRegistry(projectDir, trackCtx); err != nil {
		return tracks.TrackContext{}, err
	}
	return trackCtx, nil
}

func appendLifecycleGroup(existing []string, next string) []string {
	next = strings.TrimSpace(next)
	if next == "" {
		return append([]string(nil), existing...)
	}
	out := make([]string, 0, len(existing)+1)
	seen := map[string]struct{}{}
	for _, item := range existing {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	if _, ok := seen[next]; !ok {
		out = append(out, next)
	}
	return out
}
