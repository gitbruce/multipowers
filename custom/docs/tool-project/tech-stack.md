# Tech Stack: Claude Octopus

## Runtime and Languages

- Bash (primary): orchestration core and operational scripts.
- Markdown/YAML (primary): command, skill, workflow, and policy definitions.
- JavaScript/TypeScript (supporting): helper scripts and token-extraction subsystem.
- Node.js (supporting runtime): package metadata and JS/TS tooling execution.

## Core Execution Components

- `scripts/orchestrate.sh`: central command router and workflow engine.
- `scripts/provider-router.sh`: provider/model routing logic.
- `scripts/state-manager.sh`: workflow state persistence helpers.
- `scripts/octo-state.sh`: project/session state commands and context tiers.
- `scripts/agent-teams-bridge.sh`: coordination bridge for multi-agent execution.

## Command and Skill Surface

- `.claude/commands/*.md`: `/octo:*` command entrypoints.
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

- `hooks/*.sh` and `hooks/*.md`: quality gates, lifecycle hooks, status indicators, and safety checks.
- `SAFEGUARDS.md` and `SECURITY.md`: policy-level security and safe-operation constraints.

## Test and Verification Stack

- `make test`, `make test-smoke`, `make test-unit`, `make test-integration`.
- Shell-based suites in `tests/` validate routing, workflow contracts, hook behavior, and release integrity.

## External Dependencies

- Claude Code runtime and plugin framework.
- Optional Codex CLI (`@openai/codex`) for OpenAI-backed provider execution.
- Optional Gemini CLI (`@google/gemini-cli`) for Google-backed provider execution.

## Engineering Constraints

- Keep orchestration logic shell-first unless there is a clear maintainability benefit to move code.
- Preserve backward compatibility for command aliases and workflow triggers where feasible.
- Design for graceful provider degradation; single-provider mode must remain operational.
- Keep documentation, commands, skills, and tests synchronized for every behavioral change.
