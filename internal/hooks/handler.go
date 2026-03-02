package hooks

import (
	"strings"

	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/modelroute"
	"github.com/gitbruce/claude-octopus/pkg/api"
)

func isSpecPrompt(evt api.HookEvent) bool {
	if evt.ToolInput == nil {
		return false
	}
	raw, _ := evt.ToolInput["prompt"].(string)
	raw = strings.ToLower(strings.TrimSpace(raw))
	return strings.HasPrefix(raw, "/mp:plan") ||
		strings.HasPrefix(raw, "/mp:discover") ||
		strings.HasPrefix(raw, "/mp:define") ||
		strings.HasPrefix(raw, "/mp:develop") ||
		strings.HasPrefix(raw, "/mp:deliver") ||
		strings.HasPrefix(raw, "/mp:embrace") ||
		strings.HasPrefix(raw, "/mp:review") ||
		strings.HasPrefix(raw, "/mp:research") ||
		strings.HasPrefix(raw, "/mp:debate")
}

func missingContextGuidance(rawPrompt string, missing []string) api.HookResult {
	return api.HookResult{
		Decision:    "block",
		Reason:      "missing required .multipowers context",
		Remediation: "run /mp:init to complete guided setup, then retry the original /mp command",
		Metadata: map[string]any{
			"action_required":     "run_init",
			"recommended_command": "/mp:init",
			"resume_command":      rawPrompt,
			"missing_files":       strings.Join(missing, ","),
		},
	}
}

func Handle(projectDir string, evt api.HookEvent) api.HookResult {
	switch evt.Event {
	case "SessionStart":
		return api.HookResult{Decision: "allow", Metadata: SessionStartData(projectDir)}
	case "UserPromptSubmit":
		if isSpecPrompt(evt) && !ctxpkg.Complete(projectDir) {
			raw, _ := evt.ToolInput["prompt"].(string)
			return missingContextGuidance(raw, ctxpkg.Missing(projectDir))
		}
		if isSpecPrompt(evt) {
			raw, _ := evt.ToolInput["prompt"].(string)
			r := modelroute.ResolveForPrompt(projectDir, raw)
			return api.HookResult{
				Decision: "allow",
				Metadata: map[string]any{
					"model_routing": r,
				},
			}
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
