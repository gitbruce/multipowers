package devx

import (
	"fmt"
	"sort"
	"strings"
)

func (r Runner) ValidateStructureParity(cfg StructureRulesConfig, sourceRef, targetRef string) error {
	violations := make([]string, 0)
	for _, rule := range cfg.Rules {
		if rule.Decision != DecisionMustHomomorphic {
			continue
		}

		sourceNames, err := r.ListTreeNames(sourceRef, rule.SourceRoot)
		if err != nil {
			return fmt.Errorf("list source names for %s@%s: %w", sourceRef, rule.SourceRoot, err)
		}
		targetNames, err := r.ListTreeNames(targetRef, rule.TargetRoot)
		if err != nil {
			return fmt.Errorf("list target names for %s@%s: %w", targetRef, rule.TargetRoot, err)
		}

		sourceNames = filterIgnoredNames(sourceNames, rule.IgnoreSourceNames)
		targetNames = filterIgnoredNames(targetNames, rule.IgnoreTargetNames)

		diff := CompareNameSets(sourceNames, targetNames)
		if len(diff.MissingInTarget) == 0 && len(diff.ExtraInTarget) == 0 {
			continue
		}

		violations = append(violations, fmt.Sprintf(
			"%s -> %s | missing_in_target=%s | extra_in_target=%s",
			rule.SourceRoot,
			rule.TargetRoot,
			formatNameList(diff.MissingInTarget),
			formatNameList(diff.ExtraInTarget),
		))
	}

	if len(violations) == 0 {
		return nil
	}

	sort.Strings(violations)
	return fmt.Errorf("structure parity violations:\n- %s", strings.Join(violations, "\n- "))
}

func filterIgnoredNames(names []string, ignored []string) []string {
	if len(ignored) == 0 || len(names) == 0 {
		return names
	}

	ignoredSet := make(map[string]struct{}, len(ignored))
	for _, name := range ignored {
		ignoredSet[name] = struct{}{}
	}

	filtered := make([]string, 0, len(names))
	for _, name := range names {
		if _, skip := ignoredSet[name]; skip {
			continue
		}
		filtered = append(filtered, name)
	}
	return filtered
}

func formatNameList(names []string) string {
	if len(names) == 0 {
		return "[]"
	}
	return "[" + strings.Join(names, ", ") + "]"
}
