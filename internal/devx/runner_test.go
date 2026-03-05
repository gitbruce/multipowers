package devx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDevxRunner_SuiteUnitRunsGoTest(t *testing.T) {
	r := Runner{}
	plan, err := r.CommandPlan("unit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan) < 3 || plan[0] != "go" || plan[1] != "test" {
		t.Fatalf("unexpected plan: %#v", plan)
	}
}

func TestDevxRunner_CostReport(t *testing.T) {
	d := t.TempDir()
	file := filepath.Join(d, "model_outputs.2026-03-05.jsonl")
	if err := os.WriteFile(file, []byte("{\"model\":\"gpt-5.3-codex\",\"tokens_input\":10,\"tokens_output\":5}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := Runner{}
	rep, err := r.CostReport(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.TotalInputTokens != 10 || rep.TotalOutputTokens != 5 {
		t.Fatalf("unexpected token totals: in=%d out=%d", rep.TotalInputTokens, rep.TotalOutputTokens)
	}
}

func TestDevxRunner_ValidateRuntimeNoShell(t *testing.T) {
	r := Runner{}
	res, err := r.ValidateRuntimeNoShell(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Valid {
		t.Fatalf("expected valid empty project, got violations: %v", res.Violations)
	}
}
