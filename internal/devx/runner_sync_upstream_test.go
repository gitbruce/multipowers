package devx

import (
	"strings"
	"testing"
	"time"
)

func TestRunSyncUpstreamMain_PlansFetchAndFastForward(t *testing.T) {
	var calls []string
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			calls = append(calls, dir+" :: "+name+" "+strings.Join(args, " "))
			return []byte("ok"), nil
		},
		NowFn: func() time.Time { return time.Unix(1700000001, 0) },
	}
	err := r.RunSyncUpstreamMain(SyncOptions{DryRun: false, Push: false})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	requireContains(t, calls, "git fetch upstream --prune")
	requireContains(t, calls, "git fetch origin --prune")
	requireContains(t, calls, "git merge --ff-only upstream/main")
}

func TestRunSyncUpstreamMain_DryRunSkipsMergeAndPush(t *testing.T) {
	var calls []string
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			calls = append(calls, dir+" :: "+name+" "+strings.Join(args, " "))
			return []byte("ok"), nil
		},
		NowFn: func() time.Time { return time.Unix(1700000001, 0) },
	}
	err := r.RunSyncUpstreamMain(SyncOptions{DryRun: true, Push: true})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	requireContains(t, calls, "git fetch upstream --prune")
	requireContains(t, calls, "git fetch origin --prune")
	requireNotContains(t, calls, "git merge --ff-only upstream/main")
	requireNotContains(t, calls, "git push origin HEAD:main")
}

func requireContains(t *testing.T, got []string, needle string) {
	t.Helper()
	for _, v := range got {
		if strings.Contains(v, needle) {
			return
		}
	}
	t.Fatalf("expected call containing %q, got %#v", needle, got)
}
