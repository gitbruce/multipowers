package isolation

import (
	"os"
	"path/filepath"
	"testing"
)

type gitCall struct {
	dir  string
	name string
	args []string
}

func TestIsolationRuntime_CreateModelSandbox(t *testing.T) {
	projectDir := t.TempDir()
	calls := make([]gitCall, 0, 1)

	runtime := NewRuntimeManager(projectDir, RuntimeConfig{
		BranchPrefix: "bench",
		WorktreeRoot: ".worktrees/bench",
		LogsSubdir:   "logs",
	}, func(dir string, name string, args ...string) ([]byte, error) {
		copied := append([]string{}, args...)
		calls = append(calls, gitCall{dir: dir, name: name, args: copied})
		return []byte("ok"), nil
	})

	sandbox, err := runtime.CreateModelSandbox("run-123", "gpt-4o", "main")
	if err != nil {
		t.Fatalf("CreateModelSandbox error = %v", err)
	}

	wantBranch := "bench/run-123/gpt-4o"
	if sandbox.Branch != wantBranch {
		t.Fatalf("Branch = %q, want %q", sandbox.Branch, wantBranch)
	}

	wantWorktree := filepath.Join(projectDir, ".worktrees/bench", "run-123", "gpt-4o")
	if sandbox.WorktreePath != wantWorktree {
		t.Fatalf("WorktreePath = %q, want %q", sandbox.WorktreePath, wantWorktree)
	}

	wantLogs := filepath.Join(wantWorktree, "logs")
	if sandbox.LogsPath != wantLogs {
		t.Fatalf("LogsPath = %q, want %q", sandbox.LogsPath, wantLogs)
	}
	if _, err := os.Stat(wantLogs); err != nil {
		t.Fatalf("logs path missing: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("git call count = %d, want 1", len(calls))
	}
	if calls[0].dir != projectDir {
		t.Fatalf("git dir = %q, want %q", calls[0].dir, projectDir)
	}
	if calls[0].name != "git" {
		t.Fatalf("command name = %q, want git", calls[0].name)
	}
	wantArgs := []string{"worktree", "add", "-b", wantBranch, wantWorktree, "main"}
	assertStringSliceEqual(t, calls[0].args, wantArgs)
}

func TestIsolationRuntime_CleanupModelSandbox(t *testing.T) {
	projectDir := t.TempDir()
	calls := make([]gitCall, 0, 2)

	runtime := NewRuntimeManager(projectDir, RuntimeConfig{}, func(dir string, name string, args ...string) ([]byte, error) {
		copied := append([]string{}, args...)
		calls = append(calls, gitCall{dir: dir, name: name, args: copied})
		return []byte("ok"), nil
	})

	sandbox := ModelSandbox{
		Model:        "gpt-4o",
		Branch:       "bench/run-123/gpt-4o",
		WorktreePath: filepath.Join(projectDir, ".worktrees/bench", "run-123", "gpt-4o"),
	}
	if err := runtime.CleanupModelSandbox(sandbox); err != nil {
		t.Fatalf("CleanupModelSandbox error = %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("git call count = %d, want 2", len(calls))
	}
	assertStringSliceEqual(t, calls[0].args, []string{"worktree", "remove", "--force", sandbox.WorktreePath})
	assertStringSliceEqual(t, calls[1].args, []string{"branch", "-D", sandbox.Branch})
}

func assertStringSliceEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d; got=%v want=%v", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("index %d = %q, want %q; got=%v want=%v", i, got[i], want[i], got, want)
		}
	}
}
