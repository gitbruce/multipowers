# Tech Stack: Claude Octopus

## Runtime and Languages

- Go (primary): orchestration runtime and operational utilities.
- Markdown/YAML (primary): command, skill, workflow, and policy definitions.
- JavaScript/TypeScript (supporting): helper scripts and token-extraction subsystem.
- Node.js (supporting runtime): package metadata and JS/TS tooling execution.

## Core Execution Components

- `bin/octo`: central command router and workflow engine.
- `cmd/octo-devx`: maintainer/CI helper entrypoint.
- `internal/providers/*`: provider/model routing logic.
- `internal/tracks/*`: workflow state persistence and track status.
- `internal/hooks/*`: lifecycle and governance hooks.

## Command and Skill Surface

- `.claude/commands/*.md`: `/mp:*` command entrypoints.
- `.claude/skills/*.md`: workflow and discipline skills used by commands.
- `workflows/embrace.yaml`: structured workflow graph definition.
- `workflows/schema.yaml`: workflow schema/validation reference.

## Persona and Routing Configuration

- `agents/config.yaml`: persona catalog and routing metadata.
- `agents/personas/*.md`: persona behavior and domain expertise.
- `config/providers/{claude,codex,gemini}/CLAUDE.md`: provider-specific operational context.
- `config/workflows/CLAUDE.md`: workflow-level operational context.

## Routing Policy Baseline

- Codex (`gpt-5.3-codex`) is the default lane for planning, architecture, and important technical decisions.
- Claude Opus (mapped to GLM-5 in this environment) is the default lane for heavy implementation and heavy-token quality audits.
- Claude Sonnet (mapped to GLM-4.7 in this environment) is the default lane for documentation and test-case authoring.
- Gemini (`gemini-3-pro-preview`) is the default lane for external research tasks.
- Lighter quality checks prefer Codex; heavier quality checks prefer Claude Opus.

## Operational Hooks and Safeguards

- `hooks/*.md` and `.claude-plugin/hooks.json`: lifecycle hooks and runtime wiring.
- `SAFEGUARDS.md` and `SECURITY.md`: policy-level security and safe-operation constraints.

## Test and Verification Stack

- `make test`, `make test-smoke`, `make test-unit`, `make test-integration`.
- Go test suites validate routing, workflow contracts, hook behavior, and release integrity.

## External Dependencies

- Claude Code runtime and plugin framework.
- Optional Codex CLI (`@openai/codex`) for OpenAI-backed provider execution.
- Optional Gemini CLI (`@google/gemini-cli`) for Google-backed provider execution.

## Engineering Constraints

- Keep orchestration logic shell-first unless there is a clear maintainability benefit to move code.
- Preserve backward compatibility for command aliases and workflow triggers where feasible.
- Design for graceful provider degradation; single-provider mode must remain operational.
- Keep documentation, commands, skills, and tests synchronized for every behavioral change.
