package devx

import (
	"encoding/json"
	"fmt"
	"os"
)

type SyncDecision string

const (
	DecisionCopyFromMain   SyncDecision = "COPY_FROM_MAIN"
	DecisionMigrateToGo    SyncDecision = "MIGRATE_TO_GO"
	DecisionKeepInGo       SyncDecision = "KEEP_IN_GO"
	DecisionExclude        SyncDecision = "EXCLUDE_WITH_REASON"
	DecisionDeferCondition SyncDecision = "DEFER_WITH_CONDITION"
)

type SyncRule struct {
	Name     string       `json:"name"`
	Decision SyncDecision `json:"decision"`
	Paths    []string     `json:"paths"`
}

type SyncRulesConfig struct {
	Rules []SyncRule `json:"rules"`
}

func LoadSyncRules(path string) (SyncRulesConfig, error) {
	var cfg SyncRulesConfig
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	for _, r := range cfg.Rules {
		switch r.Decision {
		case DecisionCopyFromMain, DecisionMigrateToGo, DecisionKeepInGo, DecisionExclude, DecisionDeferCondition:
		default:
			return cfg, fmt.Errorf("invalid decision %q in rule %q", r.Decision, r.Name)
		}
	}
	return cfg, nil
}
