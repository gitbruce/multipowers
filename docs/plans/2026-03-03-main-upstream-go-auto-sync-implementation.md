# Main-Upstream Mirror And Go Auto-Sync Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a fully automated sync pipeline where `main` fast-forwards from `upstream/main`, then shared files are automatically merged into `go` without touching or reverting the developer's current local changes.

**Architecture:** Introduce deterministic sync orchestration in `mp-devx` with explicit actions (`sync-upstream-main`, `sync-main-to-go`, `sync-all`) and a rules contract (`COPY_FROM_MAIN` allowlist). All git mutations run in temporary worktrees so the active workspace is never checked out/reset. Sync behavior is validated via Go tests plus a dry-run CI workflow.

**Tech Stack:** Go 1.24, Bash, Git CLI, Git worktrees, GitHub Actions, Markdown docs.

---

## Execution Guardrails (Must Apply To Every Task)

- Required skills during implementation:
  - `@superpowers:using-git-worktrees`
  - `@superpowers:test-driven-development`
  - `@superpowers:verification-before-completion`
- Never run `git reset --hard`, `git checkout --`, or any revert-style command.
- Never switch branches in the developer's active worktree.
- All branch mutation commands run under temporary worktrees rooted at `.worktrees/sync-*`.
- If the active worktree is dirty, continue safely by isolating operations in temp worktrees; do not modify unstaged files.

## Task Status Tracker (Vibe Coding)

Rule: During vibe coding, immediately update this table after each task is finished. A completed task must be marked with `Status=DONE` and `Done=[x]`.

Allowed `Status` values: `NOT_STARTED`, `IN_PROGRESS`, `DONE`, `BLOCKED`.

| Task | Title | Status | Done |
|---|---|---|---|
| 1 | Define Sync Rules Contract | DONE | [x] |
| 2 | Add Safe Worktree Execution Primitive | DONE | [x] |
| 3 | Implement `sync-upstream-main` (Automated FF Sync) | DONE | [x] |
| 4 | Implement `sync-main-to-go` (Rules-Driven Common File Merge) | DONE | [x] |
| 5 | Expose Full Automation Via `mp-devx` CLI | DONE | [x] |
| 6 | Add CI Dry-Run Automation | DONE | [x] |
| 7 | Update Sync Playbooks To Match `go` Branch Reality | DONE | [x] |
| 8 | End-To-End Verification Before Merge | DONE | [x] |

### Task 1: Define Sync Rules Contract

**Files:**
- Create: `config/sync/main-to-go-rules.json`
- Create: `internal/devx/sync_rules.go`
- Create: `internal/devx/sync_rules_test.go`

**Step 1: Write the failing test**

```go
package devx

import "testing"

func TestLoadSyncRules_ValidAndInvalid(t *testing.T) {
	t.Run("loads valid rules", func(t *testing.T) {
		cfg, err := LoadSyncRules("testdata/rules-valid.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.Rules) == 0 {
			t.Fatalf("expected non-empty rules")
		}
	})

	t.Run("rejects unknown decision", func(t *testing.T) {
		_, err := LoadSyncRules("testdata/rules-invalid-decision.json")
		if err == nil {
			t.Fatalf("expected error for invalid decision")
		}
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestLoadSyncRules_ValidAndInvalid -v`  
Expected: FAIL with `undefined: LoadSyncRules`

**Step 3: Write minimal implementation**

```go
package devx

import (
	"encoding/json"
	"fmt"
	"os"
)

type SyncDecision string

const (
	DecisionCopyFromMain   SyncDecision = "COPY_FROM_MAIN"
	DecisionMigrateToGo    SyncDecision = "MIGRATE_TO_GO"
	DecisionKeepInGo       SyncDecision = "KEEP_IN_GO"
	DecisionExclude        SyncDecision = "EXCLUDE_WITH_REASON"
	DecisionDeferCondition SyncDecision = "DEFER_WITH_CONDITION"
)

type SyncRule struct {
	Name     string       `json:"name"`
	Decision SyncDecision `json:"decision"`
	Paths    []string     `json:"paths"`
}

type SyncRulesConfig struct {
	Rules []SyncRule `json:"rules"`
}

func LoadSyncRules(path string) (SyncRulesConfig, error) {
	var cfg SyncRulesConfig
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	for _, r := range cfg.Rules {
		switch r.Decision {
		case DecisionCopyFromMain, DecisionMigrateToGo, DecisionKeepInGo, DecisionExclude, DecisionDeferCondition:
		default:
			return cfg, fmt.Errorf("invalid decision %q in rule %q", r.Decision, r.Name)
		}
	}
	return cfg, nil
}
```

Also add concrete `COPY_FROM_MAIN` entries in `config/sync/main-to-go-rules.json` for currently shared wrapper assets (for example `deploy.sh`, `install.sh`, `scripts/lib/common.sh`, `tests/run-all.sh`, `tests/run-all-tests.sh`, `.claude/hooks/pre-commit.sh`).

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestLoadSyncRules_ValidAndInvalid -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add config/sync/main-to-go-rules.json internal/devx/sync_rules.go internal/devx/sync_rules_test.go
git commit -m "feat(sync): add main-to-go rules contract"
```

### Task 2: Add Safe Worktree Execution Primitive

**Files:**
- Modify: `internal/devx/runner.go`
- Create: `internal/devx/runner_worktree_test.go`

**Step 1: Write the failing test**

```go
func TestWithTempWorktree_UsesIsolatedPath(t *testing.T) {
	r := Runner{NowFn: func() time.Time { return time.Unix(1700000000, 0) }}
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestWithTempWorktree_UsesIsolatedPath -v`  
Expected: FAIL with `Runner.WithTempWorktree undefined`

**Step 3: Write minimal implementation**

```go
type Runner struct {
	NowFn func() time.Time
	RunFn func(dir string, name string, args ...string) ([]byte, error)
}

func (r Runner) now() time.Time {
	if r.NowFn != nil {
		return r.NowFn()
	}
	return time.Now()
}

func (r Runner) run(dir string, name string, args ...string) ([]byte, error) {
	if r.RunFn != nil {
		return r.RunFn(dir, name, args...)
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

func (r Runner) WithTempWorktree(branch, prefix string, fn func(path string) error) error {
	tmpPath := filepath.Join(".worktrees", fmt.Sprintf("%s-%d", prefix, r.now().Unix()))
	if _, err := r.run(".", "git", "worktree", "add", "--detach", tmpPath, branch); err != nil {
		return err
	}
	defer func() { _, _ = r.run(".", "git", "worktree", "remove", "--force", tmpPath) }()
	return fn(tmpPath)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestWithTempWorktree_UsesIsolatedPath -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/runner.go internal/devx/runner_worktree_test.go
git commit -m "feat(sync): add isolated worktree execution helper"
```

### Task 3: Implement `sync-upstream-main` (Automated FF Sync)

**Files:**
- Modify: `internal/devx/runner.go`
- Create: `internal/devx/runner_sync_upstream_test.go`

**Step 1: Write the failing test**

```go
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
	requireContains(t, calls, "git fetch upstream --prune origin --prune")
	requireContains(t, calls, "git merge --ff-only upstream/main")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestRunSyncUpstreamMain_PlansFetchAndFastForward -v`  
Expected: FAIL with `RunSyncUpstreamMain undefined`

**Step 3: Write minimal implementation**

```go
type SyncOptions struct {
	DryRun bool
	Push   bool
}

func (r Runner) RunSyncUpstreamMain(opts SyncOptions) error {
	return r.WithTempWorktree("main", "sync-main", func(path string) error {
		if _, err := r.run(path, "git", "fetch", "upstream", "--prune", "origin", "--prune"); err != nil {
			return err
		}
		if _, err := r.run(path, "git", "merge", "--ff-only", "upstream/main"); err != nil {
			return err
		}
		if opts.Push && !opts.DryRun {
			if _, err := r.run(path, "git", "push", "origin", "HEAD:main"); err != nil {
				return err
			}
		}
		return nil
	})
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestRunSyncUpstreamMain_PlansFetchAndFastForward -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/runner.go internal/devx/runner_sync_upstream_test.go
git commit -m "feat(sync): add automated upstream-to-main fast-forward action"
```

### Task 4: Implement `sync-main-to-go` (Rules-Driven Common File Merge)

**Files:**
- Modify: `internal/devx/runner.go`
- Create: `internal/devx/runner_sync_go_test.go`
- Modify: `config/sync/main-to-go-rules.json`

**Step 1: Write the failing test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestRunSyncMainToGo_CopiesOnlyCopyFromMainRules -v`  
Expected: FAIL with `RunSyncMainToGo undefined`

**Step 3: Write minimal implementation**

```go
func (r Runner) RunSyncMainToGo(cfg SyncRulesConfig, opts SyncOptions) error {
	return r.WithTempWorktree("go", "sync-go", func(path string) error {
		copyPaths := make([]string, 0)
		for _, rule := range cfg.Rules {
			if rule.Decision == DecisionCopyFromMain {
				copyPaths = append(copyPaths, rule.Paths...)
			}
		}
		if len(copyPaths) == 0 {
			return nil
		}
		args := append([]string{"checkout", "main", "--"}, copyPaths...)
		if _, err := r.run(path, "git", args...); err != nil {
			return err
		}
		if !opts.DryRun {
			if _, err := r.run(path, "git", "add", "--all"); err != nil {
				return err
			}
			if _, err := r.run(path, "git", "commit", "-m", "chore(sync): copy common files from main into go"); err != nil {
				return err
			}
			if opts.Push {
				if _, err := r.run(path, "git", "push", "origin", "HEAD:go"); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestRunSyncMainToGo_CopiesOnlyCopyFromMainRules -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/runner.go internal/devx/runner_sync_go_test.go config/sync/main-to-go-rules.json
git commit -m "feat(sync): add rules-driven main-to-go common file sync"
```

### Task 5: Expose Full Automation Via `mp-devx` CLI

**Files:**
- Modify: `cmd/mp-devx/main.go`
- Create: `cmd/mp-devx/main_test.go`
- Create: `scripts/sync-upstream-main.sh`
- Create: `scripts/sync-main-to-go.sh`
- Create: `scripts/sync-all.sh`

**Step 1: Write the failing test**

```go
func TestRun_ActionSyncAll(t *testing.T) {
	rc := run([]string{"-action", "sync-all", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd/mp-devx -run TestRun_ActionSyncAll -v`  
Expected: FAIL with `unknown action: sync-all`

**Step 3: Write minimal implementation**

```go
// in cmd/mp-devx/main.go
dryRun := flag.Bool("dry-run", false, "plan only, no push/commit")
push := flag.Bool("push", false, "push branch after successful sync")
rules := flag.String("rules", "config/sync/main-to-go-rules.json", "sync rules path")

case "sync-upstream-main":
	err := r.RunSyncUpstreamMain(devx.SyncOptions{DryRun: *dryRun, Push: *push})
case "sync-main-to-go":
	cfg, err := devx.LoadSyncRules(*rules)
	if err == nil {
		err = r.RunSyncMainToGo(cfg, devx.SyncOptions{DryRun: *dryRun, Push: *push})
	}
case "sync-all":
	cfg, err := devx.LoadSyncRules(*rules)
	if err == nil {
		err = r.RunSyncUpstreamMain(devx.SyncOptions{DryRun: *dryRun, Push: *push})
	}
	if err == nil {
		err = r.RunSyncMainToGo(cfg, devx.SyncOptions{DryRun: *dryRun, Push: *push})
	}
```

Wrapper scripts:

```bash
#!/usr/bin/env bash
set -euo pipefail
exec ./scripts/mp-devx -action sync-all "$@"
```

**Step 4: Run test to verify it passes**

Run: `go test ./cmd/mp-devx -run TestRun_ActionSyncAll -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add cmd/mp-devx/main.go cmd/mp-devx/main_test.go scripts/sync-upstream-main.sh scripts/sync-main-to-go.sh scripts/sync-all.sh
git commit -m "feat(sync): expose upstream/main/go sync actions via mp-devx and scripts"
```

### Task 6: Add CI Dry-Run Automation

**Files:**
- Create: `.github/workflows/upstream-sync-automation.yml`

**Step 1: Write the failing test/check**

```bash
test -f .github/workflows/upstream-sync-automation.yml
```

**Step 2: Run check to verify it fails**

Run: `test -f .github/workflows/upstream-sync-automation.yml`  
Expected: FAIL (exit code 1)

**Step 3: Write minimal implementation**

```yaml
name: upstream-sync-automation

on:
  workflow_dispatch:
  schedule:
    - cron: "0 3 * * *"

jobs:
  dry-run-sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
      - run: ./scripts/sync-all.sh -dry-run
```

**Step 4: Run check to verify it passes**

Run: `test -f .github/workflows/upstream-sync-automation.yml && rg -n "sync-all.sh -dry-run" .github/workflows/upstream-sync-automation.yml`  
Expected: PASS and prints matching line

**Step 5: Commit**

```bash
git add .github/workflows/upstream-sync-automation.yml
git commit -m "ci(sync): add scheduled dry-run upstream/main/go sync workflow"
```

### Task 7: Update Sync Playbooks To Match `go` Branch Reality

**Files:**
- Modify: `custom/docs/sync/upstream-sync-playbook.md`
- Modify: `custom/docs/sync/conflict-resolution.md`
- Modify: `custom/docs/sync/verification-transcript.md`
- Modify: `.multipowers/CLAUDE.md`

**Step 1: Write the failing doc checks**

```bash
rg -n "multipowers" custom/docs/sync/upstream-sync-playbook.md custom/docs/sync/conflict-resolution.md .multipowers/CLAUDE.md
```

**Step 2: Run check to verify it fails**

Run: `rg -n "multipowers" custom/docs/sync/upstream-sync-playbook.md custom/docs/sync/conflict-resolution.md .multipowers/CLAUDE.md`  
Expected: FAIL condition for current plan (matches exist and need migration to `go` flow wording)

**Step 3: Write minimal implementation**

- Replace flow wording with: `upstream/main -> main -> go`.
- Add explicit guard clause:
  - "Never run sync by switching current worktree branch."
  - "Use temporary worktree under `.worktrees/sync-*`."
  - "Never resolve by revert/reset of local uncommitted files."
- Add a concrete execution block:

```bash
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
```

**Step 4: Run checks to verify they pass**

Run: `rg -n "upstream/main -> main -> go|\\.worktrees/sync-|Never run sync by switching current worktree" custom/docs/sync/upstream-sync-playbook.md custom/docs/sync/conflict-resolution.md .multipowers/CLAUDE.md`  
Expected: PASS with all key phrases present

**Step 5: Commit**

```bash
git add custom/docs/sync/upstream-sync-playbook.md custom/docs/sync/conflict-resolution.md custom/docs/sync/verification-transcript.md .multipowers/CLAUDE.md
git commit -m "docs(sync): align playbooks to upstream->main->go automation and no-revert policy"
```

### Task 8: End-To-End Verification Before Merge

**Files:**
- Modify: `docs/plans/evidence/no-shell-runtime/` (add sync verification transcript)

**Step 1: Write verification command set**

```bash
./scripts/sync-all.sh -dry-run
go test ./internal/devx ./cmd/mp-devx -v
go test ./... 
scripts/verify-architecture-diff-docs.sh
```

**Step 2: Run verification and confirm failures are zero**

Run the command set above.  
Expected: all commands exit `0`; dry-run prints planned sync actions without branch switches in active worktree.

**Step 3: Record evidence**

Save transcript to:

```text
docs/plans/evidence/no-shell-runtime/sync/2026-03-03-sync-automation-verification.md
```

Include:
- command
- exit code
- key output lines
- timestamp

**Step 4: Re-run critical command to prove freshness**

Run: `./scripts/sync-all.sh -dry-run`  
Expected: PASS with fresh timestamp and no destructive commands

**Step 5: Commit**

```bash
git add docs/plans/evidence/no-shell-runtime/sync/2026-03-03-sync-automation-verification.md
git commit -m "test(sync): record end-to-end dry-run verification evidence"
```

## Final Acceptance Checklist

- `main` sync action uses `--ff-only` and fails fast on non-fast-forward.
- `go` sync action only copies `COPY_FROM_MAIN` paths from rules file.
- Active worktree is never branch-switched during sync automation.
- No command in automation path performs revert/reset of local user edits.
- `sync-all` dry-run is callable locally and in CI.
- Docs and policy files consistently describe `upstream/main -> main -> go`.
