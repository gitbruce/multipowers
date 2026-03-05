# Visual Indicators Guide

Multipowers uses **visual indicators** (emojis) to show you exactly which AI provider is responding at any given moment. This helps you understand whether you're using external CLI tools (which cost money and use your API quotas) or built-in Claude Code capabilities.

## The Indicator System

| Indicator | Meaning | Provider | Cost |
|-----------|---------|----------|------|
| 🐙 | **Parallel Mode** | Multiple CLIs orchestrated | Uses external APIs |
| 🔴 | **Codex CLI** | OpenAI Codex | Your OPENAI_API_KEY |
| 🟡 | **Gemini CLI** | Google Gemini | Your GEMINI_API_KEY |
| 🔵 | **Claude Subagent** | Claude Code Task tool | Included with Claude Code |

---

## What Triggers External CLIs

External CLI providers (Codex and Gemini) are invoked when you:

### 1. Use Explicit Commands
```
/mp:discover "research OAuth patterns"
/mp:debate "Should we use Redis or Memcached?"
```

### 2. Trigger Workflow Skills

Natural language that triggers native `mp` workflows:

- **Discover**: "research X", "explore Y", "investigate Z"
- **Define**: "define requirements for X", "clarify scope of Y"
- **Develop**: "build X", "implement Y", "create Z"
- **Deliver**: "review X", "validate Y", "test Z"

### 3. Use Direct CLI Commands

```bash
# Direct mp runtime execution
mp discover "research GraphQL vs REST"
mp develop "implement user authentication"

# Direct CLI execution
codex exec "Generate API endpoint for users"
gemini -y "What are authentication best practices?"
```

---

## What Triggers Claude Subagents

Claude subagents (built-in Claude Code Task tool) are used for simple file operations, git commands, and single-perspective tasks that don't require the Multipowers Go engine.

---

## Visual Indicator Examples

### Example 1: Research Task (External CLIs)

```
User: Research authentication best practices for React apps

Claude:
🐙 **MULTIPOWERS ACTIVATED** - Multi-provider research mode
🔍 Discover Phase: Researching authentication patterns

Providers:
🔴 Codex CLI - Technical implementation analysis
🟡 Gemini CLI - Ecosystem and community research
🔵 Claude - Strategic synthesis

[Final synthesis report]
```

### Example 2: Code Review (External CLIs)

```
User: Review my authentication code for security issues

Claude:
🐙 **MULTIPOWERS ACTIVATED** - Multi-provider validation
✅ Deliver Phase: Reviewing authentication implementation

Providers:
🔴 Codex CLI - Code quality and best practices
🟡 Gemini CLI - Security audit and edge cases
🔵 Claude - Synthesis and validation report

[Final validation report]
```

---

## Advanced: Hook-Based Indicators

Visual indicators are implemented using Claude Code's hook system via `internal/hooks`.

### PreToolUse Hooks

Whenever a command matching `mp [discover|define|develop|deliver]` is executed, the hook injects the following prompt:
`🐙 **MULTIPOWERS ACTIVATED** - Using external CLI providers`

---

## See Also
- [Triggers Guide](./TRIGGERS.md)
- [CLI Reference](./CLI-REFERENCE.md)
- [Plugin Architecture](./PLUGIN-ARCHITECTURE.md)
