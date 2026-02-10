# Multipowers: Workflow + Role Orchestration for Vibe Coding

Multipowers is a **tool project** that helps vibe-coding agents (Claude Code and similar environments) deliver changes in a predictable way: context first, then workflow, then role execution.

## Project Positioning

- **This repo (`conductor/`)** defines how to build and evolve the **Multipowers tool itself**.
- **Template repo (`templates/conductor/`)** scaffolds context for building the **user's target app**.

## Design Sources and What We Absorb

1. **conductor** (`setup` + `new track`): project-level stable context and track lifecycle.
2. **superpowers** (methodology): explicit software workflow for major changes.
3. **oh-my-opencode** (role selection): task-aware role dispatch for execution nodes.
4. **Claudecode-Codex-Gemini** (multi-CLI bridge): non-interactive `claude` / `codex` / `gemini` invocation per role.

## Core Model

### 1) Dual Context Layers

| Layer | Purpose | Subject | Tech Stack | Audience |
|---|---|---|---|---|
| `conductor/context/*` | Build Multipowers itself | How to build the Multipowers tool | Bash, Python, Node.js | Tool maintainers |
| `templates/conductor/context/*` | Initialize user projects | How to build the user's app | React/Python/Go/Rust (or user-chosen stack) | End users |

### 2) Router: Fast Lane vs Standard Lane

Router is the coordinator. It routes by task type:

- **Fast Lane**: skip skills; directly dispatch an execution role for small, low-risk, bounded tasks.
- **Standard Lane**: choose a **workflow**, not a role. Example workflows: `brainstorming`, `writing-plans`, `subagent-driven-development`, `executing-plans`.

In Standard Lane, each workflow has a default executor role, but specific nodes can require specialist roles. Example: in `subagent-driven-development`, implementation may default to Coder, while code review nodes are forced to Architect.

### 3) Role-to-CLI Mapping

- `router` → local/system coordinator in main Claude Code session
- `architect` → `gemini` CLI (planning/architecture + review/verification)
- `coder` → `codex` CLI (implementation)
- `librarian` → `gemini` CLI (fast research)

All external calls go through `bin/ask-role`, which injects context and forwards role-specific prompts to connectors.

## Major Change Governance (Required)

For significant modifications, use this checklist:

1. **Record changed files** (e.g., `git diff --name-only`).
2. **Run post-change checks**:
   - `semgrep` for security/pattern issues
   - `biome` for TS/JS formatting/linting
   - `ruff` for Python linting/format checks
3. **Fix findings** and re-run until clean (or explicitly justify exceptions).
4. **Update documentation based on changed-file scope** (design/workflow/usage docs).

## Quick Start

```bash
npm install
./bin/multipowers init --repair
./bin/multipowers doctor
./bin/multipowers track new my-feature
```

## Track Lifecycle

```bash
./bin/multipowers track new <feature-name>
./bin/multipowers track start <track-name>
# execute via fast lane or standard lane
./bin/multipowers track complete <track-name>
```

## Key Commands

```bash
# choose lane only (decision)
./bin/multipowers route --task "Refactor auth boundary" --risk-hint high --json

# route + execute in one command
./bin/multipowers run --task "Refactor auth boundary" --risk-hint high --json --allow-untracked

# run standard-lane workflow (node-level role switching is defined in workflow config)
./bin/multipowers workflow run subagent-driven-development --task "Implement feature X" --allow-untracked

# inspect and validate available workflows
./bin/multipowers workflow list
./bin/multipowers workflow validate

# role dispatch bridge
./bin/ask-role architect "Brainstorm architecture"
./bin/ask-role coder "Implement task with TDD"
./bin/ask-role architect "Review for blocking issues"

# health and governance baseline
./bin/multipowers doctor
bash scripts/run_governance_checks.sh --mode strict --changed-file README.md
bash scripts/run_governance_checks.sh --mode advisory --changed-file bin/multipowers --changed-file README.md
npm run governance -- --help

# safe update workflow
./bin/multipowers update --check --json
./bin/multipowers update --apply --yes

npm test --silent
```

## Runtime Observability

- Structured logs are written to `outputs/runs/YYYY-MM-DD.jsonl`.
- Router/workflow events include: `lane_selected`, `fast_lane_dispatched`, `fast_lane_finished`, `workflow_started`, `workflow_node_executed`, `workflow_finished`, `governance_started`, `governance_finished`.
- Use `request_id` (and optional `track_id`) to reconstruct one execution timeline.

## Repository Layout

```text
multipowers/
├── bin/                     # multipowers + ask-role entry points
├── connectors/              # codex/gemini/claude wrappers
├── config/                  # default role and schema config
├── conductor/               # maintainer-facing live context/tracks
├── templates/conductor/     # scaffold for user project context/tracks
├── skills/                  # methodology skills
├── scripts/                 # validators/governance checks
└── tests/                   # regression and integration tests
```
