package validation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandOwnershipSnapshotExists(t *testing.T) {
	root := findProjectRoot()
	snapshot := filepath.Join(root, "docs", "plans", "evidence", "command-boundary", "commands-snapshot-2026-03-05.md")
	if _, err := os.Stat(snapshot); err != nil {
		t.Fatalf("command snapshot missing: %v", err)
	}
}
