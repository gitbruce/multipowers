# Hybrid Skills Restoration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Restore flow-* and skill-* files from thin wrappers to structured reasoning content using atomic mp commands instead of shell scripts.

**Architecture:** Hybrid architecture where Go runtime provides deterministic atomic capabilities (state, validate, route, test, coverage) and Markdown skills provide stepwise reasoning and orchestration using atomic CLI calls.

**Tech Stack:** Go 1.21+, Markdown skills, mp CLI, JSON response contract

---

## Prerequisites

Before starting, verify:
- [ ] `mp` binary exists at `.claude-plugin/bin/mp`
- [ ] Atomic commands available: `state get|set|update`, `validate --type`, `route --intent`, `status`
- [ ] Current branch is `go`

---

## Response Contract Reference

All atomic commands return this structure:
```json
{
  "status": "ok" | "error" | "blocked",
  "action": "continue" | "ask_user_questions",
  "message": "Human-readable status",
  "error_code": "ERR_XXX",
  "data": { ... },
  "remediation": "Suggested fix if blocked"
}
```

---

## P0: Core Skills (Critical Path)

### Task 1: Restore flow-discover.md

**Files:**
- Modify: `.claude-plugin/.claude/skills/flow-discover.md`

**Step 1: Read current thin wrapper**

Run: Read the current file
Expected: Should be ~92 lines with atomic mp commands

**Step 2: Replace with full reasoning content**

Write the restored skill with:
- Trigger section for automatic activation
- 7-step execution contract
- Atomic mp commands instead of shell scripts
- JSON response parsing for branching

```markdown
---
name: flow-discover
aliases:
  - discover
  - discover-workflow
  - probe
  - research
description: Multi-AI research using Codex and Gemini CLIs (Double Diamond Discover phase)
agent: Explore
context: fork
task_management: true
execution_mode: enforced
pre_execution_contract:
  - visual_indicators_displayed
validation_gates:
  - mp_command_executed
  - synthesis_complete
trigger: |
  AUTOMATICALLY ACTIVATE when user requests research or exploration:
  - "research X" or "explore Y" or "investigate Z"
  - "what are the options for X"
  - Comparative analysis ("compare X vs Y")
---

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

### STEP 1: Validate Workspace (MANDATORY)

**Execute via Bash:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json
```

**Parse JSON response:**
- If `status=ok`: Proceed to STEP 2
- If `status=error`: Check `message`, guide user to `/mp:init`
- If `status=blocked`: Review `data.missing`, collect required context

**DO NOT PROCEED until workspace validated.**

---

### STEP 2: Route Providers (MANDATORY)

**Execute via Bash:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent discover --dir "$PWD" --json
```

**Parse JSON response:**
- `data.selected_providers`: List of providers to use
- `data.available_providers`: All available providers
- `data.reason`: Why these providers selected

**Display visual indicators:**
```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Multi-provider research mode
🔍 Discover Phase: [Brief description]

Providers:
🔴 Codex CLI: [from data.available_providers]
🟡 Gemini CLI: [from data.available_providers]
🔵 Claude: Available ✓ (Strategic synthesis)
```

**DO NOT PROCEED until banner displayed.**

---

### STEP 3: Update State (MANDATORY)

**Execute via Bash:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"discover","stage":"executing","started_at":"<timestamp>"}' --dir "$PWD" --json
```

**Verify `status=ok` before proceeding.**

---

### STEP 4: Execute Discovery (MANDATORY - Multi-Provider Research)

Based on routing results, perform multi-perspective research:

1. **Use selected providers** from STEP 2
2. **Synthesize findings** from multiple perspectives
3. **Structure results** in coherent format

**For each selected provider (Codex/Gemini):**
- Call provider CLI with research question
- Capture response
- Include in synthesis

**Your role (Claude):**
- Strategic synthesis
- Pattern identification
- Recommendation formulation

---

### STEP 5: Persist Findings (MANDATORY)

**Execute via Bash:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"discover","status":"complete","findings":"<summary>","completed_at":"<timestamp>"}' --dir "$PWD" --json
```

---

### STEP 6: Present Results

Format findings based on context:

**For Dev Context:**
- Technical research summary
- Recommended implementation approach
- Library/tool comparison
- Next steps

**For Knowledge Context:**
- Strategic research summary
- Business rationale
- Framework analysis
- Next steps

**Include attribution:**
```
---
*Multi-AI Research powered by Claude Octopus*
*Full state: mp state get --dir $PWD*
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Workspace invalid | Guide user to `/mp:init` |
| No providers available | Check `mp status` for provider config |
| State update fails | Retry with simpler data |
| All providers fail | Report error, do not substitute with direct research |

## Context Detection

**Dev Context Indicators:**
- Technical terms: "API", "endpoint", "database", "implementation"
- Project has `package.json`, `Cargo.toml`, etc.

**Knowledge Context Indicators:**
- Business terms: "market", "ROI", "stakeholders", "strategy"
- Research terms: "literature", "synthesis", "academic"
```

**Step 3: Verify file updated**

Run: `wc -l .claude-plugin/.claude/skills/flow-discover.md`
Expected: ~150-200 lines (not thin wrapper ~92 lines)

**Step 4: Commit**

```bash
git add .claude-plugin/.claude/skills/flow-discover.md
git commit -m "feat(skills): restore flow-discover with reasoning content

- Replace thin wrapper with full execution contract
- Use atomic mp commands (validate, route, state)
- Add JSON response parsing for branching logic
- Include visual indicators and error handling"
```

---

### Task 2: Restore skill-deep-research.md

**Files:**
- Modify: `.claude-plugin/.claude/skills/skill-deep-research.md`

**Step 1: Read current thin wrapper**

Run: Read the current file
Expected: Should be ~10 lines (thin wrapper)

**Step 2: Replace with full reasoning content**

```markdown
---
name: skill-deep-research
aliases:
  - research
  - deep-research
description: Deep multi-AI parallel research with cost transparency and synthesis
agent: Explore
context: fork
task_management: true
execution_mode: enforced
pre_execution_contract:
  - interactive_questions_answered
  - visual_indicators_displayed
validation_gates:
  - mp_command_executed
  - synthesis_complete
trigger: |
  Use when user wants deep research: "research this topic", "investigate how X works",
  "analyze the architecture", "explore approaches to Y", "what are the options for Z".
---

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

### STEP 1: Interactive Questions (BLOCKING)

**You MUST call AskUserQuestion BEFORE any other action.**

```javascript
AskUserQuestion({
  questions: [
    {
      question: "How deep should the research go?",
      header: "Research Depth",
      multiSelect: false,
      options: [
        {label: "Quick overview (Recommended)", description: "1-2 min, surface-level"},
        {label: "Moderate depth", description: "2-3 min, standard"},
        {label: "Comprehensive", description: "3-4 min, thorough"},
        {label: "Deep dive", description: "4-5 min, exhaustive"}
      ]
    },
    {
      question: "What's your primary focus area?",
      header: "Primary Focus",
      multiSelect: false,
      options: [
        {label: "Technical implementation (Recommended)", description: "Code patterns, APIs"},
        {label: "Best practices", description: "Industry standards"},
        {label: "Ecosystem & tools", description: "Libraries, community"},
        {label: "Trade-offs & comparisons", description: "Pros/cons analysis"}
      ]
    },
    {
      question: "How should the output be formatted?",
      header: "Output Format",
      multiSelect: false,
      options: [
        {label: "Detailed report (Recommended)", description: "Comprehensive write-up"},
        {label: "Summary", description: "Concise findings"},
        {label: "Comparison table", description: "Side-by-side analysis"},
        {label: "Recommendations", description: "Actionable next steps"}
      ]
    }
  ]
})
```

**Capture responses as:**
- `depth_choice` = user's depth selection
- `focus_choice` = user's focus selection
- `format_choice` = user's format selection

**DO NOT PROCEED TO STEP 2 until all questions answered.**

---

### STEP 2: Validate & Route (MANDATORY)

**Execute via Bash:**
```bash
# Validate workspace
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json

# Route providers
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent discover --dir "$PWD" --json
```

**Display visual indicators:**
```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Multi-provider research mode
🔍 Deep Research: [Brief description]

Provider Availability:
🔴 Codex CLI: [status from route response]
🟡 Gemini CLI: [status from route response]
🔵 Claude: Available ✓ (Strategic synthesis)

Research Parameters:
📊 Depth: ${depth_choice}
🎯 Focus: ${focus_choice}
📝 Format: ${format_choice}

💰 Estimated Cost: $0.01-0.05
⏱️  Estimated Time: 2-5 minutes
```

**If BOTH external providers unavailable:**
- Suggest `/mp:setup` and STOP
- DO NOT proceed with single-provider

---

### STEP 3: Execute Research (MANDATORY)

**Update state:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"research","stage":"executing","depth":"${depth_choice}","focus":"${focus_choice}"}' --dir "$PWD" --json
```

**Perform multi-provider research:**
1. Use selected providers from STEP 2
2. Apply depth/focus parameters to queries
3. Synthesize all perspectives

---

### STEP 4: Persist & Present (MANDATORY)

**Update state with results:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"research","status":"complete","completed_at":"<timestamp>"}' --dir "$PWD" --json
```

**Format according to `format_choice`:**
- **Summary**: 2-3 paragraph overview with key recommendations
- **Detailed report**: Full synthesis with all perspectives
- **Comparison table**: Side-by-side analysis in markdown table
- **Recommendations**: Actionable next steps with rationale

**Include attribution:**
```
---
*Multi-AI Research powered by Claude Octopus*
*Providers: 🔴 Codex | 🟡 Gemini | 🔵 Claude*
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Questions skipped | Use defaults: Quick overview, Technical, Detailed report |
| No providers | Suggest `/mp:setup`, STOP |
| State update fails | Retry, continue if persistent |
| Provider timeout | Continue with available providers |

## Security: External Content

When research fetches external URLs:
1. Validate URL (HTTPS only, no localhost/private IPs)
2. Wrap fetched content in security frame
3. See skill-security-framing.md for details
```

**Step 3: Verify file updated**

Run: `wc -l .claude-plugin/.claude/skills/skill-deep-research.md`
Expected: ~120-150 lines (not thin wrapper ~10 lines)

**Step 4: Commit**

```bash
git add .claude-plugin/.claude/skills/skill-deep-research.md
git commit -m "feat(skills): restore skill-deep-research with reasoning content

- Replace thin wrapper with full execution contract
- Add interactive questions with AskUserQuestion
- Use atomic mp commands for state/route/validate
- Include visual indicators and cost awareness"
```

---

## P1: Double Diamond Phases

### Task 3: Restore flow-define.md

**Files:**
- Modify: `.claude-plugin/.claude/skills/flow-define.md`

**Step 1: Read current file**

Run: `cat .claude-plugin/.claude/skills/flow-define.md`

**Step 2: Replace with full reasoning content**

Key structure:
- Trigger for scoping/clarification requests
- 6-step execution contract
- Atomic mp commands: `mp validate`, `mp route --intent define`, `mp state`
- Phase dependency check (discover should be complete)
- Visual indicators for Define phase (🎯)

**Step 3: Verify file updated**

Run: `wc -l .claude-plugin/.claude/skills/flow-define.md`
Expected: ~150-180 lines

**Step 4: Commit**

```bash
git add .claude-plugin/.claude/skills/flow-define.md
git commit -m "feat(skills): restore flow-define with reasoning content

- Replace thin wrapper with full execution contract
- Add phase dependency check for discover
- Use atomic mp commands for state/route/validate
- Include visual indicators (🎯) and consensus building"
```

---

### Task 4: Restore flow-develop.md

**Files:**
- Modify: `.claude-plugin/.claude/skills/flow-develop.md`

**Step 1: Read current file**

Run: `cat .claude-plugin/.claude/skills/flow-develop.md`

**Step 2: Replace with full reasoning content**

Key structure:
- Trigger for build/implement requests
- 6-step execution contract
- Atomic mp commands: `mp validate`, `mp route --intent develop`, `mp state`
- Phase dependency check (define should be complete)
- Visual indicators for Develop phase (🛠️)
- Quality gate integration

**Step 3: Verify file updated**

Run: `wc -l .claude-plugin/.claude/skills/flow-develop.md`
Expected: ~150-180 lines

**Step 4: Commit**

```bash
git add .claude-plugin/.claude/skills/flow-develop.md
git commit -m "feat(skills): restore flow-develop with reasoning content

- Replace thin wrapper with full execution contract
- Add phase dependency check for define
- Use atomic mp commands for state/route/validate
- Include visual indicators (🛠️) and quality gates"
```

---

### Task 5: Restore flow-deliver.md

**Files:**
- Modify: `.claude-plugin/.claude/skills/flow-deliver.md`

**Step 1: Read current file**

Run: `cat .claude-plugin/.claude/skills/flow-deliver.md`

**Step 2: Replace with full reasoning content**

Key structure:
- Trigger for review/validate requests
- 6-step execution contract
- Atomic mp commands: `mp validate`, `mp route --intent deliver`, `mp state`
- Phase dependency check (develop should be complete)
- Visual indicators for Deliver phase (✅)
- Validation gate integration

**Step 3: Verify file updated**

Run: `wc -l .claude-plugin/.claude/skills/flow-deliver.md`
Expected: ~150-180 lines

**Step 4: Commit**

```bash
git add .claude-plugin/.claude/skills/flow-deliver.md
git commit -m "feat(skills): restore flow-deliver with reasoning content

- Replace thin wrapper with full execution contract
- Add phase dependency check for develop
- Use atomic mp commands for state/route/validate
- Include visual indicators (✅) and validation gates"
```

---

## P2: Extended Workflows

### Task 6: Create flow-parallel.md (New Skill)

**Files:**
- Create: `.claude-plugin/.claude/skills/flow-parallel.md`

**Step 1: Create new skill file**

Key structure:
- Trigger for parallel/decomposition requests
- 7-step execution contract for Team of Teams
- Uses `mp state` for work package tracking
- No shell script dependencies
- Each work package runs with full plugin capabilities

```markdown
---
name: flow-parallel
aliases:
  - parallel
  - team
  - team-of-teams
description: Team of Teams — decompose compound tasks across independent claude instances
execution_mode: enforced
validation_gates:
  - wbs_generated
  - instructions_written
  - all_work_packages_complete
---

## EXECUTION CONTRACT (MANDATORY)

### STEP 1: Clarifying Questions (MANDATORY)

Ask via AskUserQuestion:
1. **Compound task**: What should be decomposed?
2. **Work package count**: How many packages? (3-5 recommended)
3. **Dependencies**: Are packages independent?

---

### STEP 2: Visual Indicators (MANDATORY)

```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Team of Teams Mode
Parallel Phase: Decomposing into N independent work packages

Architecture:
  Main (this session) - Orchestrator
  WP-1..WP-N (claude -p) - Independent workers with full plugin
```

---

### STEP 3: Generate Work Breakdown Structure (WBS)

**Execute via Bash:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"parallel","wbs":["wp1","wp2","wp3"]}' --dir "$PWD" --json
```

Decompose the compound task into N independent work packages.

---

### STEP 4: Write Work Package Instructions

For each work package, create detailed instructions:
- Context from original task
- Specific deliverables
- Success criteria
- Dependencies (if any)

---

### STEP 5: Launch Independent Processes

For each work package, spawn independent `claude -p` process:
- Each gets full plugin capabilities
- Own context, tools, quality gates
- Produces output.md + exit-code

---

### STEP 6: Monitor & Aggregate

Track work package completion:
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state get --key "parallel_status" --dir "$PWD" --json
```

---

### STEP 7: Present Aggregated Results

Synthesize all work package outputs into coherent deliverable.
```

**Step 2: Verify file created**

Run: `test -f .claude-plugin/.claude/skills/flow-parallel.md && echo "OK"`

**Step 3: Commit**

```bash
git add .claude-plugin/.claude/skills/flow-parallel.md
git commit -m "feat(skills): add flow-parallel for Team of Teams

- New skill for parallel work decomposition
- Uses atomic mp commands for state tracking
- Independent claude -p processes with full plugin
- No shell script dependencies"
```

---

### Task 7: Create flow-spec.md (New Skill)

**Files:**
- Create: `.claude-plugin/.claude/skills/flow-spec.md`

**Step 1: Create new skill file**

Key structure:
- Trigger for specification/NLSpec requests
- 8-step execution contract
- Uses `mp state` for spec tracking
- Multi-AI specification generation
- Completeness validation

```markdown
---
name: flow-spec
aliases:
  - spec
  - nlspec
  - specification
description: NLSpec authoring — structured specification from multi-AI research
execution_mode: enforced
validation_gates:
  - mp_command_executed
  - spec_completeness_validated
---

## EXECUTION CONTRACT (MANDATORY)

### STEP 1: Clarifying Questions (MANDATORY)

Ask via AskUserQuestion:
1. **What to specify**: Project/feature name + description
2. **Actors**: Who interacts with this system?
3. **Key constraints**: Performance, Security, Compatibility, Scale
4. **Complexity class**: Clear, Complicated, or Complex

---

### STEP 2: Visual Indicators (MANDATORY)

```
🐙 **CLAUDE OCTOPUS ACTIVATED** - NLSpec Authoring Mode
Spec Phase: Generating specification for [project]

Providers:
🔴 Codex CLI - Structure and patterns
🟡 Gemini CLI - Ecosystem context
🔵 Claude - Synthesis & NLSpec generation
```

---

### STEP 3: Validate & Route (MANDATORY)

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent define --dir "$PWD" --json
```

---

### STEP 4: Multi-AI Specification Research

Query each provider for:
- Architectural patterns
- Similar systems in ecosystem
- Best practices

---

### STEP 5: Generate NLSpec Structure

Create structured specification:
- Overview
- Actors and interactions
- Constraints
- Success criteria

---

### STEP 6: Validate Completeness

Check all required sections present:
- [ ] Overview complete
- [ ] Actors identified
- [ ] Constraints documented
- [ ] Success criteria defined

---

### STEP 7: Persist Specification

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"spec","status":"complete","spec_path":"<path>"}' --dir "$PWD" --json
```

---

### STEP 8: Present Results

Display generated NLSpec with attribution.
```

**Step 2: Verify file created**

Run: `test -f .claude-plugin/.claude/skills/flow-spec.md && echo "OK"`

**Step 3: Commit**

```bash
git add .claude-plugin/.claude/skills/flow-spec.md
git commit -m "feat(skills): add flow-spec for NLSpec authoring

- New skill for structured specification generation
- Uses atomic mp commands for state/route/validate
- Multi-AI research for specification content
- Completeness validation gates"
```

---

## Final Verification

### Task 8: Verify All Skills Restored

**Files:**
- All 7 skill files

**Step 1: Run no-shell validation**

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type no-shell --dir "$PWD" --json
```

Expected: `status=ok`, no shell references found

**Step 2: Check line counts**

```bash
for f in flow-discover flow-define flow-develop flow-deliver skill-deep-research flow-parallel flow-spec; do
  echo "$f: $(wc -l < .claude-plugin/.claude/skills/$f.md) lines"
done
```

Expected: All files >100 lines (not thin wrappers)

**Step 3: Verify mp commands only**

```bash
grep -r "orchestrate.sh\|state-manager.sh\|octo-state.sh" .claude-plugin/.claude/skills/ || echo "OK - no shell scripts"
```

Expected: "OK - no shell scripts"

**Step 4: Final commit**

```bash
git add -A
git commit -m "feat(skills): complete hybrid architecture restoration

P0 (Critical):
- flow-discover: Full reasoning with atomic mp commands
- skill-deep-research: Interactive questions + mp commands

P1 (Double Diamond):
- flow-define: Phase dependency + consensus building
- flow-develop: Quality gates + implementation
- flow-deliver: Validation + review

P2 (Extended):
- flow-parallel: Team of Teams decomposition
- flow-spec: NLSpec authoring

All skills now use atomic mp commands instead of shell scripts.
Go runtime provides deterministic capabilities.
Markdown skills provide reasoning and orchestration."
```

---

## Acceptance Criteria

- [ ] All 7 skills are no longer thin wrappers
- [ ] No references to `orchestrate.sh`, `state-manager.sh`, `octo-state.sh`
- [ ] All skills use atomic `mp` commands (`state`, `validate`, `route`)
- [ ] JSON response parsing for branching logic
- [ ] Visual indicators displayed before execution
- [ ] Error handling with `status` checks
- [ ] `mp validate --type no-shell` passes

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Behavior drift from main | Test each skill manually after restoration |
| Missing atomic commands | Verify `mp` binary has required subcommands |
| JSON parsing errors | Use `jq` for robust parsing in bash |

---

*Plan generated: 2026-03-02*
