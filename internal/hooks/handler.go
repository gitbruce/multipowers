package hooks

import (
	"strings"

	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/pkg/api"
)

func isSpecPrompt(evt api.HookEvent) bool {
	if evt.ToolInput == nil {
		return false
	}
	raw, _ := evt.ToolInput["prompt"].(string)
	raw = strings.ToLower(strings.TrimSpace(raw))
	return strings.HasPrefix(raw, "/octo:plan") ||
		strings.HasPrefix(raw, "/octo:discover") ||
		strings.HasPrefix(raw, "/octo:define") ||
		strings.HasPrefix(raw, "/octo:develop") ||
		strings.HasPrefix(raw, "/octo:deliver") ||
		strings.HasPrefix(raw, "/octo:embrace") ||
		strings.HasPrefix(raw, "/octo:review") ||
		strings.HasPrefix(raw, "/octo:research") ||
		strings.HasPrefix(raw, "/octo:debate")
}

func Handle(projectDir string, evt api.HookEvent) api.HookResult {
	switch evt.Event {
	case "SessionStart":
		return api.HookResult{Decision: "allow", Metadata: SessionStartData(projectDir)}
	case "UserPromptSubmit":
		if isSpecPrompt(evt) && !ctxpkg.Complete(projectDir) {
			return api.HookResult{Decision: "block", Reason: "missing required .multipowers context", Remediation: "run /octo:init first"}
		}
		return api.HookResult{Decision: "allow"}
	case "PreToolUse":
		return PreToolUse(projectDir, evt)
	case "PostToolUse":
		return PostToolUse(projectDir, evt)
	case "Stop", "SubagentStop":
		return StopDecision(ctxpkg.Complete(projectDir))
	default:
		return api.HookResult{Decision: "allow"}
	}
}
