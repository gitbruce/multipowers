package devx

import "testing"

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
