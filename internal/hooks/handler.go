package hooks

import (
	"strings"

	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/modelroute"
	"github.com/gitbruce/claude-octopus/internal/policy"
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
			// Use new policy resolver first, fall back to legacy modelroute
			metadata := resolveModelRouting(projectDir, raw)
			return api.HookResult{
				Decision: "allow",
				Metadata: metadata,
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

// resolveModelRouting resolves model routing using the new policy resolver
// with fallback to legacy modelroute for backward compatibility
func resolveModelRouting(projectDir, prompt string) map[string]any {
	// Try new policy resolver first
	resolver, err := policy.NewResolverFromProjectDir(projectDir)
	if err == nil {
		// Parse workflow name from prompt
		workflowName := extractWorkflowName(prompt)
		if workflowName != "" {
			contract, err := resolver.Resolve(policy.ResolveRequest{
				Scope: policy.ScopeWorkflow,
				Name:  workflowName,
			})
			if err == nil {
				return map[string]any{
					"model_routing": map[string]any{
						"command":            workflowName,
						"model":              contract.RequestedModel,
						"provider":           string(contract.ExecutorKind),
						"executor_profile":   contract.ExecutorProfile,
						"enforcement":        string(contract.Enforcement),
						"fallback_target":    contract.FallbackTarget,
						"source":             contract.SourceRef,
						"resolved_by_policy": true,
					},
				}
			}
		}
	}

	// Fall back to legacy modelroute
	r := modelroute.ResolveForPrompt(projectDir, prompt)
	return map[string]any{
		"model_routing": r,
	}
}

// extractWorkflowName extracts the workflow name from a /mp: command
func extractWorkflowName(prompt string) string {
	p := strings.ToLower(strings.TrimSpace(prompt))
	if !strings.HasPrefix(p, "/mp:") {
		return ""
	}
	p = strings.TrimPrefix(p, "/mp:")
	if i := strings.IndexAny(p, " \t\n"); i >= 0 {
		return p[:i]
	}
	return p
}
