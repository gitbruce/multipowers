package issues

import (
	"fmt"
	"sort"
	"strings"
)

const (
	CategoryLogicError   = "logic-error"
	CategoryIntegration  = "integration"
	CategoryQualityGate  = "quality-gate"
	CategorySecurity     = "security"
	CategoryPerformance  = "performance"
	CategoryUX           = "ux"
	CategoryArchitecture = "architecture"
)

var allowedCategories = map[string]string{
	CategoryLogicError:   CategoryLogicError,
	CategoryIntegration:  CategoryIntegration,
	CategoryQualityGate:  CategoryQualityGate,
	CategorySecurity:     CategorySecurity,
	CategoryPerformance:  CategoryPerformance,
	CategoryUX:           CategoryUX,
	CategoryArchitecture: CategoryArchitecture,
	"logic":              CategoryLogicError,
	"quality":            CategoryQualityGate,
	"qualitygate":        CategoryQualityGate,
	"perf":               CategoryPerformance,
	"ui":                 CategoryUX,
	"arch":               CategoryArchitecture,
}

func AllowedCategories() []string {
	out := []string{
		CategoryArchitecture,
		CategoryIntegration,
		CategoryLogicError,
		CategoryPerformance,
		CategoryQualityGate,
		CategorySecurity,
		CategoryUX,
	}
	sort.Strings(out)
	return out
}

func NormalizeCategory(raw string) (string, error) {
	norm := strings.ToLower(strings.TrimSpace(raw))
	norm = strings.ReplaceAll(norm, "_", "-")
	norm = strings.ReplaceAll(norm, " ", "-")
	if v, ok := allowedCategories[norm]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown issue category %q", raw)
}
