package context

import (
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/autosync/overlay"
)

func TestPolicyContext_IncludesActiveRulesWithReferences(t *testing.T) {
	d := t.TempDir()
	_, err := overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	ctx, err := BuildPolicyContext(d, "", time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build context: %v", err)
	}
	if len(ctx.ActiveRules) == 0 || ctx.ActiveRules[0].RuleID != "r1" {
		t.Fatalf("unexpected active rules: %+v", ctx.ActiveRules)
	}
}

func TestPolicyContext_ExcludesRevokedRules(t *testing.T) {
	d := t.TempDir()
	_, err := overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	now := time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC)
	if err := overlay.RevokeRule(d, "r1", "user deny", now, 24*time.Hour); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	ctx, err := BuildPolicyContext(d, "", now)
	if err != nil {
		t.Fatalf("build context: %v", err)
	}
	if len(ctx.ActiveRules) != 0 {
		t.Fatalf("expected no active rules after revoke cooldown: %+v", ctx.ActiveRules)
	}
}

func TestPolicyContext_ExcludesSessionSuppressedRules(t *testing.T) {
	d := t.TempDir()
	_, err := overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if err := overlay.SuppressForSession(d, "s1", "r1"); err != nil {
		t.Fatalf("suppress: %v", err)
	}
	ctx, err := BuildPolicyContext(d, "s1", time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("build context: %v", err)
	}
	if len(ctx.ActiveRules) != 0 {
		t.Fatalf("expected no active rules for suppressed session: %+v", ctx.ActiveRules)
	}
}
