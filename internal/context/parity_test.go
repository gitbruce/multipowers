package context

import "testing"

func TestGuardParityRequiredFiles(t *testing.T) {
	d := t.TempDir()
	if Complete(d) {
		t.Fatalf("context should be incomplete before init")
	}
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if !Complete(d) {
		t.Fatalf("context should be complete after init")
	}
}
