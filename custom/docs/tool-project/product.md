# Product: Claude Octopus Tool Project

## Mission

Claude Octopus provides structured, multi-provider AI orchestration for Claude Code through reproducible workflows, explicit quality gates, and role-aware execution.

## Users and Scope

- Primary users: maintainers of this repository and operators using `/mp:*` commands in Claude Code.
- Secondary users: teams adopting orchestration patterns from this repo's docs, workflows, and templates.
- Scope of this product context: how this repository evolves and is operated.

Naming contract:
- slash namespace: `/mp:*`
- plugin id: `multipowers`
- marketplace id: `multipowers-plugins`

## What This Repository Delivers

1. Workflow orchestration engine in `scripts/mp`.
2. Command surface in `.claude/commands/*` and skill implementations in `.claude/skills/*`.
3. Persona and routing configuration in `agents/` and `agents/config.yaml`.
4. Provider and workflow configuration in `config/providers/*` and `workflows/embrace.yaml`.
5. Operational hooks and safeguards in `hooks/*`.
6. Verification coverage in `tests/smoke`, `tests/unit`, and `tests/integration`.

## Product Boundaries

- In scope:
  - Multi-AI orchestration behavior.
  - Command/skill UX for dev and knowledge workflows.
  - Provider detection, routing, fallback, and quality gating.
  - State persistence and resume behavior.
- Out of scope:
  - Building a standalone hosted SaaS product.
  - Coupling execution to a single model vendor.
  - Hidden autonomous completion without visible checkpoints.

## Core Capability Areas

### 1) Structured Workflow Execution

- Double Diamond phases are first-class (`discover`, `define`, `develop`, `deliver`) plus full-chain `embrace`.
- `auto` routing maps user intent to workflow entrypoints.
- Major tasks are expected to pass through explicit phase logic and gate checks.

### 2) Multi-Provider Orchestration

- Supports Codex and Gemini via external CLIs plus Claude-native orchestration.
- Provider availability is detected at runtime and execution degrades gracefully.
- Provider roles vary by phase (research, implementation, critique, synthesis).

Default provider/model policy for this repository:
- Planning + architecture + important decision work: Codex (`gpt-5.3-codex`).
- Heavy coding/implementation: Claude Opus (environment-mapped to GLM-5).
- Documentation + test-case generation: Claude Sonnet (environment-mapped to GLM-4.7).
- External research: Gemini (`gemini-3-pro-preview`).
- Quality checks:
  - heavy/high-token: Claude Opus (GLM-5)
  - light/lower-token: Codex

### 3) Persona and Skill System

- Persona catalog in `agents/personas/*` provides domain-specialized behavior.
- Skills in `.claude/skills/*` encode repeatable delivery patterns (review, debug, TDD, research, docs, debate, resume, etc.).
- Routing uses context and intent to select workflows/personas rather than fixed single-agent execution.

### 4) State, Context, and Resume

- Operational state is managed by Go runtime under `.multipowers/*` artifacts.
- Resume/status/issue flows provide continuity across long-running work.
- Stable product context belongs in `conductor/context/*`; transient implementation details belong in workflow/session artifacts.

### 5) Quality and Governance

- Hooks and validation scripts enforce workflow discipline and reduce drift.
- Tests validate command contracts, routing behavior, integrations, and release expectations.
- Docs in `README.md` and `docs/*` must stay synchronized with command behavior.

## Definition of Done (Product-level)

A major change to this repository is complete only when all are true:

1. Behavior is implemented in the relevant scripts/skills/config.
2. Tests are updated or added for changed behavior.
3. Required validations and checks pass.
4. User-facing docs reflect the final behavior.
5. No regression in provider fallback, routing, or state continuity.

## Non-Goals

- Replacing Claude Code itself.
- Forcing every task into heavy multi-phase execution when a lightweight path is sufficient.
- Treating outputs as complete without verification evidence.
