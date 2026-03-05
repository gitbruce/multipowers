# Multipowers - System Instructions

> **Note:** This file provides context when working directly in the multipowers repository (Go branch).

## Core Architecture (Go-Native)

Multipowers has migrated from a Bash-scripted plugin to a **Go-native hybrid engine**.

1. **Go Runtime (Engine)**: The `mp` binary (built from `cmd/mp/main.go`) handles all heavy lifting: orchestration, worktree slot management, execution isolation, and state KV persistence.
2. **Markdown Reasoning (Brain)**: The `commands/` and `skills/` under `.claude-plugin/` provide LLM-friendly reasoning instructions that ultimately invoke the atomic `mp` binary.

### Key Directories
- `internal/orchestration/`: The core concurrent engine (Planner, Executor, Synthesizer).
- `internal/hooks/`: Event-driven lifecycle handlers (SessionStart, PreToolUse, PostToolUse).
- `internal/fsboundary/`: Strict physical isolation preventing agents from escaping the project root.
- `cmd/mp/`: The production entrypoint.
- `cmd/mp-devx/`: The development toolchain (used for parity checks, benchmarking, and test suites).

---

## Visual Indicators (MANDATORY)

When executing Multipowers workflows, you MUST display visual indicators so users know which AI providers are active.

### Required Output Format

**Before starting a workflow**, output this banner:

```
🐙 **MULTIPOWERS ACTIVATED** - [Workflow Type]
[Phase Emoji] [Phase Name]: [Brief description of what's happening]

Providers:
🔴 Codex CLI - [Provider's role in this workflow]
🟡 Gemini CLI - [Provider's role in this workflow]
🔵 Claude - [Your role in this workflow]
```

**Phase emojis by workflow**:
- 🔍 Discover - Research and exploration
- 🎯 Define - Requirements and scope
- 🛠️ Develop - Implementation
- ✅ Deliver - Validation and review
- 🐙 Debate - Multi-AI deliberation
- 🐙 Embrace - Full 4-phase workflow

---

## Workflow Development Guidelines

When modifying or adding features in the `go` branch:

1. **No New Shell Scripts**: All new runtime logic MUST be written in Go under `internal/`. Shell scripts are strictly limited to basic build/deploy wrappers in the `scripts/` root folder.
2. **Deterministic Output**: All atomic `mp` commands MUST return a structured JSON contract (`status`, `action`, `data`, `message`) when the `--json` flag is provided.
3. **Execution Isolation**: Operations that mutate code MUST leverage `internal/orchestration/worktree_slots.go` to ensure parallel agents do not corrupt the working directory.

## Command Ownership Boundary

- `mp` is runtime-only: user/agent execution paths.
- `mp-devx` is ops/devx-only: test/coverage/validation/build/parity workflows.
- Current migration baseline is tracked in `docs/architecture/command-ownership.md`.

## Testing & Parity

Before submitting changes, ensure you validate the system using the `mp-devx` toolchain:

```bash
# Run Go unit tests and smoke tests
make test

# Verify command/skill namespace parity
./mp-devx -action parity
```

---

## Cost Awareness

Always be mindful that external CLIs cost money:
- 🔴 Codex: Requires user's `OPENAI_API_KEY`
- 🟡 Gemini: Requires user's `GEMINI_API_KEY`
- 🔵 Claude: Included with Claude Code subscription

The Go engine's smart router (`mp route`) handles provider selection, but you must ensure visual indicators accurately reflect when external paid APIs are invoked.
