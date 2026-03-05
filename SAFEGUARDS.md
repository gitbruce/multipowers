# Multipowers - Safeguards & Critical Configuration

This document outlines critical configuration that must NOT be changed to ensure the stability of the Multipowers Go-native engine and its integration with Claude Code.

## 🔒 Critical: Plugin Name Lock

**Status:** ✅ **LOCKED** - Do not change without breaking all slash commands.

### Current Configuration

```json
// .claude-plugin/plugin.json
{
  "name": "mp"  // ⚠️ LOCKED - Required for /mp:* commands
}
```

### Why This Is Different From Package Name
Claude Code uses the `name` field in `plugin.json` as the command namespace. Changing this would immediately break all user workflows and internal skill references that use the `/mp:` prefix.

### Validation
Ensured via the **`make parity`** command:
```bash
# Verify plugin namespace and command integrity
make parity
```

---

## 🛡️ Sandbox Write Restrictions (Claude Code v2.1.38+)

Claude Code blocks writes to `.claude-plugin/.claude/skills` in sandbox mode.

**Multipowers Design Rule:**
1. **No Dynamic Skill Generation**: Never attempt to write to the plugin's own `.claude` directory at runtime.
2. **External State**: All runtime state and logs MUST be stored in `~/.multipowers/`.

---

## 🛡️ Go Runtime Integrity

### Command Response Contract
All atomic commands in `internal/cli` MUST return a structured JSON response when invoked with `--json`. This contract is consumed by Markdown Reasoning Skills.

**Contract Fields:**
- `status`: "ok", "error", or "blocked"
- `action`: Next action for the LLM ("continue", "ask_user_questions")
- `data`: Arbitrary payload for reasoning
- `message`: Human-readable summary

### Execution Isolation
Mutating commands (Develop, Deliver) MUST use `internal/orchestration/worktree_slots.go` to prevent race conditions and file corruption during parallel AI execution.

---

## 📋 Pre-Release Checklist

Before releasing a new version:
- [ ] Run `make test` - Ensure all Go unit tests pass.
- [ ] Run `make parity` - Verify plugin namespace is still `"mp"`.
- [ ] Run `make test-coverage` - Verify code coverage hasn't regressed.
- [ ] Check version sync in `plugin.json` and `package.json`.

---

**Last Updated:** March 2026
**Status:** Go-native safeguards active ✅
