# Multipowers Quick Summary

**Multipowers** is a workflow and role orchestration tool for vibe coding with Claude Code and similar AI environments.

## What It Does

- Keeps your project background in stable context files (no repetition per task)
- Routes tasks to **Fast Lane** (direct role execution) or **Standard Lane** (workflow-driven)
- Dispatches work to specialist roles via CLI connectors (`claude`, `codex`, `gemini`)
- Tracks work with a track lifecycle (new → start → execute → complete)
- Enforces governance checks on major changes (security, linting, formatting, docs sync)

## Key Concepts

| Concept | Purpose |
|---------|---------|
| **Fast Lane** | Small, low-risk tasks → direct role execution |
| **Standard Lane** | Significant changes → workflow with multiple nodes and role switches |
| **Track** | A unit of work with start/finish boundaries |
| **Role** | Specialist function (Router, Architect, Coder, Librarian) |
| **Workflow** | Ordered nodes with default executor and optional specialist overrides |

## Typical Commands

```bash
# Setup
./bin/multipowers init
./bin/multipowers doctor

# Start a track
./bin/multipowers track new my-feature
./bin/multipowers track start my-feature

# Execute work
./bin/multipowers run --task "Fix auth bug"              # Fast Lane
./bin/multipowers run --task "Add user API" --risk-hint high  # Standard Lane

# Finish
./bin/multipowers track complete my-feature
```

## What You Get

- Predictable task routing (no guessing which workflow to use)
- Observable execution (structured logs in `outputs/runs/`)
- Specialist role dispatch (Architect for review, Coder for implementation)
- Governance artifacts for major changes (evidence of what was checked and changed)
