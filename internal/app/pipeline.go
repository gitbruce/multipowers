package app

import (
	"github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/runtime"
	"github.com/gitbruce/claude-octopus/pkg/api"
)

type ExecFunc func() api.Response

func RunSpecPipeline(projectDir string, autoInit bool, tags []string, execFn ExecFunc) api.Response {
	missing := context.Missing(projectDir)
	if len(missing) > 0 {
		if !autoInit {
			return api.Response{Status: "blocked", Action: "run_init", ErrorCode: ErrCtxMissing, Missing: missing}
		}
		if err := context.RunInit(projectDir); err != nil {
			return api.Response{Status: "error", ErrorCode: ErrInitFailed, Message: err.Error(), Remediation: "Run /mp:init interactively."}
		}
		if missing = context.Missing(projectDir); len(missing) > 0 {
			return api.Response{Status: "error", ErrorCode: ErrCtxMissing, Missing: missing}
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
	return execFn()
}
