package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/autosync/overlay"
)

type PolicyContext struct {
	GeneratedAt time.Time              `json:"generated_at"`
	SessionID   string                 `json:"session_id,omitempty"`
	ActiveRules []autosync.OverlayRule `json:"active_rules"`
}

func BuildPolicyContext(projectDir, sessionID string, now time.Time) (PolicyContext, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	doc, err := overlay.Load(projectDir)
	if err != nil {
		return PolicyContext{}, err
	}
	rules := make([]autosync.OverlayRule, 0, len(doc.Rules))
	for _, rule := range doc.Rules {
		if c, ok := doc.Cooldowns[rule.RuleID]; ok {
			if !c.RevokedUntil.IsZero() && now.Before(c.RevokedUntil) {
				continue
			}
		}
		if overlay.IsSessionSuppressed(projectDir, sessionID, rule.RuleID) {
			continue
		}
		rules = append(rules, rule)
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].RuleID < rules[j].RuleID })
	return PolicyContext{GeneratedAt: now, SessionID: sessionID, ActiveRules: rules}, nil
}

func WriteSnapshot(projectDir string, ctx PolicyContext) (string, error) {
	path := autosync.DefaultPaths(projectDir).SnapshotFile
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
