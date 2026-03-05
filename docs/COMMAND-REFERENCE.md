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
| `/mp:test` | Run Go-native test suite and return structured results. |
| `/mp:coverage` | Check project code coverage. |

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

### `/mp:loop`
**Usage:** `/mp:loop --agent <name> "<prompt>"`
Useful for repetitive tasks (e.g., "fix all lint errors"). The engine will automatically retry and correct errors until the agent provides a "completion promise".

---

## See Also
- [CLI Reference (Direct Usage)](./CLI-REFERENCE.md)
- [Architecture Overview](./ARCHITECTURE.md)
