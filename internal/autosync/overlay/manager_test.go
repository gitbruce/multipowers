package overlay

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

func TestOverlay_AutoApplyWritesAtomicFile(t *testing.T) {
	d := t.TempDir()
	path, err := ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("overlay file missing: %v", err)
	}
}

func TestOverlay_DenyAsksDeleteOrSkipSession(t *testing.T) {
	opts := DenialOptions()
	if len(opts) != 2 {
		t.Fatalf("options=%v", opts)
	}
	if opts[0] != ChoiceDelete || opts[1] != ChoiceSkipSession {
		t.Fatalf("unexpected options: %v", opts)
	}
}

func TestOverlay_RevokeDeletesLearningDataAndSetsCooldown(t *testing.T) {
	d := t.TempDir()
	_, err := ApplyProposal(d, autosync.Proposal{RuleID: "r2", Dimension: "risk_profile", Value: "low"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	propsFile := autosync.DefaultPaths(d).ProposalsFile
	if err := os.MkdirAll(filepathDir(propsFile), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(propsFile, []byte("{\"rule_id\":\"r2\"}\n{\"rule_id\":\"r3\"}\n"), 0o644); err != nil {
		t.Fatalf("seed proposals: %v", err)
	}
	now := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)
	if err := RevokeRule(d, "r2", "user deny", now, 24*time.Hour); err != nil {
		t.Fatalf("revoke error: %v", err)
	}
	doc, err := Load(d)
	if err != nil {
		t.Fatalf("load overlay: %v", err)
	}
	if _, ok := doc.Rules["r2"]; ok {
		t.Fatal("rule r2 should be removed")
	}
	if _, ok := doc.Cooldowns["r2"]; !ok {
		t.Fatal("expected cooldown for r2")
	}
	b, err := os.ReadFile(propsFile)
	if err != nil {
		t.Fatalf("read proposals: %v", err)
	}
	if strings.Contains(string(b), "\"rule_id\":\"r2\"") {
		t.Fatalf("expected r2 entries removed: %s", string(b))
	}
}

func TestOverlay_SkipSessionSuppressesInjectionOnly(t *testing.T) {
	d := t.TempDir()
	_, err := ApplyProposal(d, autosync.Proposal{RuleID: "r9", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if err := SuppressForSession(d, "s1", "r9"); err != nil {
		t.Fatalf("suppress error: %v", err)
	}
	if !IsSessionSuppressed(d, "s1", "r9") {
		t.Fatal("expected suppressed rule")
	}
	doc, err := Load(d)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, ok := doc.Rules["r9"]; !ok {
		t.Fatal("skip session should not delete persisted rule")
	}
}

func filepathDir(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx <= 0 {
		return "."
	}
	return path[:idx]
}
