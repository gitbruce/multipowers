package faq

import (
	"strings"

	"github.com/gitbruce/multipowers/internal/issues"
)

func Classify(msg string) string {
	m := strings.ToLower(strings.TrimSpace(msg))
	category := issues.CategoryLogicError
	switch {
	case m == "":
		category = issues.CategoryLogicError
	case strings.Contains(m, "security") || strings.Contains(m, "xss") || strings.Contains(m, "secret") || strings.Contains(m, "credential") || strings.Contains(m, "authz"):
		category = issues.CategorySecurity
	case strings.Contains(m, "latency") || strings.Contains(m, "slow") || strings.Contains(m, "performance") || strings.Contains(m, "p95") || strings.Contains(m, "benchmark"):
		category = issues.CategoryPerformance
	case strings.Contains(m, "ux") || strings.Contains(m, "ui") || strings.Contains(m, "accessibility") || strings.Contains(m, "layout"):
		category = issues.CategoryUX
	case strings.Contains(m, "integration") || strings.Contains(m, "provider") || strings.Contains(m, "api") || strings.Contains(m, "http") || strings.Contains(m, "timeout"):
		category = issues.CategoryIntegration
	case strings.Contains(m, "architecture") || strings.Contains(m, "module") || strings.Contains(m, "refactor") || strings.Contains(m, "design"):
		category = issues.CategoryArchitecture
	case strings.Contains(m, "boundary") || strings.Contains(m, "policy") || strings.Contains(m, "hook") || strings.Contains(m, "checkpoint") || strings.Contains(m, "blocked"):
		category = issues.CategoryQualityGate
	default:
		category = issues.CategoryLogicError
	}
	norm, err := issues.NormalizeCategory(category)
	if err != nil {
		return issues.CategoryLogicError
	}
	return norm
}
