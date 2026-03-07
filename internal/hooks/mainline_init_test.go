package hooks

import (
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestMainlineInitGate_BlocksBrainstormWithoutContext(t *testing.T) {
	r := Handle(t.TempDir(), api.HookEvent{Event: "UserPromptSubmit", ToolInput: map[string]any{"prompt": "/mp:brainstorm mainline reset"}})
	if r.Decision != "block" {
		t.Fatalf("expected block, got %+v", r)
	}
	if got := r.Metadata["recommended_command"]; got != "/mp:init" {
		t.Fatalf("recommended_command=%v want /mp:init", got)
	}
}

func TestMainlineInitGate_BlocksExecuteWithoutContext(t *testing.T) {
	r := Handle(t.TempDir(), api.HookEvent{Event: "UserPromptSubmit", ToolInput: map[string]any{"prompt": "/mp:execute ship the branch"}})
	if r.Decision != "block" {
		t.Fatalf("expected block, got %+v", r)
	}
	if got := r.Metadata["action_required"]; got != "run_init" {
		t.Fatalf("action_required=%v want run_init", got)
	}
}
