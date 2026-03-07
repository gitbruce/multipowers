package app_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	app "github.com/gitbruce/multipowers/internal/app"
	"github.com/gitbruce/multipowers/internal/cli"
	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

func TestSpecTrackGroupLifecycleRequiresCompletionEvidenceBetweenGroups(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	planResp := app.RunSpecPipeline(d, true, []string{"plan", "all"}, "design runtime lifecycle", func() api.Response {
		trackCtx, err := prepareIntegrationTrack(d, "plan", "design runtime lifecycle")
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

	startResp := runJSONCommand(t, []string{"track", "group-start", "--dir", d, "--track-id", trackID, "--group", "g1", "--execution-mode", "workspace", "--json"})
	if startResp.Status != "ok" {
		t.Fatalf("group-start=%+v", startResp)
	}

	blocked := app.RunSpecPipeline(d, true, []string{"develop", "all"}, "implement runtime lifecycle", func() api.Response {
		return api.Response{Status: "ok", Message: "unexpected execution"}
	})
	if blocked.Status != "blocked" {
		t.Fatalf("expected blocked develop pipeline, got %+v", blocked)
	}
	if !strings.Contains(blocked.Message, "last_commit_sha") || !strings.Contains(blocked.Message, "last_verified_at") {
		t.Fatalf("expected missing evidence details, got %+v", blocked)
	}

	meta, err := tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.CurrentGroup != "g1" {
		t.Fatalf("current_group=%q want g1", meta.CurrentGroup)
	}
	if meta.GroupStatus != tracks.GroupStatusInProgress {
		t.Fatalf("group_status=%q want in_progress", meta.GroupStatus)
	}

	registryBefore, err := os.ReadFile(filepath.Join(d, ".multipowers", "tracks", "tracks.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(registryBefore), "Current Group: g1") {
		t.Fatalf("expected registry to show active group, got:\n%s", string(registryBefore))
	}

	completeResp := runJSONCommand(t, []string{"track", "group-complete", "--dir", d, "--track-id", trackID, "--group", "g1", "--commit-sha", "abc1234", "--json"})
	if completeResp.Status != "ok" {
		t.Fatalf("group-complete=%+v", completeResp)
	}

	developResp := app.RunSpecPipeline(d, true, []string{"develop", "all"}, "implement runtime lifecycle", func() api.Response {
		trackCtx, err := prepareIntegrationTrack(d, "develop", "implement runtime lifecycle")
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

	meta, err = tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.GroupStatus != tracks.GroupStatusCompleted {
		t.Fatalf("group_status=%q want completed", meta.GroupStatus)
	}
	if meta.LastCommitSHA != "abc1234" {
		t.Fatalf("last_commit_sha=%q want abc1234", meta.LastCommitSHA)
	}
	if strings.TrimSpace(meta.LastVerifiedAt) == "" {
		t.Fatal("expected last_verified_at to be recorded")
	}

	registryAfter, err := os.ReadFile(filepath.Join(d, ".multipowers", "tracks", "tracks.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(registryAfter), "Last Commit: abc1234") {
		t.Fatalf("expected registry to show completion evidence, got:\n%s", string(registryAfter))
	}
}

func prepareIntegrationTrack(projectDir, command, prompt string) (tracks.TrackContext, error) {
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
	values["CurrentGroup"] = meta.CurrentGroup
	values["GroupStatus"] = meta.GroupStatus
	values["LastCommand"] = command
	values["LastCommandAt"] = time.Now().UTC().Format(time.RFC3339)
	values["CompletedGroups"] = append([]string(nil), meta.CompletedGroups...)
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
		if strings.TrimSpace(current.ExecutionMode) == "" {
			current.ExecutionMode = "cli"
		}
		return nil
	}); err != nil {
		return tracks.TrackContext{}, err
	}
	if err := tracks.RecordCommandTouch(projectDir, trackCtx.ID, command, fmt.Sprint(values["LastCommandAt"])); err != nil {
		return tracks.TrackContext{}, err
	}
	if err := coordinator.UpdateRegistry(projectDir, trackCtx); err != nil {
		return tracks.TrackContext{}, err
	}
	return trackCtx, nil
}

func runJSONCommand(t *testing.T, args []string) api.Response {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w

	exitCode := cli.Run(args)

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

	if exitCode != 0 {
		output, _ := io.ReadAll(r)
		t.Fatalf("Run(%v) exit=%d output=%s", args, exitCode, strings.TrimSpace(string(output)))
	}

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	var resp api.Response
	if err := json.Unmarshal(output, &resp); err != nil {
		t.Fatalf("invalid JSON output: %v; output=%s", err, string(output))
	}
	return resp
}
