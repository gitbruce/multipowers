# Plugin Architecture - How Multipowers Works

This guide explains the internal architecture of Multipowers for contributors and advanced users.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Claude Code                             │
│  ┌────────────────────────────────────────────────────────┐ │
│  │            Multipowers Plugin                        │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │ │
│  │  │    Skills    │  │    Hooks     │  │   Commands   │ │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘ │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            ↓
        ┌──────────────────────────────────────┐
        │      Multipowers Go Native Engine    │
        └──────────────────────────────────────┘
                            ↓
        ┌──────────────┬────────────┬───────────┐
        │  Codex CLI   │ Gemini CLI │  Claude   │
        │  (OpenAI)    │  (Google)  │ Subagent  │
        └──────────────┴────────────┴───────────┘
```

---

## Component Overview

### 1. Plugin Manifest (plugin.json)

**Location:** `.claude-plugin/plugin.json`

**Purpose:** Defines the plugin metadata, skills, commands, and dependencies.

```json
{
  "name": "mp",
  "version": "8.1.0",
  "description": "Multi-tentacled orchestrator...",
  "skills": [
    "./.claude/skills/parallel-agents.md",
    "./.claude/skills/flow-discover.md",
    "./.claude/skills/flow-define.md",
    "./.claude/skills/flow-develop.md",
    "./.claude/skills/flow-deliver.md",
    "./.claude/skills/debate.md"
  ],
  "commands": "./.claude/commands/"
}
```

---

### 2. Skills System (Double Diamond)

Skills are markdown files with YAML frontmatter that define Claude's behavior for specific tasks.

| Phase | Skill File | Trigger Pattern | Wrapper For (CLI) |
|-------|------------|-----------------|-------------------|
| **Discover** | `flow-discover.md` | "research X" | `mp discover "X"` |
| **Define** | `flow-define.md` | "define requirements for X" | `mp define "X"` |
| **Develop** | `flow-develop.md` | "build X", "implement Y" | `mp develop "X"` |
| **Deliver** | `flow-deliver.md` | "review X", "validate Y" | `mp deliver "X"` |

---

### 3. Hooks System

Hooks inject additional context or execute commands at specific points in the workflow via `internal/hooks`.

| Hook | Purpose |
|------|---------|
| `PreToolUse` | Inject visual indicators and perform safety checks. |
| `PostToolUse` | Apply quality gates and process results. |
| `SessionStart` | Synchronize state and initialize track context. |
| `Stop` | Cleanup and save session summary. |

---

### 4. Native Engine (internal/orchestration)

The Go-native engine replaces legacy shell shims with a robust three-layer pipeline:
1. **Planner**: Resolves prompts into immutable `ExecutionPlan`.
2. **Executor**: Concurrent execution with bounded `Worktree Slots`.
3. **Synthesizer**: Multi-perspective aggregation and Progressive Synthesis.

---

## See Also
- [Architecture Deep Dive](./ARCHITECTURE.md)
- [CLI Reference](./CLI-REFERENCE.md)
