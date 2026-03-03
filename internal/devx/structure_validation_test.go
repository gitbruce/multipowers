package devx

import (
	"strings"
	"testing"
)

func TestValidateStructureParity_PassesWithExplicitIgnores(t *testing.T) {
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			if len(args) < 5 || args[0] != "ls-tree" {
				t.Fatalf("unexpected args: %#v", args)
			}
			ref := args[3]
			root := args[4]
			switch ref + ":" + root {
			case "main:.claude/commands":
				return []byte(".claude/commands/plan.md\n.claude/commands/octo.md\n"), nil
			case "go:.claude/commands":
				return []byte(".claude/commands/plan.md\n.claude/commands/mp.md\n"), nil
			default:
				t.Fatalf("unexpected lookup: %s:%s", ref, root)
				return nil, nil
			}
		},
	}
	cfg := StructureRulesConfig{
		Rules: []StructureRule{
			{
				SourceRoot:        ".claude/commands",
				TargetRoot:        ".claude/commands",
				Decision:          DecisionMustHomomorphic,
				IgnoreSourceNames: []string{"octo.md"},
				IgnoreTargetNames: []string{"mp.md"},
			},
		},
	}

	if err := r.ValidateStructureParity(cfg, "main", "go"); err != nil {
		t.Fatalf("unexpected parity error: %v", err)
	}
}

func TestValidateStructureParity_FailsOnUnignoredDifferences(t *testing.T) {
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			if len(args) < 5 || args[0] != "ls-tree" {
				t.Fatalf("unexpected args: %#v", args)
			}
			ref := args[3]
			root := args[4]
			switch ref + ":" + root {
			case "main:.claude/commands":
				return []byte(".claude/commands/plan.md\n.claude/commands/review.md\n"), nil
			case "go:.claude/commands":
				return []byte(".claude/commands/plan.md\n"), nil
			default:
				t.Fatalf("unexpected lookup: %s:%s", ref, root)
				return nil, nil
			}
		},
	}
	cfg := StructureRulesConfig{
		Rules: []StructureRule{
			{
				SourceRoot: ".claude/commands",
				TargetRoot: ".claude/commands",
				Decision:   DecisionMustHomomorphic,
			},
		},
	}

	err := r.ValidateStructureParity(cfg, "main", "go")
	if err == nil {
		t.Fatalf("expected parity validation error")
	}
	if !strings.Contains(err.Error(), "review.md") {
		t.Fatalf("expected missing file in error, got: %v", err)
	}
}
