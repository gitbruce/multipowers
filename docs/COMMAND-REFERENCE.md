# Command Reference

Complete reference for all Claude Octopus commands.

---

## Quick Reference

All commands use the `/mp:` namespace.

### System Commands

| Command | Description |
|---------|-------------|
| `/mp:setup` | Check setup status and configure providers |
| `/mp:dev` | Switch to Dev Work mode |
| `/mp:km` | Toggle Knowledge Work mode |
| `/mp:sys-setup` | Full setup command |
| `/mp:model-config` | Configure provider model selection |
| `/mp:persona` | Run a specific persona or list configured personas |

### Workflow Commands

| Command | Phase | Description |
|---------|-------|-------------|
| `/mp:discover` | Discover | Multi-AI research and exploration |
| `/mp:define` | Define | Requirements clarification and scope |
| `/mp:develop` | Develop | Multi-AI implementation |
| `/mp:deliver` | Deliver | Validation and quality assurance |
| `/mp:embrace` | All | Full 4-phase Double Diamond workflow |

### Skill Commands

| Command | Description |
|---------|-------------|
| `/mp:debate` | AI Debate Hub - 3-way debates (Claude + Gemini + Codex) |
| `/mp:review` | Expert code review with quality assessment |
| `/mp:research` | Deep research with multi-source synthesis |
| `/mp:security` | Security audit with OWASP compliance |
| `/mp:debug` | Systematic debugging with investigation |
| `/mp:tdd` | Test-driven development workflows |
| `/mp:docs` | Document delivery (PPTX/DOCX/PDF export) |

### Project Lifecycle Commands

| Command | Description |
|---------|-------------|
| `/mp:status` | Show project progress dashboard |
| `/mp:resume` | Restore context from previous session |
| `/mp:ship` | Finalize project with Multi-AI validation |
| `/mp:issues` | Track issues across sessions |
| `/mp:rollback` | Restore from checkpoint |

---

## Project Lifecycle Commands

Commands for managing project state across sessions.

### `/mp:status`

Show project progress dashboard.

**Usage:** `/mp:status`

**Output:**
- Current phase and position
- Roadmap progress with checkmarks
- Active blockers
- Suggested next action

---

### `/mp:resume`

Restore context from previous session.

**Usage:** `/mp:resume`

**Behavior:**
1. Reads `.claude-octopus/state.json` for current position
2. Loads context using adaptive tier
3. Shows restoration summary
4. Suggests next action

---

### `/mp:ship`

Finalize project with Multi-AI validation.

**Usage:** `/mp:ship`

**Behavior:**
1. Verifies project ready (all phases complete)
2. Runs Multi-AI security audit (Codex + Gemini + Claude)
3. Captures lessons learned
4. Archives project state
5. Creates shipped checkpoint

---

### `/mp:issues`

Track issues across sessions.

**Usage:** `/mp:issues [list|add|resolve|show] [args]`

**Subcommands:**
- `list` - Show all open issues (default)
- `add <description>` - Add new issue
- `resolve <id>` - Mark issue resolved
- `show <id>` - Show issue details

**Issue ID Format:** `ISS-YYYYMMDD-NNN`

**Severity Levels:** critical, high, medium, low

---

### `/mp:rollback`

Restore from checkpoint.

**Usage:** `/mp:rollback [list|<tag>]`

**Subcommands:**
- `list` - Show available checkpoints (default)
- `<tag>` - Rollback to specific checkpoint

**Safety:**
- Creates pre-rollback checkpoint automatically
- Preserves LESSONS.md (never rolled back)
- Requires explicit "ROLLBACK" confirmation

---

## System Commands

### `/mp:setup`

Check setup status and configure AI providers.

**Usage:**
```
/mp:setup
```

**What it does:**
- Auto-detects installed providers (Codex CLI, Gemini CLI)
- Shows which providers are available
- Provides installation instructions for missing providers
- Verifies API keys and authentication

**Example output:**
```
Claude Octopus Setup Status

Providers:
  Codex CLI: ready
  Gemini CLI: ready

You're all set! Try: mp research OAuth patterns
```

### `/mp:km`

Toggle between Dev Work mode and Knowledge Work mode.

**Usage:**
```
/mp:km          # Show current status
/mp:km on       # Enable Knowledge Work mode
/mp:km off      # Disable (return to Dev Work mode)
```

**Modes:**
| Mode | Focus | Best For |
|------|-------|----------|
| Dev Work (default) | Code, tests, debugging | Software development |
| Knowledge Work | Research, strategy, UX | Consulting, research, product work |

### `/mp:dev`

Shortcut to switch to Dev Work mode.

**Usage:**
```
/mp:dev
```

Equivalent to `/mp:km off`.

### `/mp:persona`

Run a specific pre-configured persona, or list available personas.

**Usage:**
```
/mp:persona list
/mp:persona <persona-name> <prompt>
```

**Examples:**
```
/mp:persona list
/mp:persona backend-architect design a scalable webhook ingestion pipeline
/mp:persona security-auditor review auth flow for OWASP risks
```

**What it does:**
- `list` prints personas defined in the repository configuration
- `<persona-name> <prompt>` executes your prompt with that persona profile

---

## Workflow Commands

### `/mp:discover`

Discovery phase - Multi-AI research and exploration.

**Usage:**
```
/mp:discover OAuth authentication patterns
```

**What it does:**
- Launches parallel research using Codex CLI + Gemini CLI
- Synthesizes findings from multiple AI perspectives
- Shows visual indicator: 🐙 🔍

**Natural language triggers:**
- `mp research X`
- `mp explore Y`
- `mp investigate Z`

### `/mp:define`

Definition phase - Clarify requirements and scope.

**Usage:**
```
/mp:define requirements for user authentication
```

**What it does:**
- Multi-AI consensus on problem definition
- Identifies success criteria and constraints
- Shows visual indicator: 🐙 🎯

**Natural language triggers:**
- `mp define requirements for X`
- `mp clarify scope of Y`
- `mp scope out Z feature`

### `/mp:develop`

Development phase - Multi-AI implementation.

**Usage:**
```
/mp:develop user authentication system
```

**What it does:**
- Generates implementation approaches from multiple AIs
- Applies 75% quality gate threshold
- Shows visual indicator: 🐙 🛠️

**Natural language triggers:**
- `mp build X`
- `mp implement Y`
- `mp create Z`

### `/mp:deliver`

Delivery phase - Validation and quality assurance.

**Usage:**
```
/mp:deliver authentication implementation
```

**What it does:**
- Multi-AI validation and review
- Quality scores and go/no-go recommendation
- Shows visual indicator: 🐙 ✅

**Natural language triggers:**
- `mp review X`
- `mp validate Y`
- `mp validate Z`

### `/mp:embrace`

Full Double Diamond workflow - all 4 phases.

**Usage:**
```
/mp:embrace complete authentication system
```

**What it does:**
1. **Discover**: Research patterns and approaches
2. **Define**: Clarify requirements
3. **Develop**: Implement with quality gates
4. **Deliver**: Validate and finalize

Shows visual indicator: 🐙 (all phases)

---

## Skill Commands

### `/mp:debate`

AI Debate Hub - Structured 3-way debates.

**Usage:**
```
/mp:debate Redis vs Memcached for caching
/mp:debate -r 3 Should we use GraphQL or REST
/mp:debate -d adversarial Review auth.ts security
```

**Options:**
| Flag | Description |
|------|-------------|
| `-r N`, `--rounds N` | Number of debate rounds (default: 2) |
| `-d STYLE`, `--debate-style STYLE` | quick, thorough, adversarial, collaborative |

**What it does:**
- Claude, Gemini CLI, and Codex CLI debate the topic
- Claude participates as both debater and moderator
- Produces synthesis with recommendations

**Natural language triggers:**
- `mp debate X vs Y`
- `run a debate about Z`
- `I want gemini and codex to review X`

### `/mp:review`

Expert code review with quality assessment.

**Usage:**
```
/mp:review auth.ts
/mp:review src/components/
```

**What it does:**
- Comprehensive code quality analysis
- Security vulnerability detection
- Architecture review
- Best practices enforcement

### `/mp:research`

Deep research with multi-source synthesis.

**Usage:**
```
/mp:research microservices patterns
```

**What it does:**
- Multi-source research using AI providers
- Documentation lookup via librarian
- Synthesizes findings into actionable insights

### `/mp:security`

Security audit with OWASP compliance.

**Usage:**
```
/mp:security auth.ts
/mp:security src/api/
```

**What it does:**
- OWASP Top 10 vulnerability scanning
- Authentication and authorization review
- Input validation checks
- Red team analysis (adversarial testing)

### `/mp:debug`

Systematic debugging with investigation.

**Usage:**
```
/mp:debug failing test in auth.spec.ts
```

**What it does:**
1. Investigate: Gather evidence
2. Analyze: Root cause identification
3. Hypothesize: Form theories
4. Implement: Fix with verification

### `/mp:tdd`

Test-driven development workflows.

**Usage:**
```
/mp:tdd implement user registration
```

**What it does:**
- Red: Write failing test first
- Green: Minimal code to pass
- Refactor: Improve while keeping tests green

### `/mp:docs`

Document delivery with export options.

**Usage:**
```
/mp:docs create API documentation
/mp:docs export report.md to PPTX
```

**Supported formats:**
- DOCX (Word)
- PPTX (PowerPoint)
- PDF

---

## Visual Indicators

When Claude Octopus activates external CLIs, you'll see visual indicators:

| Indicator | Meaning | Provider |
|-----------|---------|----------|
| 🐙 | Multi-AI mode active | Multiple providers |
| 🔴 | Codex CLI executing | OpenAI (your OPENAI_API_KEY) |
| 🟡 | Gemini CLI executing | Google (your GEMINI_API_KEY) |
| 🔵 | Claude subagent | Included with Claude Code |

**Example:**
```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Multi-provider research mode
🔍 Discover Phase: Researching authentication patterns

Providers:
🔴 Codex CLI - Technical implementation analysis
🟡 Gemini CLI - Ecosystem and community research
🔵 Claude - Strategic synthesis
```

📖 See [Visual Indicators Guide](./VISUAL-INDICATORS.md) for details.

---

## Natural Language Triggers

Instead of slash commands, you can use natural language with the "mp" prefix:

| You Say | Equivalent Command |
|---------|--------------------|
| `mp research OAuth patterns` | `/mp:discover OAuth patterns` |
| `mp build user auth` | `/mp:develop user auth` |
| `mp review my code` | `/mp:deliver my code` |
| `mp debate X vs Y` | `/mp:debate X vs Y` |

**Why "mp"?** Common words like "research" may conflict with Claude's base behaviors. The "mp" prefix ensures reliable activation.

📖 See [Triggers Guide](./TRIGGERS.md) for the complete list.

---

## See Also

- **[Visual Indicators Guide](./VISUAL-INDICATORS.md)** - Understanding what's running
- **[Triggers Guide](./TRIGGERS.md)** - What activates each workflow
- **[CLI Reference](./CLI-REFERENCE.md)** - Direct CLI usage (advanced)
- **[README](../README.md)** - Main documentation
