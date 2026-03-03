# Go Physical `.claude` Path Migration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** One-shot migrate `go` branch workspace assets from `.claude-plugin/.claude` to root `.claude` so physical paths match `main`.

**Architecture:** Move the workspace tree physically to root `.claude`, then update validation rules, scripts, and tests to reference same-path roots. Keep `.claude-plugin` only for packaging/runtime assets (`bin`, plugin manifests, custom config). Enforce no compatibility layer and no backup: old workspace path must be removed entirely.

**Tech Stack:** Go 1.24, Bash, Git CLI, GitHub Actions, Markdown docs.

---

## Global Constraints

- Required sub-skills during execution:
  - `@superpowers:test-driven-development`
  - `@superpowers:using-git-worktrees`
  - `@superpowers:verification-before-completion`
- Never revert unrelated uncommitted files.
- Never run `git reset --hard`.
- No compatibility layer, symlink, or backup for `.claude-plugin/.claude`.

### Task 1: Physically Move Workspace Tree To Root `.claude`

**Files:**
- Move: `.claude-plugin/.claude/**` -> `.claude/**`
- Delete: `.claude-plugin/.claude/`
- Test: path existence checks in shell

**Step 1: Write the failing check**

```bash
test -d .claude-plugin/.claude && test ! -d .claude
```

**Step 2: Run check to verify it fails after migration intent**

Run: `test ! -d .claude-plugin/.claude && test -d .claude`  
Expected: FAIL before move.

**Step 3: Write minimal implementation**

Run:

```bash
mkdir -p .claude
cp -a .claude-plugin/.claude/. .claude/
rm -rf .claude-plugin/.claude
```

**Step 4: Run check to verify it passes**

Run: `test ! -d .claude-plugin/.claude && test -d .claude`  
Expected: PASS

**Step 5: Commit**

```bash
git add .claude .claude-plugin
git commit -m "refactor(sync): migrate workspace path to root .claude"
```

### Task 2: Update Structure Rules To Same-Path Roots

**Files:**
- Modify: `config/sync/claude-structure-rules.json`
- Modify: `internal/devx/testdata/structure-rules-valid.json`
- Modify: `internal/devx/testdata/structure-rules-invalid.json` (if needed for root updates)
- Test: `internal/devx/structure_rules_test.go`

**Step 1: Write the failing test**

Add a case asserting target roots are `.claude/...` for `go` parity config validation.

```go
func TestLoadStructureRules_RootTargetsUseClaudeRoot(t *testing.T) {
	cfg, err := LoadStructureRules("../../config/sync/claude-structure-rules.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range cfg.Rules {
		if strings.Contains(r.TargetRoot, ".claude-plugin/.claude") {
			t.Fatalf("unexpected legacy target root: %s", r.TargetRoot)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestLoadStructureRules_RootTargetsUseClaudeRoot -v`  
Expected: FAIL with legacy target root found.

**Step 3: Write minimal implementation**

Update rules:
- `.claude/commands -> .claude/commands`
- `.claude/skills -> .claude/skills`
- `.claude/references -> .claude/references`
- `.claude/state -> .claude/state`
- keep `ALLOW_FORK_WITH_NAME_PARITY` entries on `.claude/...`

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestLoadStructureRules_RootTargetsUseClaudeRoot -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add config/sync/claude-structure-rules.json internal/devx/testdata/structure-rules-valid.json internal/devx/testdata/structure-rules-invalid.json internal/devx/structure_rules_test.go
git commit -m "feat(sync): switch structure parity rules to root .claude targets"
```

### Task 3: Update Devx Validation Tests For New Roots

**Files:**
- Modify: `internal/devx/structure_validation_test.go`
- Modify: `internal/devx/runner_structure_git_test.go` (if path literals reference old root)

**Step 1: Write the failing test**

Adjust existing test fixtures to expect `go:.claude/...` lookups and responses.

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run TestValidateStructureParity -v`  
Expected: FAIL due to old path literals.

**Step 3: Write minimal implementation**

Replace `.claude-plugin/.claude/*` literals with `.claude/*` in tests.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run TestValidateStructureParity -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/devx/structure_validation_test.go internal/devx/runner_structure_git_test.go
git commit -m "test(sync): update structure validation tests for root .claude"
```

### Task 4: Update Scripts And Validators Referencing Old Workspace Path

**Files:**
- Modify: `scripts/build-openclaw.sh`
- Modify: `scripts/fix-command-frontmatter.sh`
- Modify: `internal/validation/no_shell_runtime.go`
- Modify: `internal/validation/no_shell_runtime_test.go`
- Modify: any script under `scripts/` that reads `.claude-plugin/.claude`

**Step 1: Write the failing test/check**

```bash
rg -n "\\.claude-plugin/\\.claude" scripts internal/validation -S
```

**Step 2: Run check to verify it fails**

Run above command.  
Expected: finds legacy path references.

**Step 3: Write minimal implementation**

Update references to root `.claude` where they refer to workspace docs.
Do not change `.claude-plugin/bin`, `plugin.json`, or `marketplace.json` paths.

**Step 4: Run check to verify it passes**

Run: `rg -n "\\.claude-plugin/\\.claude" scripts internal/validation -S`  
Expected: no matches (or only explicitly accepted historical text files if any).

**Step 5: Commit**

```bash
git add scripts/build-openclaw.sh scripts/fix-command-frontmatter.sh internal/validation/no_shell_runtime.go internal/validation/no_shell_runtime_test.go
git commit -m "refactor(paths): switch workspace script references to root .claude"
```

### Task 5: Update CLI/Devx Wiring And Wrappers

**Files:**
- Modify: `scripts/validate-claude-structure.sh`
- Modify: `scripts/report-claude-content-diff.sh`
- Modify: `cmd/mp-devx/main_test.go` (if output or behavior assertions need adjustment)

**Step 1: Write the failing check**

```bash
./scripts/validate-claude-structure.sh -dry-run
```

**Step 2: Run check to verify it fails**

Expected: FAIL if rules and wrappers still assume old mapping.

**Step 3: Write minimal implementation**

- Keep entrypoint names unchanged.
- Ensure validation reads updated same-path rules.
- Update content diff script to compare `.claude` scope directly.

**Step 4: Run check to verify it passes**

Run:
- `./scripts/validate-claude-structure.sh -dry-run`
- `./scripts/report-claude-content-diff.sh`

Expected: both exit `0`.

**Step 5: Commit**

```bash
git add scripts/validate-claude-structure.sh scripts/report-claude-content-diff.sh cmd/mp-devx/main_test.go
git commit -m "chore(sync): align devx wrappers with root .claude layout"
```

### Task 6: Update Architecture And Sync Docs To Same-Path Narrative

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/other-differences.md`
- Modify: `docs/architecture/gap-remediation-tracker.md`
- Modify: `custom/docs/sync/go-upstream-diff-discipline.md`
- Modify: `custom/docs/sync/upstream-sync-playbook.md`

**Step 1: Write the failing doc check**

```bash
rg -n "\\.claude/.*-> \\.claude-plugin/\\.claude|\\.claude-plugin/\\.claude" docs/architecture custom/docs/sync -S
```

**Step 2: Run check to verify it fails**

Run above command.  
Expected: finds stale mapping text.

**Step 3: Write minimal implementation**

- Rewrite policy text to same-path `.claude -> .claude`
- Keep exclusion policy (`*.go`, `custom/**`, approved non-parity scopes)
- Remove old path-mapping language

**Step 4: Run check to verify it passes**

Run: `rg -n "\\.claude/.*-> \\.claude-plugin/\\.claude|\\.claude-plugin/\\.claude" docs/architecture custom/docs/sync -S`  
Expected: no matches (or explicitly allowed historical references only, if any).

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/other-differences.md docs/architecture/gap-remediation-tracker.md custom/docs/sync/go-upstream-diff-discipline.md custom/docs/sync/upstream-sync-playbook.md
git commit -m "docs(sync): switch architecture docs to physical root .claude model"
```

### Task 7: End-To-End Verification Evidence

**Files:**
- Create: `docs/plans/evidence/no-shell-runtime/sync/2026-03-03-go-physical-claude-path-migration-verification.md`

**Step 1: Run verification commands**

```bash
go test ./internal/devx ./cmd/mp-devx ./internal/validation -v
./scripts/validate-claude-structure.sh -dry-run
./scripts/mp-devx -action validate-structure-parity -dry-run
./scripts/report-claude-content-diff.sh
rg -n "\\.claude-plugin/\\.claude" scripts internal cmd config .github docs custom/docs -S
test ! -d .claude-plugin/.claude && test -d .claude
```

**Step 2: Confirm result**

Expected:
- all required commands exit `0`
- no active legacy workspace path references remain
- physical directory check passes

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
git add docs/plans/evidence/no-shell-runtime/sync/2026-03-03-go-physical-claude-path-migration-verification.md
git commit -m "test(sync): record physical .claude path migration verification evidence"
```

## Done Criteria

- `.claude-plugin/.claude` is removed from repository tree.
- Root `.claude` is the only workspace path for commands/skills/references/state.
- Structure parity validation runs on same-path `.claude` roots.
- Workspace-consuming scripts reference `.claude/*` (not `.claude-plugin/.claude/*`).
- Architecture/sync docs describe the new physical path model.
- Verification evidence confirms tests, validation, and path-cleanliness.
