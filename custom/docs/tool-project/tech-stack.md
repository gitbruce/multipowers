# Tech Stack: Multipowers

## Runtime and Languages

- **Go (1.21+)**: Core orchestration runtime, atomic CLI, and internal packages.
- **Markdown/YAML**: Command and skill reasoning, workflow graphs, and metadata.
- **JSON/TOML**: Configuration and deterministic command contracts.

## Core Execution Components

- `cmd/mp`: Atomic CLI entrypoint for target-project operations.
- `cmd/mp-devx`: Maintainer/CI helper entrypoint.
- `internal/cli`: CLI surface and command routing.
- `internal/providers`: Multi-provider adapters and role-based routing.
- `internal/tracks`: Workflow state persistence and task tracking.
- `internal/hooks`: First-class lifecycle hook handlers (`SessionStart`, etc.).
- `internal/validation`: Go-based quality gates and no-shell runtime enforcement.
- `pkg/api`: Shared JSON schemas and Go type definitions.

## Command and Skill Surface

- `.claude/commands/*.md`: `/mp:*` entrypoints for Claude Code.
- `.claude/skills/*.md`: High-level reasoning and orchestration skills.
- `custom/config/setup.toml`: Template and initialization protocol for `/mp:init`.

## Persona and Routing Configuration

- `custom/config/models.json`: Provider, role-based routing, and fallback lanes.
- `custom/config/persona-lanes.json`: Persona catalog and intent mappings.
- `agents/personas/*.md`: Persona system prompts and domain expertise.

## Routing Policy Baseline

- **Codex** (`gpt-5.3-codex`): Planning, architecture, and important technical decisions.
- **Claude Opus**: Heavy implementation and heavy-token quality audits.
- **Claude Sonnet**: Documentation and test-case authoring.
- **Gemini** (`gemini-3-pro-preview`): External research tasks.

## Operational Hooks and Safeguards

- `internal/hooks/handler.go`: Dispatcher for all `mp hook run` events.
- `SAFEGUARDS.md` and `SECURITY.md`: Operational security policies.

## Test and Verification Stack

- `go test ./...`: Unified Go test suite (unit, integration, parity).
- `scripts/validate-claude-structure.sh`: Parity check between `main` and `go`.
- `internal/devx/structure_validation_test.go`: Automated structure rules enforcement.

## External Dependencies

- Claude Code runtime and plugin framework.
- Codex CLI (`@openai/codex`) and Gemini CLI (`@google/gemini-cli`) for provider execution.

## Engineering Constraints

- **No-Shell Runtime**: Core control flow must remain in Go; avoid shell scripting for logic.
- **Deterministic Contracts**: All internal commands must return normalized JSON fields.
- **Upstream Parity**: Maintain semantic parity with upstream `main` via `COPY_FROM_MAIN` rules.
- **Reasoning/Logic Split**: Keep reasoning in Markdown skills; keep state/validation in Go.
