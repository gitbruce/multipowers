# Implementation Plan: Go Big-Bang Migration (Resumable Runbook)

Date: 2026-02-26
Branch target: `go` (from `multipowers`)
Execution mode: one-round big-bang during freeze window
Owner: multipowers maintainers

---

## How to Use This Plan

- Each task has a checkbox and complete execution contract: `Why / What / How / Evidence`.
- If a session stops, resume at the first unchecked task.
- Do not skip prerequisites in `How`.
- Update checkbox status immediately after evidence is produced.

Status values:
- `[ ]` Not started
- `[x]` Completed (evidence captured)

---

## Adopted Suggestions Mapping (tmp/gemini.md + tmp/cc.md)

This section explicitly maps accepted suggestions to executable tasks so any AI tool can continue safely.

1. Go single execution kernel + thin markdown wrappers
- Implemented by: 1.1, 1.2, 6.1, 6.2

2. Unified JSON contract (`status`, `action_required`, `missing_files`, `error_code`)
- Implemented by: 1.3, 2.3

3. `octo context guard --json --auto-init`
- Implemented by: 2.1, 2.2, 2.3

4. Provider interface + registry
- Implemented by: 4.1

5. Interceptor/pipeline chain (`preflight -> runtime -> execute -> post`)
- Implemented by: 1.2, 3.2, 5.5

6. Domain-based package split (avoid new monolith)
- Implemented by: 0.2, 1.1, 1.4

7. Multi-binary proposal changed to single binary subcommands
- Implemented by: 1.1, 6.1

8. JSON IPC simplified to stdout JSON
- Implemented by: 1.3, 2.3

9. Dynamic plugin loading deferred; static extension points first
- Implemented by: 4.1, 1.4

10. No runtime artifact path under `~/.claude-octopus/*`
- Implemented by: 5.4, 7.1, 7.3

11. Add parity and performance acceptance checks
- Implemented by: 8.4, 8.5

12. Explicit workflow module ownership (`internal/workflows/*`)
- Implemented by: 4.4

13. External command execution abstraction (`internal/execx/*`)
- Implemented by: 4.5

14. Stop/SubagentStop hook behavior and schema clarity
- Implemented by: 5.6, 5.7

15. Concurrent state safety for `.multipowers/temp/state.json`
- Implemented by: 7.4

16. Static analysis quality gate for Go code
- Implemented by: 8.6

17. Dual-run migration strategy for parity validation
- Implemented by: 8.7

18. Optional public API package (`pkg/api/*`)
- Implemented by: 6.4

19. Centralized render/format output package (`internal/render/*`)
- Implemented by: 6.5

---

## Mandatory Context Inputs (Any AI Must Read Before Execution)

These are required references before starting any unchecked task:

1. Design baseline:
- `docs/plans/2026-02-26-go-big-bang-migration-design.md`

2. Existing shell implementation (source of behavior truth):
- `bin/octo` (primary runtime behavior)
- `custom/lib/conductor-context.sh` (context readiness contract)
- `custom/lib/proxy-routing.sh` (proxy behavior expectations)
- `scripts/octo providers` (legacy provider routing)
- `scripts/octo state` (legacy state semantics)
- `scripts/octo context` (legacy context semantics)
- `hooks/*.sh` (legacy hook policies)

3. Existing command/skill invocation contracts:
- `.claude/commands/*.md` (especially `init/plan/discover/define/develop/deliver/embrace/research/review/debate`)
- `.claude/skills/flow-*.md`, `.claude/skills/skill-debate.md`, `.claude/skills/skill-code-review.md`

4. Existing tests (behavior guardrails):
- `tests/unit/test-conductor-context-guard.sh`
- `tests/unit/test-init-command-routing.sh`
- `tests/unit/test-spec-commands-init-guard.sh`
- `tests/unit/test-spec-skills-init-guard.sh`
- `tests/integration/test-spec-commands-auto-init.sh`
- debate/proxy/runtime related tests under `tests/`

5. Target-project policy references:
- required context files = `product.md`, `product-guidelines.md`, `tech-stack.md`, `workflow.md`, `tracks.md`, `CLAUDE.md`
- non-required context readiness files = `FAQ.md`, `context/runtime.json`
- runtime artifacts must stay under target `/.multipowers/*`

Execution rule:
- If any ambiguity exists between this runbook and legacy behavior, record decision in evidence and update this plan before proceeding.

---

## Shell-to-Go Migration Decision Matrix (Authoritative)

Each row defines what to do with an existing shell capability.

| Legacy Source | Legacy Functionality | Decision | Go Target | Validation ID |
|---|---|---|---|---|
| `bin/octo` | Command routing (`init/plan/discover/...`) | Modify and replace core runtime | `internal/cli/*`, `internal/app/pipeline.go`, `cmd/octo/main.go` | V-CLI-001 |
| `bin/octo` | Spec-driven context guard + init fallback | Adopt (with 5+CLAUDE contract) | `internal/context/checker.go`, `internal/context/init_runner.go` | V-CTX-001 |
| `bin/octo` | Init rollback on failure | Adopt | `internal/context/init_runner.go` | V-INIT-002 |
| `bin/octo` | Runtime pre-run hooks | Adopt with optional-file semantics | `internal/runtime/config.go`, `internal/runtime/prerun.go` | V-RT-001 |
| `bin/octo` | Provider routing | Modify into typed interface/registry | `internal/providers/provider.go`, `registry.go`, `router.go` | V-PROV-001 |
| `bin/octo` | Debate quorum behavior | Adopt and harden (>=2) | `internal/providers/quorum.go`, workflow debate module | V-DEBATE-001 |
| `bin/octo` | Result/event persistence | Modify to target workspace paths | `internal/tracks/*`, `internal/faq/*` | V-PATH-001 |
| `custom/lib/conductor-context.sh` | Context completeness predicate | Adopt (remove shell dependency) | `internal/context/requirements.go`, `checker.go` | V-CTX-002 |
| `custom/lib/proxy-routing.sh` | Codex/Gemini proxy injection | Adopt with dynamic host detection | `internal/providers/proxy.go` | V-PROXY-001 |
| `scripts/octo providers` | Provider command selection | Modify into registry/adapters | `internal/providers/router.go` | V-PROV-002 |
| `scripts/octo state` | Session/workflow state JSON I/O | Modify into typed state APIs | `internal/tracks/*` + app state module | V-STATE-001 |
| `scripts/octo context` | Context read/merge for prompts | Modify into context loader/summarizer | `internal/context/loader.go`, `summarizer.go` | V-CONTEXT-001 |
| `hooks/*.sh` | Hook governance logic | Replace with Go hook dispatcher | `internal/hooks/*`, `hooks/hooks.json` | V-HOOK-001 |
| `hooks/*.sh` | Boundary checks | Adopt and harden | `internal/fsboundary/*`, `internal/hooks/pre_tool_use.go` | V-BOUNDARY-001 |
| Shell markdown logic | Model-side governance in skills/commands | Replace with thin wrappers | `.claude/commands/*`, `.claude/skills/*` -> `octo ... --json` | V-THIN-001 |
| Legacy `~/.claude-octopus/*` references | Home-dir artifacts | Reject and remove | Target `/.multipowers/*` only | V-PATH-002 |
| Missing in shell (new requirement) | SessionStart stable context injection (5 files + track, <=20 lines each) | New | `internal/hooks/session_start.go`, `internal/context/summarizer.go` | V-HOOK-002 |
| Missing in shell (new requirement) | Deterministic error-code envelope | New | `internal/app/errors.go`, shared response schema | V-ERR-001 |
| Missing in shell (new requirement) | Performance benchmark gate | New | benchmark harness + CI/report scripts | V-PERF-001 |

Decision semantics:
- Adopt: preserve behavior semantics.
- Modify: preserve intent, change implementation and/or tighten policy.
- Replace: supersede legacy mechanism entirely.
- Reject: explicitly not carried forward.
- New: newly introduced capability required by this migration.

---

## Validation ID Catalog (Used Across Tasks)

- `V-CLI-001`: CLI command routing parity.
- `V-CTX-001`: context guard + auto-init + hard-stop parity.
- `V-CTX-002`: required-file contract exactly matches 5+`CLAUDE.md`.
- `V-INIT-002`: failed init rollback behavior.
- `V-RT-001`: runtime pre-run optional file + fail-fast when present.
- `V-PROV-001`: proxy applied on all codex/gemini invocation paths.
- `V-PROV-002`: provider registry routing correctness.
- `V-DEBATE-001`: debate quorum (3->2 continue, <2 fail).
- `V-STATE-001`: state read/write parity and atomicity.
- `V-CONTEXT-001`: context load/summarize behavior parity.
- `V-HOOK-001`: hooks event dispatch and policy execution correctness.
- `V-HOOK-002`: SessionStart injects 5 file summaries + track status, <=20 lines each file.
- `V-BOUNDARY-001`: write boundary blocks `$HOME` and tool-project runtime artifacts.
- `V-THIN-001`: command/skill wrappers are thin and shell-governance-free.
- `V-PATH-001`: all runtime artifacts written under target `/.multipowers/*`.
- `V-PATH-002`: no operational artifact writes under `~/.claude-octopus/*`.
- `V-ERR-001`: deterministic machine-readable error envelope.
- `V-PERF-001`: preflight benchmark thresholds met.

---

## Phase 0 - Branch and Baseline

### Task 0.1 Create `go` branch from current `multipowers` tip
- [x] Task 0.1
- Why: isolate big-bang migration from ongoing multipowers stabilization.
- What: create new working branch `go` from latest local `multipowers` and track remote.
- How:
  1. `git checkout multipowers`
  2. `git pull --ff-only`
  3. `git checkout -b go`
  4. `git push -u origin go`
- Evidence:
  - `git branch --show-current` returns `go`
  - `git status` clean or expected WIP only
  - remote branch exists

### Task 0.2 Initialize Go 1.22 project skeleton
- [x] Task 0.2
- Why: provide stable compilation boundary before migrating behavior.
- What: add `go.mod`, root CLI entry, and initial package directories.
- How:
  1. `go mod init github.com/gitbruce/claude-octopus`
  2. set `go 1.22` in `go.mod`
  3. create `cmd/octo/main.go` and core folders under `internal/*`
  4. add `make build-go` and `make test-go` stubs
- Evidence:
  - `go build ./...` succeeds
  - repository includes expected Go tree

### Task 0.3 Add file-length soft warning (>500 lines)
- [x] Task 0.3
- Why: enforce maintainability constraint without blocking migration velocity.
- What: CI/local check emits warning for Go files over 500 lines.
- How:
  1. add script/check target scanning `*.go`
  2. print warnings and exit 0
  3. wire into CI non-blocking step
- Evidence:
  - check output shows warning format
  - CI step present and non-blocking

### Task 0.4 Snapshot baseline behavior tests
- [x] Task 0.4
- Why: compare post-migration behavior against known multipowers baseline.
- What: run and save key shell/integration test outputs before rewriting logic.
- How:
  1. run selected baseline tests (context guard, init guard, debate path, runtime pre-run)
  2. save outputs under `docs/plans/evidence/go-big-bang/baseline/`
- Evidence:
  - evidence files exist with timestamps and command lines

---

## Phase 1 - Kernel and Pipeline

### Task 1.1 Implement root CLI and subcommand router
- [x] Task 1.1
- Why: all command execution must flow through one Go binary.
- What: `octo <subcommand>` parser and dispatch table.
- How:
  1. implement root command in `internal/cli/root.go`
  2. register subcommands: init/plan/discover/define/develop/deliver/embrace/research/review/debate/hook
  3. add common flags (`--dir`, `--json`, `--track-id`, `--timeout`)
- Evidence:
  - `octo --help` lists all commands
  - each command returns JSON envelope with `status`

### Task 1.2 Implement unified pipeline contract
- [x] Task 1.2
- Why: prevent markdown-layer bypass of mandatory governance.
- What: reusable execution pipeline used by all spec-driven commands.
- How:
  1. create pipeline stages in `internal/app/pipeline.go`
  2. stage order: resolve -> guard -> auto-init -> recheck -> runtime -> execute -> post
  3. integrate per-command executor callbacks
- Evidence:
  - unit tests assert stage order
  - failing guard prevents command body execution

### Task 1.3 Define deterministic error code system
- [x] Task 1.3
- Why: enable machine-parseable remediation in hooks/skills/docs.
- What: stable error code catalog (`E_CTX_MISSING`, `E_INIT_FAILED`, `E_RUNTIME_PRERUN_FAILED`, `E_PROVIDER_QUORUM`, etc.)
- How:
  1. define typed errors in `internal/app/errors.go`
  2. map to JSON response fields: `error_code`, `message`, `remediation`
  3. add tests for code-to-message mapping
- Evidence:
  - command failures produce expected `error_code`

### Task 1.4 Define core interfaces and extension contracts
- [x] Task 1.4
- Why: lock stable boundaries before broad implementation to prevent architectural drift.
- What: define interfaces for guard/provider/gate/hook/boundary and shared JSON envelopes.
- How:
  1. create interface package (or grouped files) with:
     - `ContextGuard`
     - `ProviderDetector`
     - `GovernanceGate`
     - `HookHandler`
     - `BoundaryPolicy`
  2. define shared response structs for command and hook outputs
  3. add compile-time interface conformance tests
- Evidence:
  - interfaces committed with unit tests
  - adapters compile against contracts

---

## Phase 2 - Context and Init

### Task 2.1 Implement required-context checker (5 + CLAUDE)
- [x] Task 2.1
- Why: enforce agreed minimum context contract.
- What: required files are exactly:
  - `product.md`
  - `product-guidelines.md`
  - `tech-stack.md`
  - `workflow.md`
  - `tracks.md`
  - `CLAUDE.md`
- How:
  1. implement `internal/context/requirements.go`
  2. implement `checker.go` returning missing list
  3. include tests for partial/incomplete context
- Evidence:
  - checker tests pass for full and partial contexts

### Task 2.2 Implement auto-init + recheck + hard-stop behavior
- [x] Task 2.2
- Why: stop “init failed but command continues” regressions.
- What: when missing context -> run init -> recheck -> fail hard if still missing.
- How:
  1. build init runner in `internal/context/init_runner.go`
  2. run from pipeline stage
  3. abort command body when init fails
- Evidence:
  - integration test confirms no downstream action after init failure

### Task 2.3 Implement `octo context guard --json --auto-init`
- [x] Task 2.3
- Why: shared primitive for hooks and markdown wrappers.
- What: standalone command returning guard status JSON.
- How:
  1. add `context guard` subcommand
  2. return fields: `status`, `missing_files`, `init_triggered`, `error_code`
- Evidence:
  - sample output documented and validated in tests

### Task 2.4 Preserve init rollback semantics
- [x] Task 2.4
- Why: failed init must not leave corrupted partial artifacts.
- What: rollback created files/folders on init failure.
- How:
  1. snapshot pre-init state
  2. delete newly created artifacts on failure
  3. retain pre-existing user files
- Evidence:
  - rollback integration test passes

---

## Phase 3 - Runtime Preconditions Contract

### Task 3.1 Implement optional runtime config loading
- [x] Task 3.1
- Why: `runtime.json` is optional but actionable when present.
- What: parse `.multipowers/context/runtime.json` if file exists.
- How:
  1. implement config parser in `internal/runtime/config.go`
  2. if file missing, mark `runtime_config_present=false` and continue
- Evidence:
  - tests for missing, valid, malformed file cases

### Task 3.2 Execute pre-run commands with fail-fast semantics
- [x] Task 3.2
- Why: guarantee preconditions are respected before provider execution.
- What: run matching pre-run entries and stop on failure.
- How:
  1. implement matcher and executor in `internal/runtime/prerun.go`
  2. propagate failed command and exit code into JSON error
  3. no command body execution after pre-run failure
- Evidence:
  - integration test proves fail-fast behavior

### Task 3.3 Implement validation gate primitives using target workspace paths
- [x] Task 3.3
- Why: replace shell-era ad-hoc validation checks and remove home-path dependence.
- What: validation gates read only target project workspace (under `/.multipowers/*`).
- How:
  1. implement gate checks in Go (synthesis/required artifacts/state)
  2. prohibit any default path under `~/.claude-octopus/*`
  3. expose gate result in JSON with actionable remediation
- Evidence:
  - tests assert gate lookups use target project paths only
  - tests fail if home-path fallback appears

---

## Phase 4 - Providers and Debate

### Task 4.1 Implement provider interface and registry
- [x] Task 4.1
- Why: replace monolithic shell case routing with extensible typed contracts.
- What: provider interface and static registration for codex/gemini/claude.
- How:
  1. define interface in `internal/providers/provider.go`
  2. implement registry in `registry.go`
  3. add basic provider adapters
- Evidence:
  - registry tests and provider lookup tests pass

### Task 4.2 Unify proxy routing for codex/gemini paths
- [x] Task 4.2
- Why: remove path-specific proxy drift and runtime inconsistency.
- What: all codex/gemini invocations pass through same proxy policy.
- How:
  1. centralize proxy logic in `internal/providers/proxy.go`
  2. dynamic host detection and env injection
  3. verify both direct and debate paths
- Evidence:
  - tests show env/proxy applied on every codex/gemini invocation path

### Task 4.3 Enforce debate quorum policy
- [x] Task 4.3
- Why: debate must continue with 2 providers, fail below 2.
- What: max 3 participants, minimum 2.
- How:
  1. implement quorum evaluator in `quorum.go`
  2. continue on single-provider failure when quorum remains
  3. fail with `E_PROVIDER_QUORUM` when <2
- Evidence:
  - integration tests for 3->2 continue and 2->1 fail

### Task 4.4 Implement explicit workflow modules (`internal/workflows/*`)
- [x] Task 4.4
- Why: avoid routing/business logic collapse into one file and keep domain ownership clear.
- What: create workflow modules for `discover/define/develop/deliver/embrace/debate` orchestration.
- How:
  1. create `internal/workflows/{discover,define,develop,deliver,embrace,debate}.go`
  2. each workflow module accepts typed context + provider router + runtime hooks
  3. keep command layer thin: command only calls workflow entrypoint
- Evidence:
  - workflow modules compile and are invoked by CLI commands
  - workflow-level unit tests exist per module

### Task 4.5 Implement external command abstraction (`internal/execx/*`)
- [x] Task 4.5
- Why: provider/hook execution needs consistent timeout, env, stdout/stderr capture, and cancellation behavior.
- What: central `execx` package for all subprocess invocations.
- How:
  1. define runner API (`Run`, `RunWithTimeout`, `RunJSON`) with context cancellation
  2. enforce standardized result object (`exit_code`, `stdout`, `stderr`, `duration_ms`)
  3. refactor provider and runtime pre-run execution to use `execx` only
- Evidence:
  - no direct `os/exec` use outside `internal/execx/*` (except tests/tools)
  - integration tests verify timeout and cancel semantics

### Task 4.6 Define provider degradation and fallback strategy
- [x] Task 4.6
- Why: prevent inconsistent behavior when one provider is unavailable/slow/failing.
- What: explicit degradation policy for single and multi-provider flows.
- How:
  1. define per-command fallback matrix (single-provider commands vs debate/multi-provider commands)
  2. implement retry/backoff ceilings and downgrade rules
  3. map each degradation path to stable error/warning codes
- Evidence:
  - strategy table documented under evidence folder
  - tests simulate provider failures and validate expected fallback path

---

## Phase 5 - Hooks Migration (Claude Code First-Class)

### Task 5.1 Move hooks to Go event handlers
- [x] Task 5.1
- Why: maximize Claude Code hook strengths and enforce policy before model actions.
- What: `hooks/hooks.json` calls `octo hook --event <Event>`.
- How:
  1. implement handler dispatcher in `internal/hooks/handler.go`
  2. map SessionStart/UserPromptSubmit/PreToolUse/PostToolUse/Stop/SubagentStop
- Evidence:
  - hook simulation tests pass for all configured events

### Task 5.2 Implement SessionStart context injection (5 files + track)
- [x] Task 5.2
- Why: give stable project context every session start.
- What: inject summaries for product/product-guidelines/tech-stack/workflow/CLAUDE and track status.
- How:
  1. implement summarizer in `internal/context/summarizer.go`
  2. enforce per-file summary limit <=20 lines
  3. include active track + checkbox progress + recent failures
- Evidence:
  - SessionStart payload tests validate line limits and fields

### Task 5.3 Enforce UserPromptSubmit preflight for spec-driven commands
- [x] Task 5.3
- Why: prevent entering expensive flows with invalid context.
- What: block spec-driven command start when context guard fails.
- How:
  1. detect spec-driven command intents
  2. run context guard subcommand
  3. return blocking decision + remediation text when invalid
- Evidence:
  - hook tests confirm blocking behavior

### Task 5.4 Enforce PreToolUse boundaries and command safety
- [x] Task 5.4
- Why: stop writes outside allowed target scope and risky command misuse.
- What: inspect write/edit/bash tool calls and enforce policy.
- How:
  1. implement path boundary checks in `internal/fsboundary/*`
  2. reject writes to `$HOME` and tool project for target execution artifacts
  3. attach reason in hook response
- Evidence:
  - boundary tests for allowed/blocked paths

### Task 5.5 Implement PostToolUse event + FAQ/track updates
- [x] Task 5.5
- Why: close governance loop and reduce repeated failures.
- What: classify events, update FAQ, update track status.
- How:
  1. normalize event schema
  2. classify + dedup + rewrite FAQ by error type
  3. mark track milestones
- Evidence:
  - deterministic FAQ regeneration tests

### Task 5.6 Specify `Stop` / `SubagentStop` hook behavior contract
- [x] Task 5.6
- Why: termination hooks must be deterministic to prevent premature exit or dead-end sessions.
- What: explicit allow/block conditions and remediation responses for `Stop` and `SubagentStop`.
- How:
  1. define contract with decision matrix:
     - allow stop when no mandatory checkpoint pending
     - block stop when context/init incomplete, mandatory evidence missing, or critical task active
  2. implement handlers in `internal/hooks/stop.go` and `subagent_stop.go`
  3. add test fixtures for allow/block cases
- Evidence:
  - hook contract document linked in evidence
  - test matrix covers all decision branches

### Task 5.7 Define hooks JSON I/O schema used by Go handlers
- [x] Task 5.7
- Why: any AI/tool integration needs stable schema to parse hook requests/responses reliably.
- What: versioned JSON schema for hook input and output payloads.
- How:
  1. define schema structs for each event in `pkg/api/jsonschema.go` (or equivalent internal package)
  2. include fields: `event`, `session_id`, `cwd`, `tool_name`, `tool_input`, `decision`, `reason`, `remediation`, `metadata`
  3. validate schema in hook handler tests
- Evidence:
  - schema examples committed under docs evidence
  - validation tests pass for all hook events

---

## Phase 6 - Command/Skill Thin Layer

### Task 6.1 Convert `.claude/commands/*` to thin Go wrappers
- [x] Task 6.1
- Why: remove duplicated logic and reduce upstream conflict footprint.
- What: commands only call `octo` binary and render result.
- How:
  1. update command md execution snippets
  2. remove markdown-only governance branches
- Evidence:
  - command invocation smoke tests pass

### Task 6.2 Convert `.claude/skills/*` to thin Go wrappers
- [x] Task 6.2
- Why: avoid model-side divergence from enforced execution policy.
- What: skill contracts call Go commands and parse JSON only.
- How:
  1. replace shell logic blocks with Go command calls
  2. keep instructional text concise and upstream-close
- Evidence:
  - skill smoke tests pass across spec-driven flows

### Task 6.3 Keep upstream-diff discipline
- [x] Task 6.3
- Why: preserve periodic upstream sync with minimal merge conflict.
- What: only necessary wiring deltas in high-churn files.
- How:
  1. prefer custom docs/config for behavior notes
  2. keep `.claude/commands` edits minimal and mechanical
- Evidence:
  - diff review report under `custom/docs/sync/`

### Task 6.4 Add optional public API package (`pkg/api/*`) for integration stability
- [x] Task 6.4
- Why: provide stable typed contracts for hooks, wrappers, and external automation without leaking internals.
- What: define minimal exported types/envelopes in `pkg/api/*` (optional but recommended).
- How:
  1. expose response envelopes and hook schema types
  2. keep implementation in `internal/*`, export types only
  3. add versioning note for backward compatibility
- Evidence:
  - `pkg/api` compiles and is referenced by wrappers/tests
  - compatibility note documented

### Task 6.5 Add centralized output rendering package (`internal/render/*`)
- [x] Task 6.5
- Why: avoid scattered formatting logic and ensure consistent CLI/hook output across commands.
- What: render helpers for banner/table/markdown/json status output.
- How:
  1. implement render functions with deterministic output order
  2. migrate command output formatting to render package
  3. snapshot-test key outputs
- Evidence:
  - output snapshots pass for key commands
  - command packages no longer format output ad hoc

---

## Phase 7 - Filesystem Boundary + FAQ Engine

### Task 7.1 Enforce target-project-only artifact paths
- [x] Task 7.1
- Why: user requirement forbids artifact pollution outside target project.
- What: all generated artifacts under target `/.multipowers/*`.
- How:
  1. centralize path resolver and policy checker
  2. reject prohibited write roots
- Evidence:
  - integration tests validate no writes to `$HOME` or tool project artifacts

### Task 7.2 Implement FAQ synthesis lifecycle
- [x] Task 7.2
- Why: auto-learning loop must be deterministic and maintainable.
- What: classify, dedup, and regenerate FAQ by error type.
- How:
  1. ingest structured failure events
  2. dedup key = `error_type + root_cause + fix`
  3. regenerate `.multipowers/FAQ.md`
- Evidence:
  - golden-file tests for FAQ output

### Task 7.3 Enforce and test no-home/no-tool runtime artifact policy
- [x] Task 7.3
- Why: this is a hard user requirement and a recurring historical regression.
- What: runtime artifacts and command outputs must never write to `$HOME` or tool project during target execution.
- How:
  1. add explicit denylist checks in boundary policy
  2. add integration tests running commands from target repo and asserting no external writes
  3. add static scan guard for forbidden path literals in Go code
- Evidence:
  - integration test logs proving no forbidden writes
  - static scan report with zero forbidden runtime paths

### Task 7.4 Implement concurrent-safe state handling (`.multipowers/temp/state.json`)
- [x] Task 7.4
- Why: hooks and commands can run concurrently; unsafe writes cause state corruption.
- What: atomic, lock-protected read/write strategy for shared state files.
- How:
  1. introduce file lock strategy (cross-platform compatible) for state read/write
  2. perform write-to-temp + fsync + atomic rename
  3. add concurrent writer/read stress tests
- Evidence:
  - stress test report shows no corruption under concurrent operations
  - state schema remains valid after concurrency tests

---

## Phase 8 - Verification and Cutover

### Task 8.1 Port and run critical regression suites
- [x] Task 8.1
- Why: big-bang migration must prove no critical governance regressions.
- What: unit/integration/E2E coverage for core contracts.
- How:
  1. port key shell tests to Go tests where possible
  2. keep necessary shell smoke tests for plugin integration
- Evidence:
  - test matrix report committed to docs evidence folder

### Task 8.2 Validate target project E2E behaviors
- [x] Task 8.2
- Why: acceptance is behavior in real target projects, not only local unit tests.
- What: run `/octo:init`, `/octo:plan`, `/octo:develop`, `/octo:debate` on clean target repo.
- How:
  1. test missing-context -> init -> success flow
  2. test init-failure hard-stop
  3. test debate quorum and proxy application
- Evidence:
  - E2E transcript files under evidence folder

### Task 8.3 Cut over plugin runtime to Go binary
- [x] Task 8.3
- Why: finalize migration and remove shell as execution source of truth.
- What: plugin command paths point to Go runtime.
- How:
  1. update plugin manifest/hooks
  2. verify install/uninstall/update flows
  3. document force cache refresh steps
- Evidence:
  - plugin smoke checks in clean environment pass

### Task 8.4 Run key behavior parity tests (`plan`, `debate`, `embrace`)
- [x] Task 8.4
- Why: preserve user-visible workflow semantics while changing runtime implementation.
- What: compare old shell baseline vs Go runtime for key behavior parity.
- How:
  1. use baseline outputs from Phase 0.4
  2. execute same prompts on Go runtime
  3. verify parity on:
     - generated artifact structure
     - workflow step ordering and hard-stop behavior
     - required track/intent outputs and status transitions
  4. treat strict byte-level markdown identity as non-goal
- Evidence:
  - parity report under `docs/plans/evidence/go-big-bang/parity/`
  - explicit pass/fail matrix for `plan`, `debate`, `embrace`

### Task 8.5 Run performance benchmarks for hook preflight
- [x] Task 8.5
- Why: validate expected startup/latency gains from Go migration.
- What: benchmark `octo hook preflight` response time.
- How:
  1. run hot-path benchmark (warm process cache) and compute p95
  2. run cold-start benchmark and compute p95
  3. targets:
     - hot path p95 <= 50ms
     - cold start p95 <= 120ms
  4. record machine/env details for reproducibility
- Evidence:
  - benchmark report and raw samples in `docs/plans/evidence/go-big-bang/perf/`
  - target threshold pass/fail summary

### Task 8.6 Add Go static analysis gate
- [x] Task 8.6
- Why: big-bang migration increases defect risk; static checks catch issues early.
- What: run `go vet`, staticcheck, and lint policy in CI/local.
- How:
  1. wire `go vet ./...`
  2. add staticcheck (or equivalent) in CI
  3. define severity policy (block on critical, warn on style where needed)
- Evidence:
  - CI reports include static analysis results
  - no critical static findings at cutover

### Task 8.7 Define and run dual-run migration tests (shell vs Go)
- [x] Task 8.7
- Why: validate behavior before full cutover using side-by-side execution.
- What: execute selected scenarios in both runtimes and compare normalized outputs.
- How:
  1. build dual-run harness for `plan`, `debate`, `embrace`, `init-guard`, `runtime-prerun`
  2. normalize non-deterministic fields (timestamps, transient IDs)
  3. emit parity diff report and classify intentional vs unintentional differences
- Evidence:
  - dual-run report under `docs/plans/evidence/go-big-bang/dual-run/`
  - all unintentional diffs resolved before cutover

---

## Phase 9 - Cleanup and Release

### Task 9.1 Remove shell core or reduce to compatibility wrapper
- [x] Task 9.1
- Why: prevent dual-source drift and maintenance confusion.
- What: deprecate legacy shell logic.
- How:
  1. keep minimal wrapper if needed for backward compatibility
  2. remove unused shell paths and tests
- Evidence:
  - no command path depends on old shell core

### Task 9.2 Version bump, release notes, and migration docs
- [x] Task 9.2
- Why: operational rollout needs explicit upgrade guidance.
- What: bump plugin version and update user/tool docs.
- How:
  1. update plugin/version files
  2. add migration notes and known caveats
  3. include cache refresh and verification commands
- Evidence:
  - tagged release commit with complete notes

### Task 9.3 Final sign-off checklist
- [x] Task 9.3
- Why: ensure all acceptance criteria are explicitly validated.
- What: one-page sign-off with pass/fail for each criterion.
- How:
  1. map acceptance criteria to evidence links
  2. require all mandatory checks pass
- Evidence:
  - sign-off doc under `docs/plans/evidence/go-big-bang/`

---

## Rollback Strategy

### Task R1 Pre-migration tag
- [x] Task R1
- Why: instant rollback safety.
- What: create immutable rollback tag before cutover.
- How:
  1. `git tag go-big-bang-precutover-YYYYMMDD`
  2. push tag to remote
- Evidence:
  - tag visible locally and remotely

### Task R2 Runtime fallback toggle during cutover window
- [x] Task R2
- Why: reduce blast radius if latent issues appear post-cutover.
- What: feature flag or wrapper switch to fallback runtime.
- How:
  1. add runtime selector env/config
  2. document emergency fallback commands
- Evidence:
  - fallback drill executed successfully
