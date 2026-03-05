package hooks

import (
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestHookStopAndSubagentStop(t *testing.T) {
	d := t.TempDir()
	r := Handle(d, api.HookEvent{Event: "Stop"})
	if r.Decision != "block" {
		t.Fatalf("expected block before init")
	}
	r = Handle(d, api.HookEvent{Event: "SubagentStop"})
	if r.Decision != "block" {
		t.Fatalf("expected block before init")
	}
}
