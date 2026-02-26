package hooks

import "github.com/gitbruce/claude-octopus/pkg/api"

func StopDecision(canStop bool) api.HookResult {
	if canStop {
		return api.HookResult{Decision: "allow", Reason: "no mandatory checkpoint pending"}
	}
	return api.HookResult{Decision: "block", Reason: "mandatory checkpoint pending", Remediation: "finish required init/context workflow before stop"}
}
