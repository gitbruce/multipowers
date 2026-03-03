# Config-Driven Model Routing Design

## Goal

Unify all model and executor selection into `config/` so workflow-level models and persona-level models are configured declaratively, compiled into runtime policy, and never hardcoded in commands/skills/runtime logic outside explicit fallback defaults in the policy compiler.

## Confirmed Decisions

1. Source-of-truth config lives in `./config`.
2. Workflow model policy and persona model policy are separated:
- `config/workflows.yaml`
- `config/agents.yaml`
3. Executor and fallback behavior is declarative:
- `config/executors.yaml`
4. `workflows.yaml` supports optional second-level tasks per workflow.
5. If skill/command docs do not name per-task intent, tasks are addressed as `task_1`, `task_2`, ...
6. External executors enforce model hard constraints and may auto-fallback once.
7. Fallback can cross provider and can fallback to Claude Code native execution.
8. Claude Code execution uses model hint (soft constraint), not hard enforcement.
9. Fallback is automatic, no manual confirmation.
10. `/mp:config` controls whether current model and fallback model are shown.
11. Default is to show model routing and fallback information.
12. Build stage must include both policy compilation and Go binary build to `.claude-plugin/bin`.
13. `.claude-plugin` is runtime-only (read-only in steady state), built/copied from development-time sources.

## Architecture Boundary (Dev-Time vs Run-Time)

### Development-Time (mutable)

- Author and review policy in:
- `config/workflows.yaml`
- `config/agents.yaml`
- `config/executors.yaml`
- Build tools validate and compile these files.

### Build-Time (deterministic)

Build pipeline produces runtime artifacts:

1. Policy compile:
- input: `config/*.yaml`
- output: `.claude-plugin/runtime/policy.json`

2. Binary compile:
- output: `.claude-plugin/bin/mp`
- output: `.claude-plugin/bin/mp-devx`

3. Optional parity check:
- ensure runtime bundle includes policy + binaries + command/skill assets.

### Run-Time (read-only)

Runtime loads only generated artifacts from `.claude-plugin`:

- `.claude-plugin/runtime/policy.json`
- `.claude-plugin/bin/*`
- `.claude-plugin/.claude/*`

Commands/skills pass normalized identifiers (workflow/agent/task), not model strings.

## Configuration Model

### `config/workflows.yaml`

Two-level structure:

- level 1: workflow key (`discover`, `define`, `develop`, `deliver`, ...)
- level 2: optional `tasks` mapping (`task_1`, `task_2`, or named intents)

Each workflow/task can define:

- `model`
- `executor_profile`
- `fallback_policy`
- `display_name`

Resolution order:

1. `workflow.tasks.<task>`
2. `workflow.default`
3. executor-level defaults

### `config/agents.yaml`

Persona policy mapping by agent key:

- `model`
- `executor_profile`
- `fallback_policy`
- `permission_mode` (if needed)

### `config/executors.yaml`

Model-to-executor mapping and enforcement:

- `kind: external_cli | claude_code`
- `enforcement: hard | hint`
- `model_patterns` (regex or exact)
- `command_template` (for external executors)
- `fallback` (single-hop policy)

Confirmed policy:

- `claude_code` => `enforcement=hint`
- `external_cli` => `enforcement=hard`
- on hard-fail => one automatic fallback

## Runtime Contract and Data Flow

1. Command/skill/hook provides:
- `scope`: `workflow` or `agent`
- `name`
- `task` (optional)
- `prompt`

2. Runtime resolves execution contract from `policy.json`:
- `requested_model`
- `effective_model`
- `executor_kind`
- `enforcement_mode`
- `fallback_target` (optional)

3. Executor dispatch:
- external executor: execute with explicit model flag (hard)
- claude_code executor: emit/inject hint only (soft)

4. On failure:
- external execution may fallback once according to policy
- fallback may cross provider or move to `claude_code`

5. Telemetry/metadata returned:
- `degraded` (bool)
- `fallback_from`
- `fallback_to`
- `actual_executor`

## Hooks Integration

- `UserPromptSubmit`: resolve and attach model/executor metadata for visibility.
- `PreToolUse`: re-resolve before tool execution to prevent stale routing.
- `SessionStart`: expose policy summary snapshot.

Hooks no longer infer provider by hardcoded model names in markdown content.

## UX / Config Surface

Add `/mp:config` runtime option:

- `show_model_routing: true|false`
- default `true`

When `true`, user-facing output includes current model and fallback details.
When `false`, runtime still records data but suppresses user-facing routing details.

## Validation and Guardrails

### Build-time validation

- schema validation for all `config/*.yaml`
- task-level duplicate detection
- missing executor references
- model pattern conflicts
- fallback chain validation (max one hop)
- cycle prevention

### CI guardrails

- reject hardcoded model strings outside allowed config and test fixtures
- verify `policy build` + binary build success
- verify runtime can load compiled policy

## Risks and Mitigations

1. Risk: old commands/skills still embed model behavior.
- Mitigation: convert commands/skills to identifier-only calls and add lint checks.

2. Risk: fallback behavior becomes opaque.
- Mitigation: always return `degraded/fallback_*` fields and keep `/mp:config` default visible.

3. Risk: dev/runtime drift.
- Mitigation: runtime consumes only compiled `policy.json`; release gate validates checksum and artifact freshness.

## Non-Goals

1. No hard enforcement for Claude Code native model in this phase.
2. No multi-hop fallback chains in this phase.
3. No direct editing in `.claude-plugin` during normal development workflow.

## Acceptance Criteria

1. Workflow/agent routing decisions come from `config/*.yaml` via compiled policy.
2. Non-config hardcoded model references are removed from runtime decision paths.
3. External executor paths enforce model hard constraints with one-hop automatic fallback.
4. Claude Code paths use hint-only routing.
5. `/mp:config` controls routing/fallback visibility and defaults to visible.
6. Build outputs include `.claude-plugin/runtime/policy.json` and `.claude-plugin/bin/*`.
