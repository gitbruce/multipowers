package devx

import (
	"encoding/json"
	"fmt"
	"os"
)

type StructureDecision string

const (
	DecisionMustHomomorphic StructureDecision = "MUST_HOMOMORPHIC"
	DecisionAllowFork       StructureDecision = "ALLOW_FORK"
)

type StructureRule struct {
	SourceRoot string            `json:"source_root"`
	TargetRoot string            `json:"target_root"`
	Decision   StructureDecision `json:"decision"`
	Notes      string            `json:"notes"`
}

type StructureRulesConfig struct {
	Rules []StructureRule `json:"rules"`
}

func LoadStructureRules(path string) (StructureRulesConfig, error) {
	var cfg StructureRulesConfig
	body, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(body, &cfg); err != nil {
		return cfg, err
	}
	for _, rule := range cfg.Rules {
		switch rule.Decision {
		case DecisionMustHomomorphic, DecisionAllowFork:
		default:
			return cfg, fmt.Errorf("invalid structure decision %q for %s -> %s", rule.Decision, rule.SourceRoot, rule.TargetRoot)
		}
	}
	return cfg, nil
}
