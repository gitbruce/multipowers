# Native Integration Guide

**Version:** v7.23.0+
**Last Updated:** February 2026

This guide explains how claude-octopus integrates with native Claude Code features introduced in v2.1.20+.

---

## Overview

claude-octopus v7.23.0+ uses a **hybrid approach** that combines:

1. **Native Claude Code features** (where beneficial)
2. **Claude-octopus orchestration** (where multi-AI needed)

**Key principle:** Use the right tool for the job.

---

## Feature Comparison

| Feature | Native Claude Code | Claude-Octopus | When to Use |
|---------|-------------------|----------------|-------------|
| **Task Management** | TaskCreate/TaskUpdate/TaskList | (v7.23.0+) Uses native tools | Always use native (v7.23.0+) |
| **Planning** | EnterPlanMode/ExitPlanMode | /mp:plan with intent contracts | Simple: native, Complex: octopus |
| **State Persistence** | Context summarization | .claude-octopus/state.json | Multi-session projects: octopus |
| **Multi-AI Orchestration** | Not available | Codex + Gemini + Claude | When diverse perspectives needed |
| **Workflows** | Single-phase | Double Diamond (4-phase) | Complex features: octopus |

---

## 1. Task Management Integration (v7.23.0+)

### What Changed

**Before v7.23.0:**
- Used `TodoWrite` tool for task tracking
- Tasks in `.claude/todos.md` markdown files
- No native UI integration

**After v7.23.0:**
- Uses native `TaskCreate`, `TaskUpdate`, `TaskList`, `TaskGet`
- Tasks show in Claude Code's native UI
- Better progress tracking and visualization

### Migration Path

See [MIGRATION-7.23.0.md](../MIGRATION-7.23.0.md) for complete migration guide.

**Quick migration:**
```bash
# Backup existing todos
cp .claude/todos.md .claude/todos.md.backup

# Run migration
~/.claude/plugins/cache/multipowers-plugins/claude-octopus/7.23.0/scripts/migrate-todos.sh

# Verify tasks
/tasks
```

### API Usage

**Creating tasks:**
```javascript
TaskCreate({
  subject: "Implement user authentication",
  description: "Build auth system with JWT tokens and refresh logic",
  activeForm: "Implementing authentication"
})
```

**Updating tasks:**
```javascript
TaskUpdate({
  taskId: "1",
  status: "in_progress"
})

TaskUpdate({
  taskId: "1",
  status: "completed"
})
```

**Listing tasks:**
```javascript
const tasks = TaskList()
const completed = tasks.filter(t => t.status === 'completed')
const pending = tasks.filter(t => t.status === 'pending')
```

### Task Dependencies

Native tasks support dependencies:

```javascript
TaskCreate({
  subject: "Set up database",
  description: "Configure PostgreSQL",
  activeForm: "Setting up database"
})

TaskCreate({
  subject: "Run migrations",
  description: "Create schema",
  activeForm: "Running migrations",
  addBlockedBy: ["1"]  // Blocked by task 1
})
```

---

## 2. Plan Mode Integration (v7.24.0+)

### Hybrid Planning Approach

claude-octopus v7.24.0+ uses **intelligent routing** between native plan mode and octopus workflows.

#### When to Use Native EnterPlanMode

✅ **Use native plan mode for:**
- Single-phase planning (just need a plan)
- Well-defined requirements
- Quick architectural decisions
- When context clearing after planning is OK

**Example:**
```
User: "I need a plan for implementing OAuth"

Claude detects:
- Clear scope ✓
- Single-phase ✓
- No multi-AI needed ✓

→ Suggests: "Use native EnterPlanMode"
```

#### When to Use /mp:plan

✅ **Use /mp:plan for:**
- Multi-AI orchestration (Codex + Gemini + Claude)
- Double Diamond 4-phase execution
- State needs to persist across sessions
- Complex intent capture with routing
- High-stakes decisions requiring multiple perspectives

**Example:**
```
User: "Should we use microservices or monolith?"

Claude detects:
- High-stakes decision ✓
- Multiple perspectives needed ✓
- Requires research ✓

→ Routes to: /mp:debate or /mp:embrace
```

### Routing Logic

```
IF single_phase AND well_defined AND NOT high_stakes:
    → Suggest native EnterPlanMode

IF multi_ai_needed OR complex_scope OR high_stakes:
    → Use /mp:plan with weighted phases

IF decision_between_alternatives:
    → Use /mp:debate
```

### Context Clearing Compatibility

**Native plan mode behavior:**
- `EnterPlanMode` creates isolated planning context
- `ExitPlanMode` clears/summarizes context to save tokens

**How octopus handles this:**
- State persists in `.claude-octopus/state.json`
- Workflows auto-detect context clearing
- Auto-reload state from files
- No information loss

See [State Persistence](#3-state-persistence-v725) below.

---

## 3. State Persistence (v7.25.0+)

### The Problem

Native plan mode's `ExitPlanMode` **clears Claude's memory** to save tokens. This could disrupt multi-phase octopus workflows.

### The Solution

**File-based state management:**

```
.claude-octopus/
├── state.json              # Main state (decisions, metrics, context)
├── context/                # Phase outputs
│   ├── discover-context.md
│   ├── define-context.md
│   ├── develop-context.md
│   └── deliver-context.md
└── summaries/              # Execution summaries
```

**What survives context clearing:**
- ✅ `.claude-octopus/state.json`
- ✅ Phase context files
- ✅ Native tasks (TaskList)
- ✅ Git commits and WIP checkpoints
- ✅ Multi-AI synthesis files

**What gets cleared:**
- ❌ Claude's memory of conversations
- ❌ Workflow phase outputs in memory

**But:** Workflows auto-reload from files.

### Auto-Resume Protocol

**At start of each workflow:**

```bash
# Check if state exists but memory doesn't
if [[ -f .claude-octopus/state.json ]] && [[ -z "${WORKFLOW_CONTEXT_LOADED}" ]]; then
    echo "🔄 Reloading prior session context..."

    # Load state
    state=$("${CLAUDE_PLUGIN_ROOT}/scripts/mp state" read_state)

    # Restore context
    discover_context=$(echo "$state" | python3 -r '.context.discover')
    define_context=$(echo "$state" | python3 -r '.context.define')
    # ... etc

    # Mark as loaded
    export WORKFLOW_CONTEXT_LOADED=true
fi
```

### state.json Structure

```json
{
  "version": "1.0.0",
  "project_id": "unique-hash",
  "current_workflow": "flow-develop",
  "current_phase": "develop",
  "session_start": "2026-02-03T14:30:00Z",
  "decisions": [
    {
      "phase": "define",
      "decision": "React 19 + Next.js 15",
      "rationale": "Modern stack with best DX",
      "date": "2026-02-03",
      "commit": "abc123f"
    }
  ],
  "blockers": [
    {
      "description": "Waiting for API endpoint",
      "phase": "develop",
      "status": "active",
      "created": "2026-02-03"
    }
  ],
  "context": {
    "discover": "researched auth patterns, chose JWT",
    "define": "user wants passwordless magic links",
    "develop": "implementing backend API first",
    "deliver": null
  },
  "metrics": {
    "phases_completed": 2,
    "total_execution_time_minutes": 45,
    "provider_usage": {
      "codex": 12,
      "gemini": 10,
      "claude": 25
    }
  }
}
```

### Resume Example

**Day 1:**
```bash
/mp:embrace "Build authentication system"
→ Runs discover, define phases
→ Saves state to .claude-octopus/state.json
→ User ends session
```

**Day 2 (after context cleared):**
```bash
/mp:resume  # or just continue with /mp:develop
→ Auto-detects context was cleared
→ Loads state.json
→ Restores discover + define findings
→ Continues from where left off
→ No information lost
```

---

## 4. Multi-AI Orchestration

**This is exclusive to claude-octopus** (not available in native Claude Code).

### What It Does

Runs Codex + Gemini + Claude **in parallel**, then synthesizes perspectives:

```
User: "Research authentication patterns"

Claude-Octopus:
├─ 🔴 Codex CLI → Technical implementation analysis
├─ 🟡 Gemini CLI → Ecosystem and community research
└─ 🔵 Claude → Strategic synthesis

→ Synthesizes all 3 perspectives
→ Provides multi-angle recommendation
```

### When to Use Multi-AI

✅ **Use multi-AI orchestration when:**
- High-stakes decisions (architecture, tech stack)
- Need multiple perspectives (security, design trade-offs)
- Broad research coverage (comparing 5+ options)
- Adversarial review (production-critical code)
- Complex implementations (multiple valid approaches)

❌ **Don't use multi-AI for:**
- Simple operations (file edits, basic refactoring)
- Single perspective adequate
- Quick fixes (typos, formatting)
- Cost efficiency priority
- Already know the answer

### Cost Awareness

**External API usage:**
- 🔴 Codex CLI: ~$0.01-0.05 per query (uses OPENAI_API_KEY)
- 🟡 Gemini CLI: ~$0.01-0.03 per query (uses GEMINI_API_KEY)
- 🔵 Claude: Included with Claude Code subscription

**You see cost estimates BEFORE execution:**
```
💰 Estimated Cost: $0.02-0.05
⏱️  Estimated Time: 2-5 minutes
```

---

## 5. Double Diamond Workflows

**This is exclusive to claude-octopus** (not available in native Claude Code).

### What It Is

Proven design methodology with 4 phases:

```
🔍 Discover (probe)  → Research and exploration
🎯 Define (grasp)    → Requirements and scope
🛠️ Develop (tangle)  → Implementation
✅ Deliver (ink)     → Validation and review
```

### When to Use

**Use Double Diamond workflows for:**
- Complex features requiring research → implementation
- High-stakes projects needing validation
- Features where you want multiple AI perspectives
- When you need structured quality gates

**Example:**
```bash
/mp:embrace "Build payment processing"
→ Discover: Research payment gateways, compliance requirements
→ Define: Lock scope (Stripe, PCI compliance, refund handling)
→ Develop: Implement with quality gates
→ Deliver: Security review, validation
```

### Quality Gates

Each phase includes validation:
- 75% consensus threshold (if 2 of 3 AIs disagree, you see the debate)
- Security checks
- Best practices verification
- Performance considerations

---

## 6. Best Practices

### Use Native Features When Appropriate

```
✅ Native TaskCreate/TaskUpdate for task tracking
✅ Native EnterPlanMode for simple planning
✅ Native /tasks command to view tasks
```

### Use Octopus Features When Needed

```
✅ /mp:plan for complex intent capture
✅ /mp:research for multi-AI research
✅ /mp:embrace for complete 4-phase workflows
✅ /mp:debate for high-stakes decisions
```

### State Management

```bash
# Always initialize state at workflow start
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" init_state

# Record decisions
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" write_decision \
  "define" \
  "Use React 19" \
  "Modern features and Server Components"

# Update context after each phase
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" update_context \
  "discover" \
  "Researched auth patterns, recommend JWT"
```

### Multi-Day Projects

```bash
# Day 1
/mp:embrace "Build feature X"
→ Completes discover, define
→ State saved to .claude-octopus/state.json

# Day 2 (new session, context cleared)
/mp:resume  # or just continue
→ Auto-reloads state
→ Continues seamlessly
```

---

## 7. Troubleshooting

### Issue: "TodoWrite tool not found"

**Cause:** You're on v7.23.0+ which removed TodoWrite

**Solution:** See [MIGRATION-7.23.0.md](../MIGRATION-7.23.0.md)

### Issue: "Context keeps clearing"

**Cause:** Native plan mode's ExitPlanMode behavior

**Solution:** This is expected. Octopus auto-reloads from state.json.

### Issue: "Tasks not showing"

**Cause:** May not have migrated to native tasks

**Solution:**
```bash
# Run migration
~/.claude/plugins/cache/multipowers-plugins/claude-octopus/7.23.0/scripts/migrate-todos.sh

# View native tasks
/tasks
```

### Issue: "State not persisting"

**Cause:** octo state not being called

**Solution:**
```bash
# Initialize state
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" init_state

# Verify file exists
ls .claude-octopus/state.json
```

---

## 8. Migration Timeline

| Version | Features | Release |
|---------|----------|---------|
| v7.22.01 | Last version with TodoWrite | Feb 2026 |
| **v7.23.0** | **Task migration + compatibility layer** | **Feb 2026** |
| v7.24.0 | Hybrid plan mode routing | Mar 2026 |
| v7.25.0 | Enhanced state persistence + resume | Apr 2026 |

---

## 9. API Reference

### State Manager

```bash
# Initialize
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" init_state

# Read state
state=$("${CLAUDE_PLUGIN_ROOT}/scripts/mp state" read_state)

# Set workflow
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" set_current_workflow \
  "flow-discover" "discover"

# Record decision
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" write_decision \
  "<phase>" "<decision>" "<rationale>"

# Update context
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" update_context \
  "<phase>" "<summary>"

# Track metrics
"${CLAUDE_PLUGIN_ROOT}/scripts/mp state" update_metrics \
  "phases_completed" "1"
```

### Native Tasks

```javascript
// Create task
TaskCreate({
  subject: "Task description",
  description: "Detailed info",
  activeForm: "Working on task"
})

// Update task
TaskUpdate({
  taskId: "1",
  status: "in_progress" | "completed" | "deleted"
})

// List tasks
const tasks = TaskList()

// Get specific task
const task = TaskGet({ taskId: "1" })
```

---

## 10. Resources

- [Migration Guide](../MIGRATION-7.23.0.md) - Migrate from v7.22.x
- [Integration Plan](../../analysis/NATIVE_INTEGRATION_PLAN.md) - Technical details
- [Architecture](ARCHITECTURE.md) - Overall system design
- [Command Reference](COMMAND-REFERENCE.md) - All commands

---

**Questions?** Open an issue: https://github.com/nyldn/claude-octopus/issues

---

*Native integration: Best of both worlds. 🐙*
