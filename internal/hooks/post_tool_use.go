package hooks

import "github.com/gitbruce/claude-octopus/pkg/api"

func PostToolUse() api.HookResult {
	return api.HookResult{Decision: "allow", Reason: "post-processing complete"}
}
