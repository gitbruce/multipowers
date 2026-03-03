package devx

import (
	"strings"
	"testing"
	"time"
)

func TestWithTempWorktree_UsesIsolatedPath(t *testing.T) {
	r := Runner{
		NowFn: func() time.Time { return time.Unix(1700000000, 0) },
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			return []byte("ok"), nil
		},
	}
	var gotPath string
	err := r.WithTempWorktree("go", "sync-go", func(path string) error {
		gotPath = path
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotPath, ".worktrees/sync-go-") {
		t.Fatalf("unexpected temp path: %s", gotPath)
	}
}
