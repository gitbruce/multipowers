# Main-Go Naming Parity Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enforce naming parity with `main` for Markdown and Claude Code related files while excluding `./custom/**` and keeping all Go files unconstrained in `go`.

**Architecture:** Extend the structure-rules contract with a new decision type for controlled forks that still require naming parity and explicit add/delete registration. Apply scope filters during structure validation so only `*.md` and configured Claude-related files are evaluated, with a hard exclusion for `./custom/**`. Keep existing `MUST_HOMOMORPHIC` behavior intact and add CI lanes for required parity checks plus non-blocking content-diff reporting.

**Tech Stack:** Go 1.24, Git CLI (`ls-tree`), `mp-devx`, GitHub Actions, Markdown docs.

---

## Global Constraints

- Required sub-skills during execution:
  - `@superpowers:test-driven-development`
  - `@superpowers:using-git-worktrees`
  - `@superpowers:verification-before-completion`
- Never revert local uncommitted files.
- Never run `git reset --hard` or branch-switch the active developer workspace.
- All sync-related branch mutations run in isolated `.worktrees/*`.

### Task 1: Extend Rules Contract For Naming-Parity Forks

**Files:**
- Modify: `internal/devx/structure_rules.go`
- Modify: `internal/devx/structure_rules_test.go`
- Modify: `internal/devx/testdata/structure-rules-valid.json`
- Modify: `internal/devx/testdata/structure-rules-invalid.json`
- Modify: `config/sync/claude-structure-rules.json`

**Step 1: Write the failing test**

```go
func TestLoadStructureRules_AllowsNameParityForkDecision(t *testing.T) {
	_, err := LoadStructureRules(filepath.Join("testdata", "structure-rules-name-parity-valid.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

Add invalid case:

```go
func TestLoadStructureRules_RejectsNameParityForkWithoutRegistry(t *testing.T) {
	_, err := LoadStructureRules(filepath.Join("testdata", "structure-rules-name-parity-invalid.json"))
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run "TestLoadStructureRules_AllowsNameParityForkDecision|TestLoadStructureRules_RejectsNameParityForkWithoutRegistry" -v`  
Expected: FAIL with unknown decision / missing fields validation.

**Step 3: Write minimal implementation**

Add new decision and fields:

```go
const (
	DecisionMustHomomorphic      StructureDecision = "MUST_HOMOMORPHIC"
	DecisionAllowFork            StructureDecision = "ALLOW_FORK"
	DecisionAllowForkNameParity  StructureDecision = "ALLOW_FORK_WITH_NAME_PARITY"
)

type StructureRule struct {
	SourceRoot            string            `json:"source_root"`
	TargetRoot            string            `json:"target_root"`
	Decision              StructureDecision `json:"decision"`
	EnforcePatterns       []string          `json:"enforce_patterns"`
	AllowedTargetAdds     []string          `json:"allowed_target_additions"`
	AllowedTargetRemovals []string          `json:"allowed_target_removals"`
	IgnoreSourceNames     []string          `json:"ignore_source_names"`
	IgnoreTargetNames     []string          `json:"ignore_target_names"`
	Notes                 string            `json:"notes"`
}
```

Validation rule:
- `ALLOW_FORK_WITH_NAME_PARITY` requires non-empty `enforce_patterns`.

Update baseline config:
- move `init.md`, `mp.md`, `persona.md`, `skill-persona.md` into explicit name-parity fork rules
- include `enforce_patterns: ["*.md"]`
- register allowed additions/removals explicitly.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run "TestLoadStructureRules_AllowsNameParityForkDecision|TestLoadStructureRules_RejectsNameParityForkWithoutRegistry" -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/structure_rules.go internal/devx/structure_rules_test.go internal/devx/testdata/structure-rules-valid.json internal/devx/testdata/structure-rules-invalid.json config/sync/claude-structure-rules.json
git commit -m "feat(sync): add allow-fork-with-name-parity rules contract"
```

### Task 2: Add Scope Filter For Markdown + Claude-Related Files

**Files:**
- Modify: `internal/devx/structure_validation.go`
- Create: `internal/devx/structure_scope.go`
- Create: `internal/devx/structure_scope_test.go`

**Step 1: Write the failing test**

```go
func TestScopeFilter_ExcludesGoFilesAndKeepsMarkdown(t *testing.T) {
	in := []string{"plan.md", "impl.go", "skill.md"}
	got := FilterByPatterns(in, []string{"*.md"})
	if len(got) != 2 {
		t.Fatalf("expected only markdown files, got %#v", got)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestScopeFilter_ExcludesGoFilesAndKeepsMarkdown -v`  
Expected: FAIL with `undefined: FilterByPatterns`

**Step 3: Write minimal implementation**

```go
func FilterByPatterns(names []string, patterns []string) []string
```

Use `path.Match` with stable sorted output and default passthrough when no patterns provided.
Add a fast-path exclusion helper for `custom/` prefixes.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestScopeFilter_ExcludesGoFilesAndKeepsMarkdown -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/structure_scope.go internal/devx/structure_scope_test.go internal/devx/structure_validation.go
git commit -m "feat(sync): add scope filter for markdown naming parity checks"
```

### Task 3: Implement Name-Parity Fork Validation + Explicit Registry Check

**Files:**
- Modify: `internal/devx/structure_validation.go`
- Modify: `internal/devx/structure_validation_test.go`

**Step 1: Write the failing test**

```go
func TestValidateStructureParity_NameParityForkRequiresRegisteredAddsRemovals(t *testing.T) {
	// source has persona.md; target has persona.md + init.md
	// rule permits ALLOW_FORK_WITH_NAME_PARITY with *.md but no add registry
	// expect validation error mentioning init.md unregistered addition
}
```

Add pass case with explicit registry:

```go
func TestValidateStructureParity_NameParityForkPassesWithRegistry(t *testing.T) {
	// same inputs but rule includes allowed_target_additions=["init.md"]
	// expect no error
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run "TestValidateStructureParity_NameParityForkRequiresRegisteredAddsRemovals|TestValidateStructureParity_NameParityForkPassesWithRegistry" -v`  
Expected: FAIL because decision path is not implemented.

**Step 3: Write minimal implementation**

- In `ValidateStructureParity`, handle `ALLOW_FORK_WITH_NAME_PARITY`:
  - list source/target names
  - apply ignore filters
  - apply pattern filters (e.g., `*.md`)
  - compare name sets
  - allow only differences present in explicit addition/removal registries
  - fail with per-file violation details for unregistered diffs

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run "TestValidateStructureParity_NameParityForkRequiresRegisteredAddsRemovals|TestValidateStructureParity_NameParityForkPassesWithRegistry" -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/structure_validation.go internal/devx/structure_validation_test.go
git commit -m "feat(sync): validate allow-fork name parity with explicit diff registry"
```

### Task 4: Keep `mp-devx` Structure Validation Backward Compatible

**Files:**
- Modify: `cmd/mp-devx/main_test.go`
- Modify: `cmd/mp-devx/main.go` (if needed for output clarity only)

**Step 1: Write the failing test**

```go
func TestRun_ActionValidateStructureParity_SupportsNameParityForkRules(t *testing.T) {
	rc := run([]string{"-action", "validate-structure-parity", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd/mp-devx -run TestRun_ActionValidateStructureParity_SupportsNameParityForkRules -v`  
Expected: FAIL if parsing/validation breaks with new schema.

**Step 3: Write minimal implementation**

Keep action wiring unchanged; only adjust load/validation path if schema or diagnostics changed.

**Step 4: Run test to verify it passes**

Run: `go test ./cmd/mp-devx -run TestRun_ActionValidateStructureParity_SupportsNameParityForkRules -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add cmd/mp-devx/main.go cmd/mp-devx/main_test.go
git commit -m "test(sync): ensure mp-devx accepts name parity fork rules"
```

### Task 5: Split CI Into Required Parity Gate + Informational Content Diff

**Files:**
- Modify: `.github/workflows/claude-structure-parity.yml`
- Create: `scripts/report-claude-content-diff.sh`

**Step 1: Write the failing check**

```bash
test -f scripts/report-claude-content-diff.sh
```

**Step 2: Run check to verify it fails**

Run: `test -f scripts/report-claude-content-diff.sh`  
Expected: FAIL

**Step 3: Write minimal implementation**

- Keep existing required job:
  - `./scripts/validate-claude-structure.sh -dry-run`
- Add informational job:
  - run non-blocking content-level diff script
  - upload artifact or print summary without failing the workflow

Script outline:

```bash
#!/usr/bin/env bash
set -euo pipefail
git diff --name-status main...go -- .claude .claude-plugin/.claude || true
```

**Step 4: Run checks to verify they pass**

Run:
- `test -f scripts/report-claude-content-diff.sh`
- `rg -n "validate-claude-structure\\.sh -dry-run|report-claude-content-diff\\.sh" .github/workflows/claude-structure-parity.yml`

Expected: PASS

**Step 5: Commit**

```bash
git add .github/workflows/claude-structure-parity.yml scripts/report-claude-content-diff.sh
git commit -m "ci(sync): add informational content diff reporting lane"
```

### Task 6: Update Policy Docs With Final Naming Scope

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`

**Step 1: Write the failing doc checks**

```bash
rg -n "ALLOW_FORK_WITH_NAME_PARITY|\\.go.*excluded|custom/\\*\\*.*excluded|explicitly registered" docs/architecture/commands_skills_difference.md
```

**Step 2: Run check to verify it fails**

Run above command.  
Expected: no matches.

**Step 3: Write minimal implementation**

- Document policy:
  - markdown + Claude-related file naming parity with `main`
  - `./custom/**` is excluded from naming parity enforcement
  - Go files excluded from naming parity
  - add/remove in forked scopes requires explicit registry entries

**Step 4: Run checks to verify they pass**

Run:
- `rg -n "ALLOW_FORK_WITH_NAME_PARITY" docs/architecture/commands_skills_difference.md`
- `rg -n "custom/\\*\\*.*excluded|Go files.*excluded|\\.go files.*excluded" docs/architecture/commands_skills_difference.md`

Expected: PASS

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md
git commit -m "docs(sync): codify markdown and claude-file naming parity policy"
```

### Task 7: End-To-End Verification Evidence

**Files:**
- Create: `docs/plans/evidence/no-shell-runtime/sync/2026-03-03-main-go-naming-parity-verification.md`

**Step 1: Run verification commands**

```bash
./scripts/validate-claude-structure.sh -dry-run
go test ./internal/devx ./cmd/mp-devx -v
./scripts/report-claude-content-diff.sh
```

**Step 2: Confirm result**

Expected:
- parity validate exits `0`
- tests exit `0`
- informational report may show diffs but exits `0`

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
git add docs/plans/evidence/no-shell-runtime/sync/2026-03-03-main-go-naming-parity-verification.md
git commit -m "test(sync): record naming parity policy verification evidence"
```

## Done Criteria

- Rules contract supports controlled forks with naming parity and explicit registry.
- Validation scope enforces naming parity for Markdown and Claude-related files.
- `./custom/**` is excluded from naming parity enforcement.
- Go files are excluded from naming parity checks.
- CI has required parity gate and non-blocking content diff report.
- Sync docs describe the finalized policy and exception workflow.
- Verification evidence captured and timestamped.
