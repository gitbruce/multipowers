# No-Shell Runtime + Hybrid Reasoning Design

Date: 2026-03-02  
Status: Approved

## 1. Context

Current `go` direction solved shell dependency but over-centralized logic in `bin/mp`, while many skills regressed to thin wrappers. This removed explicit LLM reasoning paths from Markdown skills and reduced maintainability.

## 2. Goal

Adopt a hybrid architecture:
- Go runtime provides deterministic atomic capabilities.
- Markdown skills provide stepwise reasoning and orchestration using atomic CLI calls.

## 3. Scope (Big-Bang)

Single-iteration full replacement across three domains:
- `state + validate`
- `hook + route`
- `test + coverage + tdd-env`

No staged rollout. Keep compatibility shims only where required for immediate command continuity.

## 4. Architecture

### 4.1 Engine (Go)

Go runtime is a deterministic execution engine, not a full workflow black box.

Required atomic command surfaces:
- `mp state get|set|update`
- `mp validate --type <tdd-env|test-run|coverage|workspace|no-shell>`
- `mp hook run --event <event>` (or equivalent `mp hook --event` with normalized JSON)
- `mp route --intent <intent> --provider-policy <policy>`
- `mp test run`
- `mp coverage check`

### 4.2 Brain (Markdown Skills)

Skills must be restored from thin wrappers to guided workflows:
- Include explicit stage-by-stage reasoning instructions.
- Call atomic `mp` commands instead of shell scripts.
- Consume JSON outputs to decide next actions.
- Preserve LLM decision space for non-deterministic reasoning.

### 4.3 Contract Layer

All atomic commands must return a normalized response schema:
- `status`
- `action`
- `error_code`
- `message`
- `data`
- `remediation`

This contract is the stable boundary between deterministic runtime and reasoning skills.

## 5. Command and Module Changes

### 5.1 CLI

Refactor `internal/cli/root.go` to atomic-first subcommands and flags. `status` must be de-stubbed and report real runtime health:
- runtime/context completeness
- provider availability
- latest validation state
- hook readiness

### 5.2 Workflows

`internal/workflows` should expose reusable atomic operations. Existing high-level commands (`discover`, `develop`, `deliver`, etc.) become compatibility facades over atomic operations.

### 5.3 Skills

Replace current thin `flow-*.md` and key `skill-*.md` wrappers with structured reasoning content. Remove `.sh` references and map each call to `mp` atomic commands.

### 5.4 Hooks

Use existing `hooks.json` and `internal/hooks` with split responsibilities:
- Deterministic policy checks in Go hooks.
- Guidance-triggering actions via `action` values (for example `ask_user_questions`) interpreted by skills.

## 6. Data Flow

1. Skill parses user intent and determines next step.
2. Skill calls an atomic `mp` command.
3. Go returns normalized JSON contract.
4. Skill branches based on returned `status/action/data`.
5. Skill persists progress through `mp state update`.
6. Loop until completion criteria are met.

## 7. Error Handling

- `status=ok`, `action=continue`: proceed.
- `status=blocked`, `action=ask_user_questions`: collect missing input.
- `status=error`: technical failure with `error_code` and `remediation`.

No command should fail with ambiguous free-form output only.

## 8. Verification Strategy

Required tests for Big-Bang acceptance:
- CLI contract tests for all new atomic commands.
- Regression tests for compatibility commands.
- No-shell runtime scans across skills/commands/hooks.
- End-to-end tests validating `Go JSON result -> Markdown reasoning branch` behavior.

## 9. Acceptance Criteria

- `flow-*` and critical `skill-*` are no longer thin wrappers.
- `mp status` returns real runtime state, not placeholder output.
- Atomic coverage exists for `state/validate/hook/route/test/coverage/tdd-env`.
- Hooks execute deterministically in no-shell runtime.
- Docs and references contain no stale shell-script execution paths.

## 10. Risks and Mitigations

- High blast radius from Big-Bang change.
Mitigation: strict contract tests + compatibility facades + e2e scenarios before merge.

- Behavior drift between old high-level flow semantics and new atomic orchestration.
Mitigation: explicit parity tests for core user journeys.

## 11. Out of Scope

- Reintroducing shell-based control scripts.
- Embedding full reasoning logic back into Go workflows.

## 12. Next Step

Invoke `writing-plans` to produce a detailed implementation plan from this approved design.
