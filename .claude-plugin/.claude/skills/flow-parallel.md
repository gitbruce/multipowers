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
trigger: |
  AUTOMATICALLY ACTIVATE when user requests parallel work decomposition:
  - "decompose X into parallel tasks"
  - "run team of teams for Y"
  - "split Z into independent work packages"
  - "parallelize this compound task"
---

# flow-parallel

Decompose compound tasks into independent work packages that can be executed in parallel by separate claude instances.

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

### STEP 1: Clarifying Questions (MANDATORY)

**Ask via AskUserQuestion:**

```javascript
AskUserQuestion({
  questions: [
    {
      question: "What compound task should be decomposed?",
      header: "Task",
      multiSelect: false,
      options: [
        {label: "Use provided description", description: "Use the task description from the command"},
        {label: "Describe now", description: "I'll provide the task description"}
      ]
    },
    {
      question: "How many work packages?",
      header: "Count",
      multiSelect: false,
      options: [
        {label: "3 (Recommended)", description: "Optimal for most tasks"},
        {label: "4", description: "For larger compound tasks"},
        {label: "5", description: "For complex multi-part tasks"},
        {label: "Custom", description: "Specify up to 10"}
      ]
    },
    {
      question: "Are the work packages independent?",
      header: "Dependencies",
      multiSelect: false,
      options: [
        {label: "Fully independent (Recommended)", description: "No dependencies between packages"},
        {label: "Some dependencies", description: "Packages may share interfaces"},
        {label: "Sequential dependencies", description: "Packages must complete in order"}
      ]
    }
  ]
})
```

**DO NOT PROCEED TO STEP 2 until questions answered.**

---

### STEP 2: Visual Indicators (MANDATORY)

Display the orchestration banner:

```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Team of Teams Mode
Parallel Phase: Decomposing into N independent work packages

Architecture:
  Main (this session) - Orchestrator: decompose, launch, monitor, aggregate
  WP-1..WP-N (claude -p) - Independent workers with full plugin capabilities

Each worker:
  - Runs as independent claude -p process
  - Loads full Octopus plugin
  - Has own context, tools, and quality gates
  - Produces output.md + exit-code

Estimated Time: 5-15 minutes (depending on task complexity)
```

---

### STEP 3: Generate Work Breakdown Structure (WBS)

**Initialize parallel state:**

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"parallel","stage":"decomposing","work_package_count":"<N>"}' --dir "$PWD" --json
```

**Decompose the compound task:**

1. Analyze the compound task
2. Identify natural boundaries for decomposition
3. Create N independent work packages
4. Document dependencies (if any)

**WBS format:**
| WP # | Description | Deliverable | Dependencies |
|------|-------------|-------------|--------------|
| WP-1 | [Description] | [Expected output] | None |
| WP-2 | [Description] | [Expected output] | None |
| ... | ... | ... | ... |

---

### STEP 4: Write Work Package Instructions

For each work package, create detailed instructions:

**Instruction template:**
```markdown
## Work Package N: [Title]

### Context
[Relevant context from original task]

### Objective
[Specific goal for this work package]

### Deliverables
- [ ] [Deliverable 1]
- [ ] [Deliverable 2]

### Success Criteria
- [Criterion 1]
- [Criterion 2]

### Dependencies
[List any dependencies on other work packages]

### Commands
Run: `claude -p "[full instruction]"`
```

**Update state:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"parallel","stage":"instructions_written","wps":["wp1","wp2","wp3"]}' --dir "$PWD" --json
```

---

### STEP 5: Launch Independent Processes

For each work package, spawn independent process:

**Process launch pattern:**
```bash
claude -p "[Work package instruction]" > outputs/wp-N-output.md 2>&1 &
```

**Tracking:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"parallel","stage":"executing","launched":["wp1","wp2","wp3"]}' --dir "$PWD" --json
```

**Each process:**
- Gets full plugin capabilities
- Own context, tools, quality gates
- Produces output.md + exit-code

---

### STEP 6: Monitor & Aggregate

**Monitor completion:**

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state get --key "parallel_status" --dir "$PWD" --json
```

**Track work package status:**
| WP # | Status | Exit Code | Output File |
|------|--------|-----------|-------------|
| WP-1 | Complete | 0 | outputs/wp-1-output.md |
| WP-2 | Complete | 0 | outputs/wp-2-output.md |
| ... | ... | ... | ... |

**Aggregate results:**
1. Collect all work package outputs
2. Check for failures (non-zero exit codes)
3. Merge deliverables
4. Resolve any conflicts

---

### STEP 7: Present Aggregated Results

**Output format:**

```
## Parallel Execution Complete

### Work Package Summary
| WP | Status | Key Deliverables |
|----|--------|-----------------|
| WP-1 | ✅ Complete | [Deliverables] |
| WP-2 | ✅ Complete | [Deliverables] |
| WP-3 | ✅ Complete | [Deliverables] |

### Aggregated Deliverables
[Combined outputs from all work packages]

### Issues Found
[Any issues or conflicts discovered]

### Next Steps
[Recommended follow-up actions]
```

**Update final state:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"parallel","status":"complete","wp_results":"<summary>"}' --dir "$PWD" --json
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Decomposition fails | Task may not be suitable for parallelization |
| Work package fails | Retry individually or report issue |
| Conflicts in aggregation | Manual resolution required |
| Timeout | Check individual process logs |

## Architecture Notes

**Why independent processes?**
- Task tool subagents do NOT load plugins
- Independent `claude -p` processes DO load plugins
- Each work package gets full Octopus capabilities
- Includes Double Diamond workflows, agents, quality gates
