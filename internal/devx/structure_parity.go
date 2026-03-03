package devx

import "sort"

type NameSetDiff struct {
	MissingInTarget []string
	ExtraInTarget   []string
}

func CompareNameSets(source []string, target []string) NameSetDiff {
	sourceSet := make(map[string]struct{}, len(source))
	targetSet := make(map[string]struct{}, len(target))
	for _, s := range source {
		sourceSet[s] = struct{}{}
	}
	for _, t := range target {
		targetSet[t] = struct{}{}
	}

	diff := NameSetDiff{
		MissingInTarget: make([]string, 0),
		ExtraInTarget:   make([]string, 0),
	}

	for s := range sourceSet {
		if _, ok := targetSet[s]; !ok {
			diff.MissingInTarget = append(diff.MissingInTarget, s)
		}
	}
	for t := range targetSet {
		if _, ok := sourceSet[t]; !ok {
			diff.ExtraInTarget = append(diff.ExtraInTarget, t)
		}
	}

	sort.Strings(diff.MissingInTarget)
	sort.Strings(diff.ExtraInTarget)
	return diff
}
