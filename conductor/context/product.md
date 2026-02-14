# Product: Multipowers Tool

## Mission

Multipowers is an orchestration layer for AI coding workflows. It combines Conductor-style context and track management, Superpowers-style methodology workflows, role-driven node routing, and multi-CLI execution into one operational system for tool maintainers.

The product target is not an end-user app. It is a maintainer toolchain that keeps AI-driven development reproducible, auditable, and fast.

## Users and Scope

- Primary users: maintainers of Multipowers and maintainers of projects bootstrapped by `templates/conductor/`.
- Subject: how to run project delivery through stable context, explicit workflows, and role-specialized execution.
- In-scope languages/runtime: Bash, Python 3.x, Node.js.

## Product Boundaries

- `conductor/` in this repository governs how Multipowers itself is developed.
- `templates/conductor/` defines reusable scaffolding for downstream projects.
- `conductor/context/*.md` are stable background artifacts; track-specific decisions belong in track specs/plans.
- External LLM SDK coupling is out of scope; execution is CLI-connector based.

## Feature Split: Tool Project vs Target Project

### Tool Project Features (this repository)

- Maintainer context is split and explicit: `product.md`, `product-guidelines.md`, `workflow.md`, `tech-stack.md`, `product-vision.md`.
- `setup`/`init` must produce strict governance-ready context for building Multipowers itself.
- Standard-lane workflow enforcement is mandatory for major tool changes.
- Node-level role routing and connector execution are first-class behavior and test-covered.
- Governance and template-sync checks are required before completion claims.

### Target Project Features (generated via templates)

- Target context is intentionally simplified: `product.md` + `tech-stack.md`.
- Workflow and guideline policies are merged into target `product.md` to reduce setup overhead.
- `new track` creates delivery artifacts for app work, while preserving stable context separation.
- Fast-lane and standard-lane routing rules still apply; only documentation shape is simplified.
- Non-interactive role execution (`claude`/`codex`/`gemini`) is reused through the same connector model.

## Product Pillars and Detailed Features

### 1) Project Stabilization via `setup` and `new track` (Conductor baseline)

Multipowers must provide a reliable baseline before implementation work begins.

- `setup` (or equivalent scaffold repair) initializes and validates context by project type:
  - Tool project: `product.md`, `product-guidelines.md`, `workflow.md`, `tech-stack.md`, `product-vision.md`
  - Target project: `product.md`, `tech-stack.md` (with merged workflow/guideline policy in `product.md`)
  - Shared: `conductor/tracks.md`
- `new track` creates a delivery unit with:
  - `conductor/tracks/<track_id>/spec.md`
  - `conductor/tracks/<track_id>/plan.md`
  - `conductor/tracks/<track_id>/metadata.json`
- Track lifecycle states are explicit (`planned`, `in_progress`, `blocked`, `done`) and synchronized into `conductor/tracks.md`.
- Setup and track creation are idempotent: reruns must repair drift without duplicating artifacts.
- Track context is isolated: only stable background stays in context files; change intent stays in `spec.md` and `plan.md`.

### 2) Workflow-first methodology for major changes (Superpowers baseline)

For major changes, Multipowers must enforce explicit methodology workflows rather than ad-hoc role execution.

- Router chooses lane per task:
  - Fast Lane: direct role dispatch for small, low-risk, low-coupling edits.
  - Standard Lane: workflow graph execution for significant changes.
- Standard Lane default workflow chain:
  1. `brainstorming` (intent clarification and design alternatives)
  2. `writing-plans` (task graph with file-level actions and verification steps)
  3. `subagent-driven-development` or `executing-plans` (implementation loop)
  4. `test-driven-development` (RED-GREEN-REFACTOR where applicable)
  5. `requesting-code-review` / `receiving-code-review`
  6. `verification-before-completion`
  7. `finishing-a-development-branch`
- Workflow nodes are checkpointed. A node cannot be marked complete without required artifacts (tests, review notes, verification outputs).
- Major change quality gates are mandatory:
  - changed file inventory (`git diff --name-only`)
  - `semgrep`
  - `biome` for TS/JS
  - `ruff` for Python
  - documentation sync for impacted areas

### 3) Role-driven routing at node level (Role systems baseline)

Multipowers routes each workflow node to the most suitable specialist role, not one fixed role for an entire task.

- Core roles:
  - Router: lane/workflow selection and execution graph control
  - Architect: design, spec validation, review and quality gates
  - Coder: implementation and local refactors
  - Librarian: focused evidence gathering and external research
- Role binding is configuration-driven (`config/roles.default.json` + user overrides).
- Node-to-role mapping is explicit and observable in track artifacts (who executed what, at which node, with which output).
- Review nodes default to Architect unless overridden by policy.
- Router may escalate to specialist roles for high-risk subdomains (security, data migration, API contracts) while preserving workflow ordering.

### 4) Non-interactive multi-CLI bridging (Claude/Codex/Gemini baseline)

Multipowers must execute role calls through non-interactive CLI adapters for `claude`, `codex`, and `gemini`.

- `bin/ask-role` is the single dispatch entrypoint:
  - injects stable context + track context
  - resolves role -> connector mapping
  - executes connector and normalizes output
- Connector requirements:
  - non-interactive execution only (no human prompt loops)
  - deterministic input packaging (prompt template + context bundle)
  - bounded execution (timeouts, retry policy, exit-code handling)
  - structured result envelope (`status`, `summary`, `artifacts`, `raw_output_ref`)
- Bridging pattern follows CLI-to-CLI principles:
  - isolated sub-executions to protect main context window
  - role-specialized prompts for planner/reviewer/implementer style calls
  - final-result return contract (avoid dumping verbose intermediate logs into main thread)
- Optional MCP bridge compatibility is supported as an extension path, but native connector adapters remain first-class.

## Technical Architecture

- `bin/multipowers`: lifecycle commands (`setup`/`init`, `doctor`, `new track`/`track`)
- `bin/ask-role`: routing bridge and context injector
- `connectors/*.py`: CLI adapters (`claude`, `codex`, `gemini`, future adapters)
- `scripts/*.py`: validation and governance checks
- `config/roles.default.json`: role-to-connector policy map

Design constraints:

- Execution metadata must be machine-readable for audits.
- Connectors must degrade safely (clear error envelope, no silent success).
- Workflow state transitions must be replayable from track files.

## Definition of Done (Product-level)

A major track is complete only if all conditions are true:

1. Setup baseline exists and context files are valid.
2. Track spec and plan are present and updated to final status.
3. Workflow nodes show explicit role execution records.
4. Required checks (`semgrep`, `biome`, `ruff`) pass after fixes.
5. Impacted documentation is updated.
6. Completion claim includes evidence artifact references.

## Non-Goals

- No hidden autonomous completion without visible workflow state.
- No mandatory heavyweight workflow for trivial low-risk edits.
- No provider lock-in to a single model vendor.
- No replacing CLI adapters with tightly coupled embedded SDK flows.

## Acceptance Signals

- Maintainers consistently bootstrap with `setup` and deliver through `new track`.
- Major changes follow explicit methodology workflows and produce checkpoints.
- Routing shows workflow-first selection and specialist role dispatch by node.
- `claude`, `codex`, and `gemini` executions run non-interactively via connectors.
- Post-change governance artifacts are present for every major track.
