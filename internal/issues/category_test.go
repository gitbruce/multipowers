package issues

import "testing"

func TestNormalizeCategory_AllowsKnownValues(t *testing.T) {
	cases := map[string]string{
		"logic-error":  CategoryLogicError,
		"integration":  CategoryIntegration,
		"quality-gate": CategoryQualityGate,
		"security":     CategorySecurity,
		"performance":  CategoryPerformance,
		"ux":           CategoryUX,
		"architecture": CategoryArchitecture,
		"logic":        CategoryLogicError,
		"perf":         CategoryPerformance,
		"ui":           CategoryUX,
		"arch":         CategoryArchitecture,
	}
	for in, want := range cases {
		got, err := NormalizeCategory(in)
		if err != nil {
			t.Fatalf("normalize %q: %v", in, err)
		}
		if got != want {
			t.Fatalf("normalize %q = %q want %q", in, got, want)
		}
	}
}

func TestNormalizeCategory_RejectsUnknown(t *testing.T) {
	if _, err := NormalizeCategory("random"); err == nil {
		t.Fatal("expected error for unknown category")
	}
}
