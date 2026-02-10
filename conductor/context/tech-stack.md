# Tech Stack: Multipowers Tooling

## Subject / Audience

- Subject: How to build the Multipowers tool
- Audience: Tool Maintainers

## Languages and Runtime

- Bash: CLI entrypoints and orchestration shell scripts
- Python 3.x: model connectors and validation scripts
- Node.js: plugin helpers, test harness, tooling integration

## Core Components

- `bin/multipowers`: lifecycle commands (`init`, `doctor`, `track`)
- `bin/ask-role`: role dispatch bridge with context injection
- `connectors/codex.py`: non-interactive codex-cli invocation wrapper
- `connectors/gemini.py`: non-interactive gemini-cli invocation wrapper
- `scripts/*.py`: schema/context/evidence checks

## Role and Connector Model

- Role definitions live in `conductor/config/roles.json` (project) or fallback config.
- Each role maps to a tool connector (`system`, `codex`, `gemini`, etc.).
- Router coordinates; connectors execute external CLIs and return output.

## Required Governance Tooling for Major Changes

- `semgrep` for security/pattern scanning
- `biome` for TS/JS lint/format checks
- `ruff` for Python lint/format checks

These checks are post-change quality gates for significant modifications.
