package hooks

import (
	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/pkg/api"
)

func Handle(projectDir string, evt api.HookEvent) api.HookResult {
	switch evt.Event {
	case "SessionStart":
		return api.HookResult{Decision: "allow", Metadata: SessionStartData(projectDir)}
	case "UserPromptSubmit":
		return api.HookResult{Decision: "allow"}
	case "PreToolUse":
		return PreToolUse(projectDir, evt)
	case "PostToolUse":
		return PostToolUse()
	case "Stop", "SubagentStop":
		return StopDecision(ctxpkg.Complete(projectDir))
	default:
		return api.HookResult{Decision: "allow"}
	}
}
