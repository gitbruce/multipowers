package policy

import (
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

func TestResolver_ExcludesRevokedOrCoolingRules(t *testing.T) {
	rules := []autosync.OverlayRule{{RuleID: "r1"}, {RuleID: "r2"}}
	cool := map[string]autosync.CooldownEntry{
		"r2": {RuleID: "r2", RevokedUntil: time.Date(2026, 3, 7, 0, 0, 0, 0, time.UTC)},
	}
	got := ExcludeRevokedOrCoolingRules(rules, cool, time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC))
	if len(got) != 1 || got[0].RuleID != "r1" {
		t.Fatalf("unexpected filtered rules: %+v", got)
	}
}
