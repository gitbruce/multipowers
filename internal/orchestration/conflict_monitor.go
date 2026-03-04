package orchestration

import (
	"path/filepath"
	"sort"
	"strings"
)

// ConflictMonitor checks for structural overlap between changed files and active touched files.
type ConflictMonitor struct{}

// HasOverlap returns whether overlap exists and the sorted overlap list.
func (ConflictMonitor) HasOverlap(changedFiles []string, activeTouched []string) (bool, []string) {
	touched := make(map[string]struct{}, len(activeTouched))
	for _, file := range activeTouched {
		norm := normalizePath(file)
		if norm == "" {
			continue
		}
		touched[norm] = struct{}{}
	}
	overlap := make([]string, 0)
	seen := map[string]struct{}{}
	for _, file := range changedFiles {
		norm := normalizePath(file)
		if norm == "" {
			continue
		}
		if _, ok := touched[norm]; !ok {
			continue
		}
		if _, dup := seen[norm]; dup {
			continue
		}
		seen[norm] = struct{}{}
		overlap = append(overlap, norm)
	}
	sort.Strings(overlap)
	return len(overlap) > 0, overlap
}

func normalizePath(file string) string {
	norm := strings.TrimSpace(file)
	if norm == "" {
		return ""
	}
	norm = strings.ReplaceAll(norm, "\\", "/")
	norm = filepath.ToSlash(filepath.Clean(norm))
	norm = strings.TrimPrefix(norm, "./")
	if norm == "." {
		return ""
	}
	return strings.ToLower(norm)
}
