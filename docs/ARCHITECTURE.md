# Architecture: Models, Providers, and Execution Flow

This document explains how Multipowers orchestrates multiple AI providers and the execution flow of each workflow.

---

## Overview

Multipowers coordinates **three AI providers** to give you multi-perspective analysis:

```
    +------------------+
    |   Claude Code    |  <-- Orchestrator (you talk to this)
    +--------+---------+
             |
    +--------v---------+
    | Multipowers   |  <-- Plugin coordinates providers
    +--------+---------+
             |
     +-------+-------+
     |       |       |
+----v--+ +--v---+ +-v----+
| Codex | |Gemini| |Claude|
|  CLI  | | CLI  | |  AI  |
+-------+ +------+ +------+
    |         |        |
+---v---+ +---v---+ +--v---+
|OpenAI | |Google | |Anthr.|
| API   | |  API  | | API  |
+-------+ +-------+ +------+
```

---

## Provider → Model Mapping

| Provider | CLI Tool | Underlying Model | Cost Source |
|----------|----------|------------------|-------------|
| **Codex CLI** | `codex exec --model gpt-5.3-codex` | GPT-5.3-Codex (high-capability) | Your `OPENAI_API_KEY` |
| **Gemini CLI** | `gemini -y -m gemini-3-pro-preview` | Gemini 3.0 Pro Preview | Your `GEMINI_API_KEY` |
| **Claude** | Built-in | Claude Sonnet 4.5 / Opus 4.6 | Your Claude Code subscription |

> **Note:** Models are as of February 2026. The mp runtime script uses the latest available models.

### What Each Provider Excels At

| Provider | Strengths | Best For |
|----------|-----------|----------|
| **Codex (OpenAI)** | Code generation, structured output, technical analysis | Implementation approaches, code patterns, API design |
| **Gemini (Google)** | Research synthesis, documentation, broad knowledge | Ecosystem research, best practices, alternative perspectives |
| **Claude** | Strategic synthesis, nuanced analysis, code review | Final synthesis, quality assessment, moderation |

---

## Orchestration Architecture (v8.0+)

The Go-native orchestration engine (v8.0+) replaces legacy shell shims with a robust three-layer pipeline:

1.  **Planner Layer**: Resolves global defaults (`config/orchestration.yaml`), workflow overrides (`config/workflows.yaml`), and task-specific overrides into an immutable `ExecutionPlan`. It decomposes prompts into multiple perspectives for parallel phases.
2.  **Executor Layer**: Executes the plan using a bounded worker pool (Goroutines). It manages lifecycle events, handles context cancellation, and performs one-hop automatic fallbacks via the `internal/policy` dispatcher.
3.  **Synthesizer Layer**: Aggregates parallel outputs. It triggers **Progressive Synthesis** based on completed step thresholds and generates a **Final Synthesis** report using high-capability models.

### Key Components

- **`internal/policy`**: Unified model and executor selection logic.
- **`internal/orchestration`**: The core orchestration runtime (Planner, Executor, Synthesizer).
- **`config/*.yaml`**: Declarative source of truth for orchestration, providers, and workflows.

### Hybrid Mailbox Control Plane (v8.1+)

For isolated external-command runs, orchestration also supports a filesystem-backed control plane:

1. **Mailbox IPC**: JSON envelope messages are written to `mailbox/tmp` and atomically renamed into inbox directories.
2. **Watcher Tier**: A polling `MailboxWatcher` emits high-priority control events for semantic invalidation and structural overlap aborts.
3. **Boundary Gate Tier**: Executor gate logic applies deterministic priority ordering (`abort -> invalidate -> overlap -> requeue -> continue`).
4. **Resume Modes**: Requeue decisions explicitly choose `RESUME_IN_PLACE` or `RESTART_FROM_SCRATCH`; stale artifacts force restart.
5. **Resource Guardrail**: Active worktree slots are capped (`active_worktree_cap`) to prevent unbounded sandbox growth.

### Policy Auto Sync Control Loop (v8.2+)

Runtime now includes a Go-native policy learning loop:

1. Multi-entry raw event ingestion to `.multipowers/policy/autosync/events.raw.*.jsonl`.
2. Universal detector mapping (`branching/workspace/command_contract/risk_profile`).
3. Proposal scoring with confidence/conflict/session gates.
4. Overlay activation with revoke + cooldown controls.
5. Deterministic policy context injection into `/mp:*` and external CLI prompt paths.

---

## Execution Flow by Workflow (v8.0+)

### Discover Phase

**Trigger:** `mp discover X` or `/mp:discover`

```
User Request
     |
     v
+------------------------+
| Orchestration Planner  | <-- Decomposes into perspectives
+-----------+------------+
            |
    +-------+-------+
    |               |
    v               v
+-------+       +-------+
| Step 1|       | Step 2|   <-- Parallel Workers (Goroutines)
| AgentA|       | AgentB|
+---+---+       +---+---+
    |               |
    v               v
"Perspective A" "Perspective B"
    |               |
    +-------+-------+
            |
            v
    +---------------+
    |  Synthesizer  |   <-- Progressive & Final
    +---------------+
            |
            v
      Final Report
```

**New Naming Convention:**
- **Discover** (formerly probe)
- **Define** (formerly grasp)
- **Develop** (formerly tangle)
- **Deliver** (formerly ink)

---

### Define Phase (grasp)

**Trigger:** `mp define X` or `/mp:define`

```
User Request
     |
     v
+--------------------+
|   Multipowers   |
+---------+----------+
          |
          v
    +-----------+
    |   Codex   |   <- Step 1: Problem statement
    +-----------+
          |
          v
    +-----------+
    |  Gemini   |   <- Step 2: Success criteria
    +-----------+
          |
          v
    +-----------+
    |  Gemini   |   <- Step 3: Constraints
    +-----------+
          |
          v
    +-----------+
    |  Gemini   |   <- Step 4: Build consensus
    | Consensus |
    +-----------+
          |
          v
   Problem Definition
    + Requirements
```

**Execution:** (Sequential for coherent problem definition)
1. Codex defines the core problem statement (2-3 sentences)
2. Gemini defines success criteria (3-5 measurable criteria)
3. Gemini defines constraints and boundaries
4. Gemini synthesizes all perspectives into unified requirements

---

### Develop Phase (tangle)

**Trigger:** `mp build X` or `/mp:develop`

```
User Request
     |
     v
+--------------------+
|   Multipowers   |
+---------+----------+
          |
    +-----+-----+
    |           |
    v           v
+-------+   +-------+
| Codex |   |Gemini |   <- PARALLEL: Implementation proposals
+---+---+   +---+---+
    |           |
    v           v
"Approach A" "Approach B"
    |           |
    +-----+-----+
          |
          v
    +----------+
    |  Claude  |
    |  Merge   |
    +----------+
          |
          v
    +----------+
    | Quality  |   <- 75% CONSENSUS GATE
    |   Gate   |
    +----------+
       |     |
   PASS?    FAIL?
       |     |
       v     v
   Continue  Revise
```

**Execution:**
1. Codex and Gemini each propose implementation approaches
2. Claude merges the best elements from both
3. **Quality Gate** checks if merged approach meets 75% consensus threshold
4. If failed: Loop back for revision
5. If passed: Proceed to implementation

**Quality Gate:**
The quality gate is based on subtask success rate:
- Measures: percentage of subtasks that completed successfully
- Threshold: 75% (configurable via `CLAUDE_OCTOPUS_QUALITY_THRESHOLD`)
- If failed: Can retry, escalate to human review, or abort

---

### Deliver Phase (ink)

**Trigger:** `mp review X` or `/mp:deliver`

```
User Request
     |
     v
+--------------------+
|   Multipowers   |
+---------+----------+
          |
    +-----+-----+
    |           |
    v           v
+-------+   +-------+
| Codex |   |Gemini |   <- PARALLEL: Different review angles
+---+---+   +---+---+
    |           |
    v           v
"Code quality""Security &
 review"       edge cases"
    |           |
    +-----+-----+
          |
          v
    +----------+
    |  Claude  |
    | Validate |
    +----------+
          |
          v
    +----------+
    | Quality  |
    |  Score   |
    +----------+
          |
          v
   Validation Report
   + Go/No-Go Decision
```

**Execution:**
1. Codex reviews code quality, patterns, maintainability
2. Gemini reviews security, edge cases, compliance
3. Claude synthesizes into validation report
4. Quality score determines go/no-go recommendation

**Validation Thresholds:**
| Score | Status | Recommendation |
|-------|--------|----------------|
| >= 90% | PASSED | Ship it |
| 75-89% | WARNING | Ship with caution |
| < 75% | FAILED | Do not ship |

---

### Debate (grapple)

**Trigger:** `mp debate X vs Y` or `/mp:debate`

```
User Question
     |
     v
+--------------------+
|   Multipowers   |
+---------+----------+
          |
     Round 1
    +-----+-----+
    |     |     |
    v     v     v
+-----+ +-----+ +-----+
|Codex| |Gemin| |Claud|  <- All 3 PARALLEL
+--+--+ +--+--+ +--+--+
   |       |       |
   v       v       v
"Pro X"  "Pro Y" "Moderator
                  analysis"
   |       |       |
   +---+---+---+---+
       |       |
     Round 2 (optional)
       |
       v
   +-------+
   |Claude |
   |Synth. |
   +-------+
       |
       v
  Final Verdict
  + Recommendation
```

**Execution:**
1. **Round 1:** All three providers argue their positions in parallel
2. **Round 2+ (optional):** Rebuttals and counter-arguments
3. **Synthesis:** Claude moderates and produces final verdict

**Debate Styles:**
| Style | Rounds | Approach |
|-------|--------|----------|
| quick | 1 | Fast positions, immediate synthesis |
| thorough | 2-3 | Multiple rounds of debate |
| adversarial | 3 | Providers actively critique each other |
| collaborative | 2 | Providers build on each other's ideas |

---

### Full Workflow (embrace)

**Trigger:** `/mp:embrace`

```
User Request
     |
     v
+---------+
| DISCOVER|  <- Phase 1
+---------+
     |
     v
+---------+
|  DEFINE |  <- Phase 2
+---------+
     |
     v
+---------+
| DEVELOP |  <- Phase 3 (with quality gate)
+---------+
     |
     v
+---------+
| DELIVER |  <- Phase 4
+---------+
     |
     v
  Complete
  Feature
```

**Execution:**
All four phases run sequentially. Each phase uses the output of the previous phase as context.

**Typical duration:** 2-5 minutes  
**Typical cost:** $0.10-0.30

---

## Cost Breakdown

### Per-Query Estimates

| Workflow | Codex Cost | Gemini Cost | Total |
|----------|------------|-------------|-------|
| discover | $0.01-0.02 | $0.01-0.02 | $0.02-0.04 |
| define | $0.01-0.02 | $0.01-0.02 | $0.02-0.04 |
| develop | $0.02-0.05 | $0.02-0.05 | $0.04-0.10 |
| deliver | $0.01-0.03 | $0.01-0.03 | $0.02-0.06 |
| debate | $0.02-0.05 | $0.02-0.05 | $0.05-0.15 |
| embrace | $0.05-0.10 | $0.05-0.10 | $0.10-0.30 |

**Note:** Claude costs are included in your Claude Code subscription (Pro, Max 5x, Max 20x).

### Cost Optimization

| Strategy | How |
|----------|-----|
| **Use one provider** | Only install Codex OR Gemini (not both) |
| **Skip unnecessary phases** | Use `/mp:develop` instead of `/mp:embrace` for simple tasks |
| **Use Claude-only** | For simple tasks, don't use "mp" prefix - just ask directly |

---

## Provider Detection

Multipowers auto-detects which providers are available:

```bash
# Check status
/mp:setup

# Output example:
# Providers:
#   Codex CLI: ready (OPENAI_API_KEY found)
#   Gemini CLI: ready (OAuth authenticated)
```

### Graceful Degradation

| Available Providers | Behavior |
|--------------------|----------|
| Codex + Gemini | Full multi-AI orchestration |
| Codex only | Dual perspective (Codex + Claude) |
| Gemini only | Dual perspective (Gemini + Claude) |
| Neither | Claude-only mode (basic functionality) |

---

## Visual Indicators

When multi-AI mode is active, you'll see these indicators:

| Indicator | Meaning |
|-----------|---------|
| 🐙 | Multipowers orchestration active |
| 🔴 | Codex CLI executing (OpenAI) |
| 🟡 | Gemini CLI executing (Google) |
| 🔵 | Claude subagent processing |

**Example output:**
```
🐙 CLAUDE OCTOPUS ACTIVATED - Multi-provider research mode
🔍 Discover Phase: Researching authentication patterns

🔴 Codex CLI: Analyzing implementation patterns...
🟡 Gemini CLI: Researching ecosystem best practices...
🔵 Claude: Synthesizing perspectives...

[Final synthesis report]
```

---

## Under the Hood: mp runtime

All workflows are powered by the native Go binary at `.claude-plugin/bin/mp`:

```bash
# Direct CLI usage
./.claude-plugin/bin/mp discover "research OAuth patterns"
./.claude-plugin/bin/mp develop "implement authentication"
./.claude-plugin/bin/mp deliver "review auth code"
./.claude-plugin/bin/mp embrace "complete auth feature"
```

The plugin wraps these commands and provides:
- Natural language triggers
- Session management
- Result storage
- Quality gates

---

## See Also

- **[Visual Indicators Guide](./VISUAL-INDICATORS.md)** - Visual feedback system
- **[Triggers Guide](./TRIGGERS.md)** - What activates each workflow
- **[CLI Reference](./CLI-REFERENCE.md)** - Direct CLI usage
- **[Command Reference](./COMMAND-REFERENCE.md)** - All available commands
