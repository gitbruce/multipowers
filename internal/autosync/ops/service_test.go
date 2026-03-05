package ops

import (
	"testing"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/autosync/overlay"
)

func TestPolicySync_DefaultDryRun(t *testing.T) {
	d := t.TempDir()
	svc := NewService(d)
	res, err := svc.Sync(SyncOptions{})
	if err != nil {
		t.Fatalf("sync err: %v", err)
	}
	if !res.DryRun {
		t.Fatal("expected dry-run by default")
	}
}

func TestPolicySync_ApplyIgnoreRollbackRevoke(t *testing.T) {
	d := t.TempDir()
	_, err := overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("seed overlay: %v", err)
	}
	svc := NewService(d)
	if _, err := svc.Sync(SyncOptions{Apply: true}); err != nil {
		t.Fatalf("apply err: %v", err)
	}
	if _, err := svc.Sync(SyncOptions{IgnoreID: "p1"}); err != nil {
		t.Fatalf("ignore err: %v", err)
	}
	if _, err := svc.Sync(SyncOptions{RollbackID: "p1"}); err != nil {
		t.Fatalf("rollback err: %v", err)
	}
	if _, err := svc.Sync(SyncOptions{RevokeID: "r1"}); err != nil {
		t.Fatalf("revoke err: %v", err)
	}
}

func TestPolicyStatsAndTune(t *testing.T) {
	d := t.TempDir()
	svc := NewService(d)
	if _, err := svc.Stats(); err != nil {
		t.Fatalf("stats err: %v", err)
	}
	if err := svc.Tune("balanced"); err != nil {
		t.Fatalf("tune err: %v", err)
	}
}
