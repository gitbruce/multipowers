package faq

import (
	"testing"

	"github.com/gitbruce/multipowers/internal/issues"
)

func TestClassify_ReturnsNormalizedCategory(t *testing.T) {
	cases := map[string]string{
		"security finding: leaked secret": issues.CategorySecurity,
		"slow benchmark p95 spike":        issues.CategoryPerformance,
		"UI accessibility issue":          issues.CategoryUX,
		"provider api timeout":            issues.CategoryIntegration,
		"architecture module split":       issues.CategoryArchitecture,
		"boundary blocked by policy":      issues.CategoryQualityGate,
		"plain runtime error":             issues.CategoryLogicError,
	}
	for msg, want := range cases {
		if got := Classify(msg); got != want {
			t.Fatalf("classify(%q)=%q want %q", msg, got, want)
		}
	}
}
