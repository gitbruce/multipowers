# Strict No-Shell Runtime Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Remove shell as runtime dependency completely by migrating all invoked `.sh` logic to Go, switching all call sites to `bin/octo`, and deleting all `.sh` files only after parity verification passes.

**Architecture:** Keep `cmd/octo` as the only runtime executable. Move remaining shell-owned behaviors (persona routing, state/context/provider utilities, hook helpers, CI/test harness entrypoints) into focused Go packages and subcommands. Enforce with a strict validator that blocks any runtime/docs/CI references to `.sh` execution.

**Tech Stack:** Go 1.24 (`cmd/octo`, `internal/*`, `pkg/api/*`), GitHub Actions, Markdown command/skill files, Makefile.

---

## Scope Guardrails

- In scope: runtime entrypoints, plugin commands/skills call sites, hooks, CI/test harness, docs references that instruct users to run `.sh`, deletion of shell files.
- Out of scope: adding new end-user features unrelated to shell removal.
- Success definition: user-facing and CI-facing execution paths run only through `bin/octo` (or `go run ./cmd/octo` in build contexts), with zero `.sh` files remaining.

## Baseline Evidence to Capture First

- Current shell count: `rg --files | rg '\\.sh$' | wc -l`
- Current shell callsites: `rg -n "orchestrate\\.sh|\\.sh\\b" .claude .claude-plugin scripts hooks Makefile .github docs`
- Save outputs to: `docs/plans/evidence/no-shell-runtime/baseline/`

---

## Target Go Directory Structure (must be established before migration)

```text
cmd/
  octo/                  # end-user runtime binary
  octo-devx/             # maintainer/test/ci helper binary

internal/
  app/                   # pipeline orchestration and error codes
  cli/                   # command parsing and wiring
  context/               # .multipowers context guard/init/summary
  workflows/             # discover/define/develop/deliver/debate/embrace/persona
  providers/             # provider registry/proxy/routing/degrade/quorum
  hooks/                 # SessionStart/PreToolUse/PostToolUse/Stop/SubagentStop
  runtime/               # runtime pre-run and workspace/runtime settings
  tracks/                # track_id, plan/intent persistence, checkbox status
  faq/                   # FAQ extraction/dedup/classification
  validation/            # gates + strict no-shell validator
  execx/                 # safe external command execution abstraction
  render/                # banners/tables/human-readable rendering
  fsboundary/            # target-project filesystem boundary checks
  devx/                  # CI/test harness logic used by octo-devx
  util/                  # json/path/time helpers

pkg/
  api/                   # shared JSON contracts/schema
```

### Mapping Rules (current -> target)

- `scripts/orchestrate*.sh` -> `cmd/octo` + `internal/cli` + `internal/workflows`
- `scripts/state-manager.sh` -> `internal/tracks`
- `scripts/context-manager.sh` -> `internal/context`
- `scripts/provider-router.sh` -> `internal/providers`
- `hooks/*.sh` runtime logic -> `internal/hooks`
- `tests/*.sh` harness -> `cmd/octo-devx` + `internal/devx`
- shell-based validation scripts -> `internal/validation` + `cmd/octo validate`

---

- [x] ### Task 0: Normalize Go Package Layout Before Feature Migration

**Why:** Without a stable package map, migration work will scatter and create second-round refactors.

**Files:**
- Modify: `cmd/octo/main.go`
- Modify: `internal/cli/root.go`
- Create: `internal/devx/.keep` (replace when Task 5 starts)
- Create: `internal/workflows/persona.go` (stub only in this task)
- Modify: `pkg/api/types.go` (if package contracts need relocation support)
- Test: `go test ./...`

**Step 1: Add a failing architecture guard test**

Create a test that asserts required package directories exist and are non-empty for core domains (`context`, `workflows`, `providers`, `hooks`, `tracks`, `validation`, `execx`).

**Step 2: Run test to verify it fails**

Run: `go test ./... -run ArchitectureLayout -v`
Expected: FAIL if required package path is missing or misplaced.

**Step 3: Move files to target locations only when needed**

- Keep package boundaries strict (no cross-domain utility dumping).
- Keep each file < 500 lines; split aggressively when near threshold.
- Ensure `cmd/octo` remains the only user runtime entrypoint.

**Step 4: Run full compile and tests**

Run:
1. `go build ./cmd/octo`
2. `go test ./...`
Expected: PASS.

**Step 5: Commit**

```bash
git add cmd/octo/main.go internal pkg
git commit -m "refactor(arch): normalize go package layout for strict no-shell migration"
```

---

- [x] ### Task 1: Add Strict Runtime No-Shell Validator (red -> green)

**Files:**
- Create: `internal/validation/no_shell_runtime.go`
- Create: `internal/validation/no_shell_runtime_test.go`
- Modify: `internal/cli/root.go`
- Test: `internal/validation/no_shell_runtime_test.go`

**Step 1: Write failing tests for forbidden shell references**

```go
func TestNoShellRuntimeValidator_FailsOnShellInvocation(t *testing.T) {
    refs := []string{".claude-plugin/commands/persona.md:${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona"}
    got := ValidateNoShellRuntimeRefs(refs)
    if got.Valid {
        t.Fatalf("expected invalid, got valid")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run NoShellRuntime -v`
Expected: FAIL (missing validator / wrong behavior).

**Step 3: Implement minimal validator + CLI wiring**

- Add `ValidateNoShellRuntimeRefs` and structured result:

```go
type NoShellRuntimeResult struct {
    Valid        bool     `json:"valid"`
    Violations   []string `json:"violations,omitempty"`
    CheckedFiles int      `json:"checked_files"`
}
```

- Add CLI command: `octo validate --strict-no-shell --json`.

**Step 4: Run tests to verify pass**

Run: `go test ./internal/validation -run NoShellRuntime -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/validation/no_shell_runtime.go internal/validation/no_shell_runtime_test.go internal/cli/root.go
git commit -m "feat(validation): add strict no-shell runtime validator"
```

---

- [x] ### Task 2: Replace Persona Command Runtime Path with Go

**Files:**
- Modify: `.claude-plugin/commands/persona.md`
- Create: `internal/workflows/persona.go`
- Modify: `internal/cli/root.go`
- Test: `internal/workflows/persona_test.go`

**Step 1: Write failing persona workflow tests**

```go
func TestPersonaList_OneLineWithModelAndDescription(t *testing.T) {
    out := RenderPersonaList()
    if !strings.Contains(out, "name | description | model") {
        t.Fatalf("missing table header")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/workflows -run Persona -v`
Expected: FAIL.

**Step 3: Implement persona subcommand in Go and switch command file**

- Add `octo persona --prompt "..." --json` in `root.go`.
- Change `.claude-plugin/commands/persona.md` from `scripts/orchestrate.sh` to `bin/octo persona ...`.

**Step 4: Run tests**

Run: `go test ./internal/workflows -run Persona -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add .claude-plugin/commands/persona.md internal/workflows/persona.go internal/workflows/persona_test.go internal/cli/root.go
git commit -m "feat(persona): route persona command to go runtime"
```

---

- [ ] ### Task 3: Port Remaining Shell Utility Logic to Go (state/context/provider)

**Files:**
- Modify: `internal/context/checker.go`
- Modify: `internal/context/init_runner.go`
- Modify: `internal/providers/router.go`
- Modify: `internal/tracks/state.go`
- Create: `internal/tracks/state_lock_test.go`
- Test: existing `*_test.go` in these packages

**Step 1: Write failing parity tests against legacy behaviors**

- Add table tests for:
1. context required files detection
2. provider routing fallback (2-of-3 quorum)
3. state update atomicity under concurrent writes

**Step 2: Run tests to verify failures**

Run: `go test ./internal/context ./internal/providers ./internal/tracks -run 'Parity|Concurrent|Guard' -v`
Expected: FAIL at least one case.

**Step 3: Implement minimal parity logic in Go**

- Ensure no shell-backed helper is needed for context/state/provider decisions.
- Ensure lock-safe writes for state file operations.

**Step 4: Run tests**

Run: `go test ./internal/context ./internal/providers ./internal/tracks -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/context/checker.go internal/context/init_runner.go internal/providers/router.go internal/tracks/state.go internal/tracks/state_lock_test.go
git commit -m "feat(core): complete go parity for context state and provider routing"
```

---

- [ ] ### Task 4: Move Hook Script Behaviors to Go Hook Subcommands

**Files:**
- Modify: `internal/hooks/handler.go`
- Modify: `internal/hooks/pre_tool_use.go`
- Modify: `internal/hooks/post_tool_use.go`
- Modify: `internal/hooks/stop.go`
- Modify: `hooks/hooks.json`
- Modify: `.claude-plugin/hooks.json`
- Test: `internal/hooks/handler_test.go`

**Step 1: Add failing tests for Stop/SubagentStop expected decisions**

```go
func TestHook_Stop_EmitsTrackAndFaqPostprocess(t *testing.T) {
    // expect no shell execution and deterministic JSON decision
}
```

**Step 2: Run hook tests (expect fail)**

Run: `go test ./internal/hooks -run 'Stop|SubagentStop|PreToolUse|PostToolUse' -v`
Expected: FAIL.

**Step 3: Implement behavior in Go and validate hooks.json schema**

- Ensure all hook commands point to `${CLAUDE_PLUGIN_ROOT}/bin/octo hook ...`.
- Ensure JSON output includes stable keys for Claude hook parser.

**Step 4: Re-run tests**

Run: `go test ./internal/hooks -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/hooks/*.go hooks/hooks.json .claude-plugin/hooks.json
git commit -m "feat(hooks): complete go hook pipeline and remove shell dependencies"
```

---

- [ ] ### Task 5: Replace Shell-Based Test Harness and Makefile Entrypoints

**Files:**
- Modify: `Makefile`
- Create: `cmd/octo-devx/main.go`
- Create: `internal/devx/runner.go`
- Create: `internal/devx/runner_test.go`
- Deprecate references in: `tests/run-all.sh`, `tests/helpers/generate-coverage-report.sh`

**Step 1: Write failing tests for test-suite runner selection**

```go
func TestDevxRunner_SuiteUnitRunsGoTest(t *testing.T) {
    // verify command plan contains: go test ./...
}
```

**Step 2: Run test (expect fail)**

Run: `go test ./internal/devx -v`
Expected: FAIL.

**Step 3: Implement Go test harness command + Makefile switch**

- New command examples:
1. `go run ./cmd/octo-devx --suite smoke`
2. `go run ./cmd/octo-devx --suite unit`
3. `go run ./cmd/octo-devx --suite integration`
- Replace Makefile `.sh` invocations with `go run ./cmd/octo-devx ...`.

**Step 4: Run tests and smoke**

Run:
1. `go test ./internal/devx -v`
2. `make test-unit`
Expected: PASS.

**Step 5: Commit**

```bash
git add Makefile cmd/octo-devx/main.go internal/devx/runner.go internal/devx/runner_test.go
git commit -m "feat(devx): replace shell harness with go runner"
```

---

- [ ] ### Task 6: Remove Shell from CI Workflows

**Files:**
- Modify: `.github/workflows/test.yml`
- Modify: `.github/workflows/claude-octopus.yml`
- Modify: `.github/workflows/go-ci.yml`

**Step 1: Add failing CI lint check in repo to catch shell invocations in workflow YAML**

- Extend validator to parse workflow `run:` blocks for forbidden `.sh` execution.

**Step 2: Run validator (expect fail)**

Run: `go run ./cmd/octo validate --strict-no-shell --dir . --json`
Expected: `status=error`, violations include workflow files.

**Step 3: Replace CI run commands with Go entrypoints**

- Replace `chmod +x scripts/orchestrate.sh` and `./scripts/orchestrate.sh ...` with:
1. `go build -o bin/octo ./cmd/octo`
2. `./bin/octo <cmd> --dir "$PWD" --json`
- Replace test `.sh` harness calls with `go run ./cmd/octo-devx --suite ...`.

**Step 4: Re-run validator**

Run: `go run ./cmd/octo validate --strict-no-shell --dir . --json`
Expected: PASS (no workflow violations).

**Step 5: Commit**

```bash
git add .github/workflows/test.yml .github/workflows/claude-octopus.yml .github/workflows/go-ci.yml
git commit -m "ci: remove shell runtime invocations from workflows"
```

---

- [ ] ### Task 7: Update Command/Skill Markdown Runtime Calls to Go-Only

**Files:**
- Modify: `.claude/commands/*.md` (all spec-driven + persona-related)
- Modify: `.claude/skills/*.md` (runtime invocation examples)
- Modify: `docs/COMMAND-REFERENCE.md`
- Modify: `custom/docs/**/*.md` where `.sh` is presented as runtime path

**Step 1: Add failing doc-reference test**

- Add validator rule: docs under command/skill runtime sections cannot instruct `.sh` execution.

**Step 2: Run validator (expect fail)**

Run: `go run ./cmd/octo validate --strict-no-shell --dir . --json`
Expected: doc violations.

**Step 3: Update docs to Go-only runtime paths**

- Replace `scripts/orchestrate.sh` examples with `bin/octo` equivalents.
- Keep migration history notes in archive docs only (non-runtime section clearly labeled).

**Step 4: Re-run validator**

Run: `go run ./cmd/octo validate --strict-no-shell --dir . --json`
Expected: no docs runtime violations.

**Step 5: Commit**

```bash
git add .claude/commands .claude/skills docs/COMMAND-REFERENCE.md custom/docs
git commit -m "docs: switch runtime guidance from shell to go commands"
```

---

- [ ] ### Task 8: Dual-Run Parity and Performance Gate Before Deletion

**Files:**
- Modify: `scripts/go/dual-run-parity.sh` (replace shell dependency with stored legacy evidence inputs)
- Modify: `scripts/go/benchmark-preflight.sh` (or port to Go benchmark tool)
- Create: `docs/plans/evidence/no-shell-runtime/parity/README.md`
- Create: `docs/plans/evidence/no-shell-runtime/perf/README.md`

**Step 1: Define failing acceptance checks**

- Parity must pass for: `plan`, `develop`, `deliver`, `debate` JSON shape.
- Preflight performance target: p95 < 50ms for `octo context guard --json` (warm cache).

**Step 2: Run checks (expect fail until migrated)**

Run:
1. `./scripts/go/dual-run-parity.sh`
2. `./scripts/go/benchmark-preflight.sh`
Expected: at least one failure before fixes.

**Step 3: Implement required fixes in Go where parity/perf fail**

- Optimize hot paths in context guard and hook preflight.
- Normalize response fields via `pkg/api/types.go`.

**Step 4: Re-run checks and save artifacts**

Expected: parity pass + perf threshold pass with reports committed.

**Step 5: Commit**

```bash
git add scripts/go docs/plans/evidence/no-shell-runtime
# plus any touched go files for parity/perf
git commit -m "test: enforce parity and perf gates before shell deletion"
```

---

- [ ] ### Task 9: Delete All .sh Files in One Batch

**Files:**
- Delete: all `*.sh` tracked files (`rg --files | rg '\\.sh$'`)
- Modify: any references left by deletion fallout

**Step 1: Pre-delete safety check**

Run:
1. `go test ./...`
2. `go run ./cmd/octo validate --strict-no-shell --dir . --json`
Expected: all pass.

**Step 2: Delete files**

Run:
1. `rg --files | rg '\\.sh$' > /tmp/sh-files.txt`
2. `xargs -a /tmp/sh-files.txt git rm`

**Step 3: Build and test**

Run:
1. `go build ./cmd/octo`
2. `go test ./...`
3. `make test-unit`
Expected: PASS without shell runtime.

**Step 4: Re-run strict validator**

Run: `go run ./cmd/octo validate --strict-no-shell --dir . --json`
Expected: zero violations.

**Step 5: Commit**

```bash
git add -A
git commit -m "refactor: remove all shell scripts after go runtime migration completion"
```

---

- [ ] ### Task 10: Final Verification, Version Bump, Release Notes

**Files:**
- Modify: `package.json`
- Modify: `.claude-plugin/plugin.json`
- Modify: `.claude-plugin/marketplace.json`
- Create: `RELEASE_NOTES_NO_SHELL_RUNTIME.md`
- Modify: `custom/README.md`

**Step 1: Run full verification matrix**

Run:
1. `go test ./...`
2. `go vet ./...`
3. `go run ./cmd/octo validate --strict-no-shell --dir . --json`
4. `go run ./cmd/octo init --dir /tmp/octo-target --json`
5. `go run ./cmd/octo develop --dir /tmp/octo-target --prompt "smoke" --json`
Expected: all pass.

**Step 2: Update version + release notes**

- Bump plugin/tool version for no-shell-runtime cut.
- Document migration notes and rollback strategy (git tag only, no shell fallback).

**Step 3: Commit**

```bash
git add package.json .claude-plugin/plugin.json .claude-plugin/marketplace.json RELEASE_NOTES_NO_SHELL_RUNTIME.md custom/README.md
git commit -m "release: publish strict no-shell runtime migration"
```

**Step 4: Push**

Run: `git push origin go`
Expected: remote updated.

---

## Definition of Done Checklist

- [ ] `rg --files | rg '\\.sh$'` returns `0`
- [ ] `go test ./...` passes
- [ ] `go run ./cmd/octo validate --strict-no-shell --dir . --json` returns `status=ok`
- [ ] No `.claude`, `.claude-plugin`, `docs`, `Makefile`, `.github` runtime instructions reference `.sh` execution
- [ ] `/octo:init`, `/octo:plan`, `/octo:develop`, `/octo:deliver`, `/octo:debate`, `/octo:persona` all run through Go runtime only
- [ ] Evidence committed under `docs/plans/evidence/no-shell-runtime/`

## Risks and Mitigations

- Risk: deleting helper `.sh` used only by local maintainer workflows.
  - Mitigation: port essential maintenance flows to Go `cmd/octo-devx`; document removed scripts in release notes.
- Risk: hidden shell invocation in markdown snippets causes operator confusion.
  - Mitigation: strict validator scans docs and command/skill markdown.
- Risk: parity regressions with legacy orchestrator behavior.
  - Mitigation: lock parity fixtures and compare JSON schema + key behavior before deletion.
