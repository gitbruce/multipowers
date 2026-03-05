# Command Reference

Complete reference for all Multipowers commands. All commands use the `/mp:` namespace.

---

## System Commands

| Command | Description |
|---------|-------------|
| `/mp:init` | Initialize Multipowers track context via a smart wizard. |
| `/mp:setup` | Check provider status (Codex, Gemini) and configuration. |
| `/mp:status` | Show project progress dashboard and track state. |
| `/mp:route` | Debug current intelligent routing logic for a specific intent. |
| `/mp:doctor` | Run governance diagnostics (16 checks) through the shared doctor engine. |
| `/mp:policy` | Operate autosync policy lifecycle (`sync/stats/gc/tune`). |

---

## Workflow Commands (The Double Diamond)

These commands trigger the Go-native orchestration engine.

| Command | Phase | Natural Language Trigger | Description |
|---------|-------|--------------------------|-------------|
| `/mp:discover` | **Discover** | `mp research X` | Multi-AI research and exploration. |
| `/mp:define` | **Define** | `mp define requirements` | Requirements clarification and scoping. |
| `/mp:develop` | **Develop** | `mp build Y` | Multi-perspective implementation with quality gates. |
| `/mp:deliver` | **Deliver** | `mp review Z` | Final validation and quality assurance. |
| `/mp:embrace` | **All** | N/A | Full 4-phase sequential workflow execution. |

---

## Utility & Skill Commands

| Command | Description |
|---------|-------------|
| `/mp:debate` | Structured 3-way debates between Claude, Gemini, and Codex. |
| `/mp:loop` | Ralph Wiggum iterative executor - loops until task is completed. |
| `/mp:persona` | Run a specific agent persona directly. |
| `/mp:test` | Deprecated in `mp`; use `mp-devx --action suite`. |
| `/mp:coverage` | Deprecated in `mp`; use `mp-devx --action coverage`. |

---

## Detailed Usage

### `/mp:init`
**Usage:** `/mp:init`
Launches the interactive wizard. It will ask about project goals, tech stack, and success criteria to bootstrap the `.multipowers/` governance directory.

### `/mp:status`
**Usage:** `/mp:status`
Reads the current track state and returns:
- Completion percentage
- Active blockers
- Next recommended actions

### `/mp:doctor`
**Usage:** `/mp:doctor [--list|--check-id <id>|--timeout <duration>|--save|--verbose|--json]`

Runs runtime governance checks using the same engine as:
```bash
mp-devx --action doctor ...
```

Key behavior:
- `--list` prints `check_id/purpose/fail_capable`
- default timeout: all checks `30s`, single check `45s`
- non-zero exit only when one or more checks return `fail`

### `/mp:policy`
**Usage:** `/mp:policy <sync|stats|gc|tune> [flags]`

Backed by `mp policy` runtime commands:
- `mp policy sync [--apply] [--ignore-id <id>] [--rollback-id <id>] [--revoke-id <rule_id>]`
- `mp policy stats`
- `mp policy gc`
- `mp policy tune --mode [balanced|accuracy|storage]`

### `/mp:loop`
**Usage:** `/mp:loop --agent <name> "<prompt>"`
Useful for repetitive tasks (e.g., "fix all lint errors"). The engine will automatically retry and correct errors until the agent provides a "completion promise".

---

## See Also
- [CLI Reference (Direct Usage)](./CLI-REFERENCE.md)
- [Architecture Overview](./ARCHITECTURE.md)
