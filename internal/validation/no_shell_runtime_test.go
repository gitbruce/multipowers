package validation

import "testing"

func TestNoShellRuntimeValidator_FailsOnShellInvocation(t *testing.T) {
	refs := []string{".claude-plugin/commands/persona.md:/scripts/mp persona"}
	got := ValidateNoShellRuntimeRefs(refs)
	if got.Valid {
		t.Fatalf("expected invalid, got valid")
	}
	if len(got.Violations) == 0 {
		t.Fatalf("expected violations to be reported")
	}
}

func TestNoShellRuntimeValidator_PassesWithoutShellInvocation(t *testing.T) {
	refs := []string{".claude-plugin/commands/persona.md:${CLAUDE_PLUGIN_ROOT}/bin/mp persona --json"}
	got := ValidateNoShellRuntimeRefs(refs)
	if !got.Valid {
		t.Fatalf("expected valid, got violations: %v", got.Violations)
	}
}
