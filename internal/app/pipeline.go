package app

import (
	"github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/runtime"
	"github.com/gitbruce/multipowers/internal/validation"
	"github.com/gitbruce/multipowers/pkg/api"
)

type ExecFunc func() api.Response

func RunSpecPipeline(projectDir string, autoInit bool, tags []string, execFn ExecFunc) api.Response {
	missing := context.Missing(projectDir)
	if len(missing) > 0 {
		_ = autoInit // kept for CLI compatibility, but no silent init is allowed.
		return api.Response{
			Status:      "blocked",
			Action:      "run_init",
			ErrorCode:   ErrCtxMissing,
			Missing:     missing,
			Remediation: "Run /mp:init wizard first; context files are never generated silently.",
		}
	}
	cfg, present, err := runtime.Load(projectDir)
	if err != nil {
		return api.Response{Status: "error", ErrorCode: ErrPreRunFailed, Message: err.Error()}
	}
	if present {
		if err := runtime.RunPreRun(cfg, tags); err != nil {
			return api.Response{Status: "error", ErrorCode: ErrPreRunFailed, Message: err.Error()}
		}
	}
	enforcement := validation.EnsureTrackExecution(projectDir)
	if !enforcement.Valid {
		return api.Response{Status: "blocked", ErrorCode: ErrInvalidArgument, Message: enforcement.Reason}
	}
	return execFn()
}
