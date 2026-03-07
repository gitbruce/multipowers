package app_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	app "github.com/gitbruce/multipowers/internal/app"
	"github.com/gitbruce/multipowers/internal/cli"
	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

func TestAdmissionLifecycleBlockPlanThenContinueInWorktree(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	blocked := runJSONCommandAllowErrorApp(t, []string{"develop", "--dir", d, "--prompt", "refactor the entire authentication flow", "--json"})
	if blocked.Status != "blocked" {
		t.Fatalf("expected initial block, got %+v", blocked)
	}
	trackID, _ := blocked.Data["track_id"].(string)
	if trackID == "" {
		t.Fatalf("expected track_id after block, got %+v", blocked.Data)
	}

	planned := runJSONCommandApp(t, []string{"plan", "--dir", d, "--prompt", "plan the runtime changes", "--json"})
	gotTrackID, _ := planned.Data["track_id"].(string)
	if gotTrackID != trackID {
		t.Fatalf("expected plan to reuse %q, got %q", trackID, gotTrackID)
	}

	meta, err := tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.ComplexityScore <= 0 {
		t.Fatalf("expected plan to record complexity score, got %+v", meta)
	}
	if !meta.WorktreeRequired {
		t.Fatalf("expected plan to require worktree for high complexity, got %+v", meta)
	}

	wt := filepath.Join(t.TempDir(), "wt")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: /tmp/fake-linked-worktree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(d, ".multipowers"), filepath.Join(wt, ".multipowers")); err != nil {
		t.Fatal(err)
	}

	continued := runJSONCommandApp(t, []string{"develop", "--dir", wt, "--prompt", "refactor the entire authentication flow", "--json"})
	if continued.Status != "ok" {
		t.Fatalf("expected worktree execution to continue, got %+v", continued)
	}
	continuedTrackID, _ := continued.Data["track_id"].(string)
	if continuedTrackID != trackID {
		t.Fatalf("expected continued command to stay on %q, got %q", trackID, continuedTrackID)
	}
}

func TestAdmissionLifecycleLegacySpecPlanDoNotSatisfyPlanning(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	trackID := "track-legacy"
	if err := tracks.SetActiveTrack(d, trackID); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteMetadata(d, trackID, tracks.Metadata{ID: trackID, ComplexityScore: 9, WorktreeRequired: true}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"spec.md", "plan.md"} {
		path := filepath.Join(tracks.Dir(d, trackID), name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	resp := app.RunSpecPipeline(d, true, []string{"develop", "all"}, "refactor the entire authentication flow", func() api.Response {
		return api.Response{Status: "ok"}
	})
	if resp.Status != "blocked" {
		t.Fatalf("expected legacy planning shape to block, got %+v", resp)
	}
	if got, _ := resp.Data["requires_planning"].(bool); !got {
		t.Fatalf("expected requires_planning=true, got %+v", resp.Data)
	}
	missing, _ := resp.Data["missing_artifacts"].([]string)
	if len(missing) == 0 {
		t.Fatalf("expected canonical missing_artifacts, got %+v", resp.Data)
	}
}

func runJSONCommandApp(t *testing.T, args []string) api.Response {
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

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("Run(%v) exit=%d output=%s", args, exitCode, strings.TrimSpace(string(output)))
	}

	var resp api.Response
	if err := json.Unmarshal(output, &resp); err != nil {
		t.Fatalf("invalid JSON output: %v; output=%s", err, string(output))
	}
	return resp
}

func runJSONCommandAllowErrorApp(t *testing.T, args []string) api.Response {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w

	_ = cli.Run(args)

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

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
