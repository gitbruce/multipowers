package hooks

import (
	"fmt"

	"github.com/gitbruce/multipowers/internal/decisions"
	"github.com/gitbruce/multipowers/internal/fsboundary"
	"github.com/gitbruce/multipowers/pkg/api"
)

func PreToolUse(projectDir string, evt api.HookEvent) api.HookResult {
	if evt.ToolName == "Write" || evt.ToolName == "Edit" || evt.ToolName == "MultiEdit" {
		if p, ok := evt.ToolInput["file_path"].(string); ok {
			if err := fsboundary.ValidateWritePath(p, projectDir); err != nil {
				_ = decisions.AppendQualityGate(projectDir, "PreToolUse", fmt.Sprintf("boundary violation: %v", err), "write-path")
				return api.HookResult{Decision: "block", Reason: fmt.Sprintf("boundary violation: %v", err)}
			}
		}
	}
	return api.HookResult{Decision: "allow"}
}
