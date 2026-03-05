package hooks

import (
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestHook_HardSafetyRuleStillBlocks(t *testing.T) {
	d := t.TempDir()
	r := PreToolUse(d, api.HookEvent{
		Event:    "PreToolUse",
		ToolName: "Bash",
		ToolInput: map[string]any{
			"safety_block": true,
		},
	})
	if r.Decision != "block" {
		t.Fatalf("decision=%s want block", r.Decision)
	}
}
