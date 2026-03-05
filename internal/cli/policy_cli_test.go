package cli

import "testing"

func TestPolicySync_DefaultDryRun(t *testing.T) {
	d := t.TempDir()
	code := Run([]string{"policy", "sync", "--dir", d, "--json"})
	if code != 0 {
		t.Fatalf("expected rc=0 got %d", code)
	}
}

func TestPolicySync_ApplyIgnoreRollbackRevoke(t *testing.T) {
	d := t.TempDir()
	if code := Run([]string{"policy", "sync", "--dir", d, "--apply", "--json"}); code != 0 {
		t.Fatalf("apply rc=%d", code)
	}
	if code := Run([]string{"policy", "sync", "--dir", d, "--ignore-id", "p1", "--json"}); code != 0 {
		t.Fatalf("ignore rc=%d", code)
	}
	if code := Run([]string{"policy", "sync", "--dir", d, "--rollback-id", "p1", "--json"}); code != 0 {
		t.Fatalf("rollback rc=%d", code)
	}
	if code := Run([]string{"policy", "sync", "--dir", d, "--revoke-id", "r1", "--json"}); code != 0 {
		t.Fatalf("revoke rc=%d", code)
	}
}

func TestPolicyStatsAndTune(t *testing.T) {
	d := t.TempDir()
	if code := Run([]string{"policy", "stats", "--dir", d, "--json"}); code != 0 {
		t.Fatalf("stats rc=%d", code)
	}
	if code := Run([]string{"policy", "tune", "--dir", d, "--mode", "balanced", "--json"}); code != 0 {
		t.Fatalf("tune rc=%d", code)
	}
}
