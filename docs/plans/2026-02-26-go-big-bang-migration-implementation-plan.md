# Implementation Plan: Go Big-Bang Migration

Date: 2026-02-26
Branch target: `go` (from `multipowers`)
Execution mode: one-round big-bang during freeze window

## Phase 0 - Branch and Baseline

- [ ] Create `go` branch from current `multipowers` tip.
- [ ] Add `go.mod` (Go 1.22) and build/test scaffolding.
- [ ] Add CI soft check for file length > 500 lines (warning only).
- [ ] Snapshot baseline behavior tests before migration.

## Phase 1 - Kernel and Pipeline

- [ ] Implement `cmd/octo` root + subcommand wiring.
- [ ] Implement unified pipeline (`resolve -> guard -> init -> runtime -> exec -> post`).
- [ ] Implement deterministic error codes (`E_CTX_MISSING`, `E_INIT_FAILED`, etc.).
- [ ] Add structured JSON response envelope for all subcommands.

## Phase 2 - Context and Init

- [ ] Implement context checker for required files (5 + `CLAUDE.md`).
- [ ] Implement auto-init trigger + re-check + hard-stop semantics.
- [ ] Implement `octo context guard --json --auto-init`.
- [ ] Ensure init failure rollback semantics for generated artifacts.

## Phase 3 - Runtime Contract

- [ ] Implement optional runtime config loader for `.multipowers/context/runtime.json`.
- [ ] Execute configured pre-run hooks when runtime file exists.
- [ ] Enforce fail-fast behavior on pre-run failures.

## Phase 4 - Providers and Debate

- [ ] Implement provider interface and registry.
- [ ] Implement codex/gemini proxy routing via unified provider router.
- [ ] Implement debate quorum policy (max 3, min 2).
- [ ] Add retry/timeout envelopes with actionable error payloads.

## Phase 5 - Hooks Migration

- [ ] Replace shell hook logic with `octo hook --event ...` handlers.
- [ ] Implement `SessionStart` context injection (5 files + track status, <=20 lines each summary).
- [ ] Implement `UserPromptSubmit` preflight enforcement for spec-driven commands.
- [ ] Implement `PreToolUse` boundary gate and command guard.
- [ ] Implement `PostToolUse` FAQ/event/track updates.
- [ ] Implement `Stop` and `SubagentStop` termination guards.

## Phase 6 - Command/Skill Thin Layer

- [ ] Update `.claude/commands/*` to thin Go invocations.
- [ ] Update `.claude/skills/*` to thin Go invocations.
- [ ] Remove duplicated markdown-side governance constraints.
- [ ] Keep upstream text close; restrict diffs to required execution wiring.

## Phase 7 - Filesystem Boundary and FAQ

- [ ] Enforce no business artifact writes to `$HOME` or tool project.
- [ ] Keep target artifacts under `/.multipowers/*`.
- [ ] Implement FAQ classifier/dedup/regenerator by error type.
- [ ] Ensure FAQ update path is deterministic and bounded.

## Phase 8 - Test and Cutover

- [ ] Port critical shell tests to Go unit/integration/E2E suites.
- [ ] Add regression: init failure must hard-stop and block downstream actions.
- [ ] Add regression: context guard must not bypass with partial context.
- [ ] Run target-project E2E verification.
- [ ] Cut over default plugin runtime to Go binary.

## Phase 9 - Cleanup and Release

- [ ] Remove or reduce shell scripts to compatibility wrappers only.
- [ ] Bump plugin version and update release notes.
- [ ] Force plugin cache refresh instructions in docs.
- [ ] Final verification checklist sign-off.

## Rollback Strategy

- [ ] Keep pre-migration tag for immediate rollback.
- [ ] Keep compatibility wrapper path toggle during cutover window.
- [ ] Gate release on mandatory acceptance checks only.

