package api

import "testing"

func TestHookTypes(t *testing.T) {
	e := HookEvent{Event: "SessionStart", CWD: "/tmp"}
	if e.Event == "" {
		t.Fatal("event empty")
	}
	r := HookResult{Decision: "allow"}
	if r.Decision != "allow" {
		t.Fatal("bad decision")
	}
}
