package context

import "testing"

func TestGuardParityRequiredFiles(t *testing.T) {
	d := t.TempDir()
	if Complete(d) {
		t.Fatalf("context should be incomplete before init")
	}
	if err := RunInit(d); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if !Complete(d) {
		t.Fatalf("context should be complete after init")
	}
}
