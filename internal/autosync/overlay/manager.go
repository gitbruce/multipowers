package overlay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

const (
	ChoiceDelete      = "delete"
	ChoiceSkipSession = "skip-this-session"
)

type Document struct {
	Rules     map[string]autosync.OverlayRule   `json:"rules"`
	Cooldowns map[string]autosync.CooldownEntry `json:"cooldowns,omitempty"`
}

func DenialOptions() []string {
	return []string{ChoiceDelete, ChoiceSkipSession}
}

func Load(projectDir string) (Document, error) {
	path := autosync.DefaultPaths(projectDir).OverlayFile
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Document{Rules: map[string]autosync.OverlayRule{}, Cooldowns: map[string]autosync.CooldownEntry{}}, nil
		}
		return Document{}, err
	}
	var d Document
	if err := json.Unmarshal(b, &d); err != nil {
		return Document{}, err
	}
	if d.Rules == nil {
		d.Rules = map[string]autosync.OverlayRule{}
	}
	if d.Cooldowns == nil {
		d.Cooldowns = map[string]autosync.CooldownEntry{}
	}
	return d, nil
}

func save(projectDir string, doc Document) (string, error) {
	path := autosync.DefaultPaths(projectDir).OverlayFile
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, path); err != nil {
		return "", err
	}
	return path, nil
}

func ApplyProposal(projectDir string, p autosync.Proposal) (string, error) {
	doc, err := Load(projectDir)
	if err != nil {
		return "", err
	}
	ruleID := strings.TrimSpace(p.RuleID)
	if ruleID == "" {
		ruleID = strings.ToLower(strings.TrimSpace(p.Dimension + ":" + p.Value))
	}
	doc.Rules[ruleID] = autosync.OverlayRule{
		RuleID:    ruleID,
		Dimension: strings.TrimSpace(p.Dimension),
		Value:     strings.TrimSpace(p.Value),
		Weight:    p.Confidence,
	}
	return save(projectDir, doc)
}

func RevokeRule(projectDir, ruleID, reason string, now time.Time, cooldown time.Duration) error {
	if strings.TrimSpace(ruleID) == "" {
		return fmt.Errorf("rule id is required")
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if cooldown <= 0 {
		cooldown = 24 * time.Hour
	}
	doc, err := Load(projectDir)
	if err != nil {
		return err
	}
	delete(doc.Rules, ruleID)
	doc.Cooldowns[ruleID] = autosync.CooldownEntry{
		RuleID:       ruleID,
		RevokedAt:    now,
		RevokedUntil: now.Add(cooldown),
		Reason:       strings.TrimSpace(reason),
	}
	if _, err := save(projectDir, doc); err != nil {
		return err
	}
	paths := autosync.DefaultPaths(projectDir)
	_ = purgeRuleFromJSONL(paths.ProposalsFile, ruleID)
	_ = purgeRuleFromJSONL(paths.SamplesFile, ruleID)
	_ = appendAppliedAudit(paths.AppliedFile, map[string]any{
		"action":    "revoked_by_user",
		"rule_id":   ruleID,
		"reason":    reason,
		"timestamp": now.UTC().Format(time.RFC3339),
	})
	return nil
}

func appendAppliedAudit(path string, row map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(row)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

func purgeRuleFromJSONL(path, ruleID string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	s := bufio.NewScanner(strings.NewReader(string(b)))
	out := make([]string, 0)
	needle := `"rule_id":"` + ruleID + `"`
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, needle) {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		out = append(out, line)
	}
	if err := s.Err(); err != nil {
		return err
	}
	if len(out) == 0 {
		return os.WriteFile(path, []byte{}, 0o644)
	}
	return os.WriteFile(path, []byte(strings.Join(out, "\n")+"\n"), 0o644)
}
