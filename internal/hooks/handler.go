package hooks

import (
	"strings"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/benchmark"
	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/isolation"
	"github.com/gitbruce/multipowers/internal/orchestration"
	"github.com/gitbruce/multipowers/internal/policy"
	"github.com/gitbruce/multipowers/pkg/api"
)

func isSpecPrompt(evt api.HookEvent) bool {
	if evt.ToolInput == nil {
		return false
	}
	raw, _ := evt.ToolInput["prompt"].(string)
	raw = strings.ToLower(strings.TrimSpace(raw))
	return strings.HasPrefix(raw, "/mp:brainstorm") ||
		strings.HasPrefix(raw, "/mp:design") ||
		strings.HasPrefix(raw, "/mp:plan") ||
		strings.HasPrefix(raw, "/mp:execute") ||
		strings.HasPrefix(raw, "/mp:debug") ||
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
			"resume_prompt":       rawPrompt,
			"missing_files":       strings.Join(missing, ","),
		},
	}
}

func handleEnterPlanMode(evt api.HookEvent) api.HookResult {
	raw := ""
	if evt.ToolInput != nil {
		if s, ok := evt.ToolInput["prompt"].(string); ok {
			raw = strings.TrimSpace(strings.ToLower(s))
		}
	}
	if strings.HasPrefix(raw, "/mp:plan") {
		return api.HookResult{
			Decision: "allow",
			Reason:   "plan-mode intent confirmed via /mp:plan",
		}
	}
	return api.HookResult{
		Decision:    "block",
		Reason:      "plan mode requires explicit /mp:plan intent",
		Remediation: "run /mp:plan <goal> before entering plan mode",
		Metadata: map[string]any{
			"required_command": "/mp:plan",
		},
	}
}

func Handle(projectDir string, evt api.HookEvent) api.HookResult {
	_, _ = autosync.EmitRawEvent(projectDir, "hook", evt.Event, map[string]any{
		"tool_name": evt.ToolName,
	})

	switch evt.Event {
	case "SessionStart":
		return api.HookResult{Decision: "allow", Metadata: SessionStartData(projectDir)}
	case "EnterPlanMode":
		return handleEnterPlanMode(evt)
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
	case "WorktreeCreate", "WorktreeRemove":
		path, err := appendWorktreeEvent(projectDir, evt, nil)
		if err != nil {
			return api.HookResult{
				Decision: "allow",
				Reason:   "worktree event accepted (persistence failed)",
				Metadata: map[string]any{"worktree_event_error": err.Error()},
			}
		}
		return api.HookResult{
			Decision: "allow",
			Reason:   "worktree event persisted",
			Metadata: map[string]any{"worktree_events_log": path},
		}
	case "Stop", "SubagentStop":
		return StopDecision(projectDir, evt.Event, ctxpkg.Complete(projectDir))
	default:
		return api.HookResult{Decision: "allow"}
	}
}

// resolveModelRouting resolves model routing using the compiled policy resolver.
func resolveModelRouting(projectDir, prompt string) map[string]any {
	workflowName := extractWorkflowName(prompt)
	intent := benchmark.ClassifyCodeIntent(benchmark.IntentRequest{
		WhitelistHits:       extractIntentWhitelistHits(prompt),
		HasLLMSemantic:      false,
		LLMDecisionPriority: true,
	})
	isolationMetadata := resolveExecutionIsolation(projectDir, workflowName, intent.CodeRelated)
	resolver, err := policy.NewResolverFromProjectDir(projectDir)
	if err != nil {
		return map[string]any{
			"model_routing_error": err.Error(),
			"model_routing": map[string]any{
				"command": workflowName,
			},
			"benchmark_code_intent": map[string]any{
				"code_related": intent.CodeRelated,
				"source":       intent.Source,
				"whitelist":    intent.WhitelistHits,
			},
			"execution_isolation": isolationMetadata,
		}
	}

	if workflowName == "" {
		return map[string]any{
			"model_routing_error": "workflow name missing from prompt",
			"benchmark_code_intent": map[string]any{
				"code_related": intent.CodeRelated,
				"source":       intent.Source,
				"whitelist":    intent.WhitelistHits,
			},
			"execution_isolation": isolationMetadata,
		}
	}

	contract, err := resolver.Resolve(policy.ResolveRequest{
		Scope: policy.ScopeWorkflow,
		Name:  workflowName,
	})
	if err != nil {
		return map[string]any{
			"model_routing_error": err.Error(),
			"model_routing": map[string]any{
				"command": workflowName,
			},
			"benchmark_code_intent": map[string]any{
				"code_related": intent.CodeRelated,
				"source":       intent.Source,
				"whitelist":    intent.WhitelistHits,
			},
			"execution_isolation": isolationMetadata,
		}
	}

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
		"benchmark_code_intent": map[string]any{
			"code_related": intent.CodeRelated,
			"source":       intent.Source,
			"whitelist":    intent.WhitelistHits,
		},
		"execution_isolation": isolationMetadata,
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

func extractIntentWhitelistHits(prompt string) []string {
	p := strings.ToLower(prompt)
	hits := make([]string, 0, 8)
	seen := map[string]struct{}{}

	add := func(tag string) {
		if _, ok := seen[tag]; ok {
			return
		}
		seen[tag] = struct{}{}
		hits = append(hits, tag)
	}

	for token, tag := range map[string]string{
		"/mp:brainstorm": "task_type:brainstorm",
		"/mp:design":     "task_type:design",
		"/mp:execute":    "task_type:execute",
		"/mp:debug":      "task_type:debug",
		"/mp:plan":       "task_type:plan",
		"/mp:debate":     "task_type:debate",
		"go":             "language:go",
		"python":         "language:python",
		"typescript":     "language:typescript",
		"javascript":     "language:javascript",
		"java":           "language:java",
		"rust":           "language:rust",
		"react":          "framework:react",
		"vue":            "framework:vue",
		"angular":        "framework:angular",
		"next.js":        "framework:nextjs",
		"django":         "framework:django",
		"flask":          "framework:flask",
		"spring":         "framework:spring",
		"api":            "tech:api",
		"endpoint":       "tech:endpoint",
		"database":       "tech:database",
		"sql":            "tech:sql",
		"schema":         "tech:schema",
		"function":       "tech:function",
		"class":          "tech:class",
		"test":           "tech:test",
		"benchmark":      "tech:benchmark",
	} {
		if strings.Contains(p, token) {
			add(tag)
		}
	}
	return hits
}

func resolveExecutionIsolation(projectDir, workflowName string, codeRelated bool) map[string]any {
	cfg, err := orchestration.LoadConfigFromProjectDir(projectDir)
	if err != nil {
		return map[string]any{
			"enforced": false,
			"reason":   "config_load_error",
			"error":    err.Error(),
		}
	}

	decision := isolation.ResolveExternalCommandIsolation(isolation.ExternalCommandIsolationInput{
		IsolationEnabled: cfg.ExecutionIsolation.Enabled,
		ExternalCommand:  strings.TrimSpace(workflowName) != "",
		MayEditFiles:     mayEditFilesForWorkflow(workflowName),
		CodeRelated:      codeRelated,
		Command:          workflowName,
		CommandWhitelist: cfg.ExecutionIsolation.CommandWhitelist,
		BenchmarkProfile: isolation.BenchmarkProfileInput{
			Enabled:           cfg.BenchmarkMode.ExecutionProfile.Enabled,
			RequireCodeIntent: cfg.BenchmarkMode.ExecutionProfile.RequireCodeIntent,
			CommandWhitelist:  cfg.BenchmarkMode.ExecutionProfile.CommandWhitelist,
		},
	})

	return map[string]any{
		"enforced":                decision.Enforced,
		"reason":                  decision.Reason,
		"shared_whitelist_match":  decision.SharedWhitelistMatch,
		"profile_whitelist_match": decision.ProfileWhitelistMatch,
		"command":                 strings.TrimSpace(workflowName),
		"code_related":            codeRelated,
		"may_edit_files":          mayEditFilesForWorkflow(workflowName),
		"branch_prefix":           cfg.ExecutionIsolation.BranchPrefix,
		"worktree_root":           cfg.ExecutionIsolation.WorktreeRoot,
	}
}

func mayEditFilesForWorkflow(workflowName string) bool {
	switch strings.ToLower(strings.TrimSpace(workflowName)) {
	case "develop", "review", "embrace", "deliver", "debug", "tdd", "security":
		return true
	default:
		return false
	}
}
