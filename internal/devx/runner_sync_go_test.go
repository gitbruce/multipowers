package devx

import (
	"strings"
	"testing"
	"time"
)

func TestRunSyncMainToGo_CopiesOnlyCopyFromMainRules(t *testing.T) {
	var calls []string
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			calls = append(calls, name+" "+strings.Join(args, " "))
			return []byte("ok"), nil
		},
		NowFn: func() time.Time { return time.Unix(1700000002, 0) },
	}
	cfg := SyncRulesConfig{
		Rules: []SyncRule{
			{Name: "wrappers", Decision: DecisionCopyFromMain, Paths: []string{"deploy.sh", "install.sh"}},
			{Name: "runtime", Decision: DecisionMigrateToGo, Paths: []string{"scripts/orchestrate.sh"}},
		},
	}
	if err := r.RunSyncMainToGo(cfg, SyncOptions{DryRun: false, Push: false}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	requireContains(t, calls, "git checkout main -- deploy.sh install.sh")
	requireNotContains(t, calls, "scripts/orchestrate.sh")
}

func requireNotContains(t *testing.T, got []string, needle string) {
	t.Helper()
	for _, v := range got {
		if strings.Contains(v, needle) {
			t.Fatalf("did not expect call containing %q, got %#v", needle, got)
		}
	}
}
