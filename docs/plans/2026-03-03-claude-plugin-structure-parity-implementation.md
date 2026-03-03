# Claude Plugin Structure Parity Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enforce a deterministic policy for `.claude-plugin/.claude` where required directories stay structurally aligned with `main` while explicitly allowing approved `go`-only divergence.

**Architecture:** Introduce a rules contract (`must-homomorphic` vs `allow-fork`) and validate it in `mp-devx` using git tree inspection of `main` and `go`. Keep enforcement read-only by default (`dry-run`/validate), then integrate into CI and sync docs so maintainers no longer rely on overlay mechanisms.

**Tech Stack:** Go 1.24, Git CLI, Bash wrappers, GitHub Actions, Markdown docs.

---

## Global Constraints

- Required sub-skills during execution:
  - `@superpowers:test-driven-development`
  - `@superpowers:using-git-worktrees`
  - `@superpowers:verification-before-completion`
- Never revert local uncommitted files.
- Never run `git reset --hard` or branch-switch the active developer workspace.
- All branch mutation/sync operations must run in isolated `.worktrees/*`.

## Task Status Tracker (Vibe Coding)

Rule: During vibe coding, immediately update this table after each task is finished. A completed task must be marked with `Status=DONE` and `Done=[x]`.

Allowed `Status` values: `NOT_STARTED`, `IN_PROGRESS`, `DONE`, `BLOCKED`.

| Task | Title | Status | Done |
|---|---|---|---|
| 1 | Add Structure Parity Rules Contract | NOT_STARTED | [ ] |
| 2 | Implement Structural Diff Engine (Names + Paths) | NOT_STARTED | [ ] |
| 3 | Add Git Tree Resolver For `main` vs `go` | NOT_STARTED | [ ] |
| 4 | Implement `validate-structure-parity` Action In `mp-devx` | NOT_STARTED | [ ] |
| 5 | Add CI Gate For Structure Parity | NOT_STARTED | [ ] |
| 6 | Update Sync Docs To Reflect No-Overlay + Parity Rules | NOT_STARTED | [ ] |
| 7 | End-To-End Verification Evidence | NOT_STARTED | [ ] |

### Task 1: Add Structure Parity Rules Contract

**Files:**
- Create: `config/sync/claude-structure-rules.json`
- Create: `internal/devx/structure_rules.go`
- Test: `internal/devx/structure_rules_test.go`

**Step 1: Write the failing test**

```go
func TestLoadStructureRules_ValidAndInvalid(t *testing.T) {
	_, err := LoadStructureRules("internal/devx/testdata/structure-rules-valid.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = LoadStructureRules("internal/devx/testdata/structure-rules-invalid.json")
	if err == nil {
		t.Fatalf("expected invalid rule error")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestLoadStructureRules_ValidAndInvalid -v`  
Expected: FAIL with `undefined: LoadStructureRules`

**Step 3: Write minimal implementation**

```go
type StructureDecision string

const (
	DecisionMustHomomorphic StructureDecision = "MUST_HOMOMORPHIC"
	DecisionAllowFork       StructureDecision = "ALLOW_FORK"
)

type StructureRule struct {
	SourceRoot string            `json:"source_root"`
	TargetRoot string            `json:"target_root"`
	Decision   StructureDecision `json:"decision"`
	Notes      string            `json:"notes"`
}
```

Add baseline rules:
- `main:.claude/commands -> go:.claude-plugin/.claude/commands` = `MUST_HOMOMORPHIC`
- `main:.claude/skills -> go:.claude-plugin/.claude/skills` = `MUST_HOMOMORPHIC`
- `main:.claude/references -> go:.claude-plugin/.claude/references` = `MUST_HOMOMORPHIC`
- `go-only init/mp/persona + skill-persona` in explicit `ALLOW_FORK` section.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestLoadStructureRules_ValidAndInvalid -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add config/sync/claude-structure-rules.json internal/devx/structure_rules.go internal/devx/structure_rules_test.go internal/devx/testdata/structure-rules-valid.json internal/devx/testdata/structure-rules-invalid.json
git commit -m "feat(sync): add claude structure parity rules contract"
```

### Task 2: Implement Structural Diff Engine (Names + Paths)

**Files:**
- Create: `internal/devx/structure_parity.go`
- Test: `internal/devx/structure_parity_test.go`

**Step 1: Write the failing test**

```go
func TestCompareStructure_MustHomomorphicDetectsMissingAndExtra(t *testing.T) {
	got := CompareNameSets(
		[]string{"plan.md", "review.md"},
		[]string{"plan.md", "persona.md"},
	)
	if len(got.MissingInTarget) != 1 || got.MissingInTarget[0] != "review.md" {
		t.Fatalf("unexpected missing set: %#v", got.MissingInTarget)
	}
	if len(got.ExtraInTarget) != 1 || got.ExtraInTarget[0] != "persona.md" {
		t.Fatalf("unexpected extra set: %#v", got.ExtraInTarget)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestCompareStructure_MustHomomorphicDetectsMissingAndExtra -v`  
Expected: FAIL with `undefined: CompareNameSets`

**Step 3: Write minimal implementation**

```go
type NameSetDiff struct {
	MissingInTarget []string
	ExtraInTarget   []string
}

func CompareNameSets(source []string, target []string) NameSetDiff
```

Add stable sort of result slices.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestCompareStructure_MustHomomorphicDetectsMissingAndExtra -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/structure_parity.go internal/devx/structure_parity_test.go
git commit -m "feat(sync): add structure parity diff engine"
```

### Task 3: Add Git Tree Resolver For `main` vs `go`

**Files:**
- Modify: `internal/devx/runner.go`
- Test: `internal/devx/runner_structure_git_test.go`

**Step 1: Write the failing test**

```go
func TestListTreeNames_UsesGitLsTree(t *testing.T) {
	var called bool
	r := Runner{
		RunFn: func(dir, name string, args ...string) ([]byte, error) {
			called = true
			return []byte(".claude/commands/plan.md\n.claude/commands/review.md\n"), nil
		},
	}
	_, err := r.ListTreeNames("main", ".claude/commands")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected git invocation")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestListTreeNames_UsesGitLsTree -v`  
Expected: FAIL with `Runner.ListTreeNames undefined`

**Step 3: Write minimal implementation**

```go
func (r Runner) ListTreeNames(ref string, root string) ([]string, error) {
	out, err := r.run(".", "git", "ls-tree", "-r", "--name-only", ref, root)
	// trim root prefix + return basenames
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestListTreeNames_UsesGitLsTree -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/runner.go internal/devx/runner_structure_git_test.go
git commit -m "feat(sync): add git tree resolver for structure parity checks"
```

### Task 4: Implement `validate-structure-parity` Action In `mp-devx`

**Files:**
- Modify: `cmd/mp-devx/main.go`
- Test: `cmd/mp-devx/main_test.go`
- Create: `scripts/validate-claude-structure.sh`

**Step 1: Write the failing test**

```go
func TestRun_ActionValidateStructureParity(t *testing.T) {
	rc := run([]string{"-action", "validate-structure-parity", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd/mp-devx -run TestRun_ActionValidateStructureParity -v`  
Expected: FAIL with `unknown action`

**Step 3: Write minimal implementation**

```go
case "validate-structure-parity":
	cfg, err := devx.LoadStructureRules(*rulesPath)
	if err != nil { ... }
	err = r.ValidateStructureParity(cfg, "main", "go")
	if err != nil { ... }
	fmt.Fprintln(stdout, "structure parity ok")
```

Wrapper:

```bash
#!/usr/bin/env bash
set -euo pipefail
exec ./scripts/mp-devx -action validate-structure-parity "$@"
```

**Step 4: Run test to verify it passes**

Run: `go test ./cmd/mp-devx -run TestRun_ActionValidateStructureParity -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add cmd/mp-devx/main.go cmd/mp-devx/main_test.go scripts/validate-claude-structure.sh
git commit -m "feat(sync): add validate-structure-parity action and wrapper"
```

### Task 5: Add CI Gate For Structure Parity

**Files:**
- Create: `.github/workflows/claude-structure-parity.yml`

**Step 1: Write the failing check**

```bash
test -f .github/workflows/claude-structure-parity.yml
```

**Step 2: Run check to verify it fails**

Run: `test -f .github/workflows/claude-structure-parity.yml`  
Expected: FAIL

**Step 3: Write minimal implementation**

```yaml
name: claude-structure-parity
on:
  pull_request:
  push:
    branches: [ "go" ]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
      - run: ./scripts/validate-claude-structure.sh -dry-run
```

**Step 4: Run check to verify it passes**

Run: `test -f .github/workflows/claude-structure-parity.yml && rg -n "validate-claude-structure\\.sh -dry-run" .github/workflows/claude-structure-parity.yml`  
Expected: PASS

**Step 5: Commit**

```bash
git add .github/workflows/claude-structure-parity.yml
git commit -m "ci(sync): add claude structure parity validation workflow"
```

### Task 6: Update Sync Docs To Reflect No-Overlay + Parity Rules

**Files:**
- Modify: `custom/docs/sync/upstream-sync-playbook.md`
- Modify: `custom/docs/sync/conflict-resolution.md`
- Modify: `custom/docs/sync/go-upstream-diff-discipline.md`
- Modify: `custom/docs/sync/verification-transcript.md`
- Modify: `docs/architecture/commands_skills_difference.md`

**Step 1: Write the failing doc checks**

```bash
rg -n "overlay|mp-devx overlay|multipowers" custom/docs/sync docs/architecture/commands_skills_difference.md
```

**Step 2: Run check to verify it fails**

Run above command.  
Expected: finds stale entries.

**Step 3: Write minimal implementation**

- Replace sync narrative with `upstream/main -> main -> go`.
- Add parity policy section:
  - `MUST_HOMOMORPHIC`: `commands/skills/references/state`.
  - `ALLOW_FORK`: `init/mp/persona/skill-persona` and explicit go-only runtime artifacts.
- Update transcript examples to `sync-upstream-main/sync-main-to-go/sync-all/validate-structure-parity`.

**Step 4: Run checks to verify they pass**

Run: `rg -n "overlay|mp-devx overlay" custom/docs/sync docs/architecture/commands_skills_difference.md`  
Expected: no matches (except unrelated CSS `z-index overlay` content outside sync docs).

**Step 5: Commit**

```bash
git add custom/docs/sync/upstream-sync-playbook.md custom/docs/sync/conflict-resolution.md custom/docs/sync/go-upstream-diff-discipline.md custom/docs/sync/verification-transcript.md docs/architecture/commands_skills_difference.md
git commit -m "docs(sync): codify must-homomorphic vs allow-fork parity policy"
```

### Task 7: End-To-End Verification Evidence

**Files:**
- Create: `docs/plans/evidence/no-shell-runtime/sync/2026-03-03-structure-parity-verification.md`

**Step 1: Run verification commands**

```bash
./scripts/validate-claude-structure.sh -dry-run
./scripts/sync-all.sh -dry-run
go test ./internal/devx ./cmd/mp-devx -v
scripts/verify-architecture-diff-docs.sh
```

**Step 2: Confirm result**

Expected: all exit `0`.

**Step 3: Record evidence**

For each command record:
- command
- exit code
- first 60 lines output
- UTC timestamp

**Step 4: Freshness re-check**

Run: `./scripts/validate-claude-structure.sh -dry-run`  
Expected: PASS with current HEAD.

**Step 5: Commit**

```bash
git add docs/plans/evidence/no-shell-runtime/sync/2026-03-03-structure-parity-verification.md
git commit -m "test(sync): record structure parity verification evidence"
```

## Done Criteria

- `.claude-plugin/.claude` parity policy is machine-readable and versioned.
- `mp-devx` can validate must-homomorphic subsets against `main`.
- CI runs parity validation on `go` changes.
- Sync docs no longer describe overlay mechanisms.
- Verification evidence exists for dry-run sync + parity checks.
