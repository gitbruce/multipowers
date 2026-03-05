# Product: Multipowers Tool Project

## Mission

Multipowers provides structured, multi-provider AI orchestration for Claude Code through reproducible workflows, explicit quality gates, and role-aware execution.

## Users and Scope

- Primary users: maintainers of this repository and operators using `/mp:*` commands in Claude Code.
- Secondary users: teams adopting orchestration patterns from this repo's docs, workflows, and templates.
- Scope of this product context: how this repository evolves and is operated in the `go` branch.

Naming contract:
- slash namespace: `/mp:*`
- plugin id: `mp`
- marketplace id: `multipowers-plugins`

## What This Repository Delivers (Go Branch)

1. **Go Runtime Engine**: Atomic CLI surfaces in `cmd/mp` and `internal/` packages for state, validation, and routing.
2. **Command Surface**: Entrypoints in `.claude/commands/*.md` and reasoning/orchestration in `.claude/skills/*.md`.
3. **Deterministic Contracts**: Go commands return normalized JSON (`status`, `action`, `error_code`, `message`, `data`, `remediation`).
4. **Hook-Driven Control Plane**: First-class lifecycle events (`SessionStart`, `PreToolUse`, `PostToolUse`, `Stop`) for policy enforcement.
5. **Migration Traceability**: Explicit script classification and parity tracking in `docs/architecture/`.

## Product Boundaries

- In scope:
  - Multi-AI orchestration behavior via no-shell hybrid runtime.
  - Go-based state persistence (`.multipowers/` and `internal/tracks`).
  - Provider detection, routing, fallback, and quality gating.
  - Hook-level governance and safeguards.
- Out of scope:
  - Reintroducing shell runtime logic as the primary control plane.
  - Hidden autonomous completion without visible checkpoints.
  - Modifying the `main` branch during `go` branch delivery.

## Core Capability Areas

### 1) Structured Workflow Execution

- Double Diamond phases are first-class (`discover`, `define`, `develop`, `deliver`) plus full-chain `embrace`.
- `mp route --intent <intent>` maps user intent to workflow entrypoints and provider policies.
- Major tasks must pass through explicit phase logic and Go-based gate checks.

### 2) Multi-Provider Orchestration

- Supports Codex, Gemini, and Claude-native execution via Go provider adapters.
- Provider roles vary by phase (research, implementation, critique, synthesis).

Default provider/model policy (via `custom/config/models.json`):
- Planning + architecture + important decision work: Codex (`gpt-5.3-codex`).
- Heavy coding/implementation: Claude Opus (`claude-opus`).
- Documentation + test-case generation: Claude Sonnet (`claude-sonnet`).
- External research: Gemini (`gemini-3-pro-preview`).
- Quality checks:
  - heavy/high-token: Claude Opus
  - light/lower-token: Codex

### 3) Persona and Skill System

- Persona catalog in `agents/personas/*` provides domain-specialized behavior.
- Skills in `.claude/skills/*` encode reasoning patterns while calling Go commands for state/validation.
- Routing uses `internal/providers/router_intent.go` to select lanes based on task intent.

### 4) State, Context, and Resume

- Operational state is managed by Go runtime under `.multipowers/*` artifacts and `internal/tracks`.
- Resume/status/issue flows provide continuity across long-running work.
- Stable product context belongs in `.multipowers/` in target projects.

### 5) Quality and Governance

- Go-based hooks (`internal/hooks`) and validation packages enforce discipline.
- Tests (Go unit/integration) validate command contracts, routing behavior, and release integrity.

## Definition of Done (Product-level)

A major change to this repository is complete only when all are true:

1. Behavior is implemented in the relevant Go packages or Markdown skills.
2. Go tests (`go test ./...`) are updated or added for changed behavior.
3. Hook-level validations and no-shell runtime checks pass.
4. User-facing docs and migration trackers reflect the final behavior.
5. No regression in provider fallback, routing, or state continuity.

## Non-Goals

- Replacing Claude Code itself.
- One-to-one mechanical translation of every shell script (prefer semantic migration to Go).
- Treating outputs as complete without verification evidence.
