package autosync_test

import (
	"testing"
	"time"

	autosync "github.com/gitbruce/multipowers/internal/autosync"
	autosyncctx "github.com/gitbruce/multipowers/internal/autosync/context"
	"github.com/gitbruce/multipowers/internal/autosync/overlay"
)

func TestE2E_AutoApplyInjectRevokeResetFlow(t *testing.T) {
	d := t.TempDir()
	_, err := overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo", Confidence: 0.99})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	ctx1, err := autosyncctx.BuildPolicyContext(d, "", time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build context: %v", err)
	}
	if len(ctx1.ActiveRules) == 0 {
		t.Fatal("expected active rule after apply")
	}
	if err := overlay.SuppressForSession(d, "s1", "r1"); err != nil {
		t.Fatalf("suppress: %v", err)
	}
	ctx2, err := autosyncctx.BuildPolicyContext(d, "s1", time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build context s1: %v", err)
	}
	if len(ctx2.ActiveRules) != 0 {
		t.Fatalf("expected suppressed session to inject no rules: %+v", ctx2.ActiveRules)
	}
	now := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	if err := overlay.RevokeRule(d, "r1", "user delete", now, 24*time.Hour); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	ctx3, err := autosyncctx.BuildPolicyContext(d, "", now)
	if err != nil {
		t.Fatalf("build context after revoke: %v", err)
	}
	if len(ctx3.ActiveRules) != 0 {
		t.Fatalf("expected no rules during cooldown: %+v", ctx3.ActiveRules)
	}
	_, err = overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo", Confidence: 0.99})
	if err != nil {
		t.Fatalf("re-apply: %v", err)
	}
	ctx4, err := autosyncctx.BuildPolicyContext(d, "", now.Add(25*time.Hour))
	if err != nil {
		t.Fatalf("build context after cooldown: %v", err)
	}
	if len(ctx4.ActiveRules) == 0 {
		t.Fatal("expected relearned rule after cooldown reset")
	}
}
