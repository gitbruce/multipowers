package devx

import (
	"testing"
)

func TestListTreeNames_UsesGitLsTree(t *testing.T) {
	var called bool
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			called = true
			return []byte(".claude/commands/plan.md\n.claude/commands/review.md\n"), nil
		},
	}
	got, err := r.ListTreeNames("main", ".claude/commands")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected git invocation")
	}
	if len(got) != 2 || got[0] != "plan.md" || got[1] != "review.md" {
		t.Fatalf("unexpected names: %#v", got)
	}
}
