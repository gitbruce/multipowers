# CLI Reference - Direct mp runtime Usage

This guide documents the direct CLI usage of the Go-native `mp` runtime for advanced users and automation scenarios.

---

## Global Options

The following flags are available for almost all commands:

- `--dir <path>`: Specify the target project directory (default: `.`)
- `--prompt "<text>"`: The primary input text or JSON context
- `--json`: Format the output as a machine-readable JSON response
- `--auto-init`: Whether to automatically initialize track context (default: `true`)

---

## Core Commands

### 1. Project & Track Initialization
```bash
mp init --prompt '{"project_name": "...", "tech_stack": "..."}'
```
Runs the initialization wizard or applies a pre-filled JSON prompt to bootstrap the `.multipowers/` governance directory.

### 2. State & KV Management
```bash
mp state get --key <key_name>
mp state set --key <key_name> --value <value>
mp state update --data '{"key": "value"}'
```
Performs atomic read/write operations on the track state. Useful for CI/CD status tracking.

### 3. Workflow Execution (The Double Diamond)
- **Discover**: `mp discover "research topic"`
- **Define**: `mp define "problem criteria"`
- **Develop**: `mp develop "implementation plan"`
- **Deliver**: `mp deliver "final review"`
- **Embrace**: `mp embrace "full feature flow"`

These commands trigger the Go orchestration engine, resolving the execution plan via `internal/orchestration`.

### 4. Intelligent Routing
```bash
mp route --intent [discover|define|develop|deliver]
```
Checks which AI providers (Codex, Gemini, Claude) are available and how the current routing policy would distribute the workload.

### 5. Validation & Guardrails
```bash
mp validate --type [workspace|no-shell|tdd-env|test-run|coverage]
```
Runs specific architectural or environmental checks. Use `--type no-shell` to ensure no legacy Bash dependencies are present in the active path.

Note:
- `mp validate --type no-shell` has moved to `mp-devx --action validate-runtime`.

### 6. Interactive Loops
```bash
mp loop --agent <agent_name> --prompt "instruction" --max-iterations 5
```
Triggers a Ralph Wiggum loop that continues execution until the agent provides a "completion promise" or the iteration limit is reached.

### 7. Runtime Doctor
```bash
mp doctor [--list] [--check-id <id>] [--timeout <duration>] [--save] [--verbose] [--json]
mp-devx --action doctor [same flags]
```

Doctor runs 16 governance checks (9 upstream-compatible + 7 local).  
Execution defaults:
- all checks timeout: `30s`
- single check timeout (`--check-id`): `45s`
- explicit `--timeout` overrides defaults

Exit code:
- non-zero only when at least one check returns `fail`

Save report:
- `--save` writes `.multipowers/doctor/reports/doctor-YYYYMMDD-HHMMSS.json`
- single check writes `.multipowers/doctor/reports/doctor-<check_id>-YYYYMMDD-HHMMSS.json`

### 8. Lifecycle Hooks
```bash
mp hook --event [SessionStart|EnterPlanMode|UserPromptSubmit|PreToolUse|PostToolUse|WorktreeCreate|WorktreeRemove|Stop|SubagentStop] --prompt "..."
```
Dispatches a lifecycle event to the Go hook handler. Returns `ok` or `blocked` based on current governance policies.

### 9. Policy Autosync Operations
```bash
mp policy sync [--apply] [--ignore-id <id>] [--rollback-id <id>] [--revoke-id <rule_id>]
mp policy stats
mp policy gc
mp policy tune --mode [balanced|accuracy|storage]
```
Manages autosync policy lifecycle, observability, garbage-collection, and tuning profiles.

---

## Debugging

Enable verbose Go runtime debugging:
```bash
MP_DEBUG=1 mp <command>
```

---

## See Also
- [Command Reference (Plugin UI)](./COMMAND-REFERENCE.md)
- [Architecture Overview](./ARCHITECTURE.md)
