package validation

import (
	"os"
	"path/filepath"
	"strings"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/tracks"
)

type Result struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

func EnsureTargetWorkspace(projectDir string) Result {
	if _, err := os.Stat(filepath.Join(projectDir, ".multipowers")); err != nil {
		return Result{Valid: false, Reason: "missing .multipowers workspace"}
	}
	if !ctxpkg.Complete(projectDir) {
		return Result{Valid: false, Reason: "required context incomplete"}
	}
	return Result{Valid: true}
}

func EnsureTrackExecution(projectDir string) Result {
	active, err := tracks.ActiveTrack(projectDir)
	if err != nil {
		return Result{Valid: false, Reason: err.Error()}
	}
	if strings.TrimSpace(active) == "" {
		return Result{Valid: true}
	}

	meta, err := tracks.ReadMetadata(projectDir, active)
	if err != nil {
		return Result{Valid: false, Reason: err.Error()}
	}
	if meta.GroupStatus != tracks.GroupStatusInProgress {
		return Result{Valid: true}
	}

	missing := make([]string, 0, 3)
	if strings.TrimSpace(meta.LastCommitSHA) == "" {
		missing = append(missing, "last_commit_sha")
	}
	if strings.TrimSpace(meta.LastVerifiedAt) == "" {
		missing = append(missing, "last_verified_at")
	}
	if len(meta.CompletedGroups) == 0 {
		missing = append(missing, "completed_groups")
	}
	if len(missing) > 0 {
		return Result{Valid: false, Reason: "track group enforcement incomplete: missing " + strings.Join(missing, ", ")}
	}
	return Result{Valid: true}
}
