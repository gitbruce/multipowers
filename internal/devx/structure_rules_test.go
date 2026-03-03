package devx

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadStructureRules_ValidAndInvalid(t *testing.T) {
	_, err := LoadStructureRules(filepath.Join("testdata", "structure-rules-valid.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = LoadStructureRules(filepath.Join("testdata", "structure-rules-invalid.json"))
	if err == nil {
		t.Fatalf("expected invalid rule error")
	}
}

func TestLoadStructureRules_RootTargetsUseClaudeRoot(t *testing.T) {
	cfg, err := LoadStructureRules(filepath.Join("..", "..", "config", "sync", "claude-structure-rules.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, rule := range cfg.Rules {
		if strings.Contains(rule.TargetRoot, ".claude-plugin/") {
			t.Fatalf("unexpected legacy target root: %s", rule.TargetRoot)
		}
	}
}
