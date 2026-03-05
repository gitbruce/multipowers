package cli

import (
	"testing"
)

func TestGetRuntimeStatus_NoContext(t *testing.T) {
	d := t.TempDir()
	status := GetRuntimeStatus(d)

	if status.ContextComplete {
		t.Error("expected context to be incomplete for empty directory")
	}
	if status.Ready {
		t.Error("expected status to not be ready without context")
	}
	if status.Status == "ready" {
		t.Error("expected status to not be 'ready' without context")
	}
}

func TestGetRuntimeStatus_HasProviders(t *testing.T) {
	d := t.TempDir()
	status := GetRuntimeStatus(d)

	// Providers count should be >= 0
	if status.ProvidersCount < 0 {
		t.Error("providers count should not be negative")
	}
}

func TestGetRuntimeStatus_HookReady(t *testing.T) {
	d := t.TempDir()
	status := GetRuntimeStatus(d)

	if !status.HookReady {
		t.Error("hooks should always be ready (they're part of runtime)")
	}
}

func TestGetRuntimeStatus_HasHookEvents(t *testing.T) {
	d := t.TempDir()
	status := GetRuntimeStatus(d)

	if len(status.HookEvents) == 0 {
		t.Error("expected non-empty hook events list")
	}
	expectedEvents := map[string]bool{
		"SessionStart":     true,
		"EnterPlanMode":    true,
		"UserPromptSubmit": true,
		"PreToolUse":       true,
		"PostToolUse":      true,
		"WorktreeCreate":   true,
		"WorktreeRemove":   true,
		"Stop":             true,
		"SubagentStop":     true,
	}
	for _, event := range status.HookEvents {
		if !expectedEvents[event] {
			t.Errorf("unexpected hook event: %s", event)
		}
	}
}

func TestGetRuntimeStatus_ValidationStatus(t *testing.T) {
	d := t.TempDir()
	status := GetRuntimeStatus(d)

	// Without .multipowers, validation should fail
	if status.ValidationStatus == "passed" {
		t.Error("expected validation to fail without .multipowers")
	}
}

func TestRuntimeStatus_JSONFields(t *testing.T) {
	status := RuntimeStatus{
		Status:             "ready",
		Ready:              true,
		ContextComplete:    true,
		ProvidersAvailable: []string{"claude"},
		ProvidersCount:     1,
		HookReady:          true,
	}

	if status.Status != "ready" {
		t.Errorf("expected status ready, got %s", status.Status)
	}
	if !status.Ready {
		t.Error("expected ready to be true")
	}
	if status.ProvidersCount != 1 {
		t.Errorf("expected 1 provider, got %d", status.ProvidersCount)
	}
}
