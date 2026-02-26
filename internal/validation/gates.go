package validation

import (
	"os"
	"path/filepath"

	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
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
