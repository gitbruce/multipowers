package devx

import "testing"

func TestCompareStructure_MustHomomorphicDetectsMissingAndExtra(t *testing.T) {
	got := CompareNameSets(
		[]string{"plan.md", "review.md"},
		[]string{"plan.md", "persona.md"},
	)
	if len(got.MissingInTarget) != 1 || got.MissingInTarget[0] != "review.md" {
		t.Fatalf("unexpected missing set: %#v", got.MissingInTarget)
	}
	if len(got.ExtraInTarget) != 1 || got.ExtraInTarget[0] != "persona.md" {
		t.Fatalf("unexpected extra set: %#v", got.ExtraInTarget)
	}
}
