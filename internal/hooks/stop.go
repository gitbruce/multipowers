package hooks

import (
	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/decisions"
	"github.com/gitbruce/multipowers/pkg/api"
)

func StopDecision(projectDir, source string, canStop bool) api.HookResult {
	_, _ = autosync.EmitRawEvent(projectDir, "hook.stop", source, map[string]any{
		"can_stop": canStop,
	})
	if canStop {
		return api.HookResult{Decision: "allow", Reason: "no mandatory checkpoint pending"}
	}
	_ = decisions.AppendQualityGate(projectDir, source, "mandatory checkpoint pending", "session-stop")
	return api.HookResult{Decision: "block", Reason: "mandatory checkpoint pending", Remediation: "finish required init/context workflow before stop"}
}
