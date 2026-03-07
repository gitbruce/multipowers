package roles

import (
	"path/filepath"
	"testing"
)

func TestRoleDispatch_MapsMainlineRolesOnly(t *testing.T) {
	dispatcher, err := NewDispatcher(filepath.Join("..", "..", "custom", "config", "mainline-surface.yaml"))
	if err != nil {
		t.Fatalf("new dispatcher: %v", err)
	}
	cases := map[string]string{
		"init":       "initializer",
		"brainstorm": "facilitator",
		"design":     "facilitator",
		"plan":       "planner",
		"execute":    "executor",
		"debug":      "debugger",
		"debate":     "debater",
		"status":     "reviewer",
	}
	for command, want := range cases {
		got, err := dispatcher.RoleForCommand(command)
		if err != nil {
			t.Fatalf("role for %s: %v", command, err)
		}
		if got != want {
			t.Fatalf("role for %s = %s, want %s", command, got, want)
		}
	}
	if _, err := dispatcher.RoleForCommand("persona"); err == nil {
		t.Fatal("expected persona to be absent from fixed-role dispatcher")
	}
}
