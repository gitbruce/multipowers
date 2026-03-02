---
name: flow-develop
aliases:
  - develop
  - develop-workflow
  - tangle
description: Multi-AI implementation (Double Diamond Develop phase)
agent: general-purpose
context: fork
task_management: true
execution_mode: enforced
validation_gates:
  - mp_command_executed
  - quality_gates_passed
trigger: |
  AUTOMATICALLY ACTIVATE when user requests building or implementation:
  - "build X" or "implement Y" or "create Z"
  - "develop a feature for X"
  - "generate implementation for X"
---

# Flow Develop - Double Diamond Develop Phase

Multi-provider implementation workflow with TDD discipline and quality gates.

## Pre-Development: Phase Dependency Check

Before starting development:
1. Check if define phase is complete via `mp state get`
2. If not complete, warn user but allow proceeding

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state get --key "phase" --dir "$PWD" --json
```

**Branch logic:**
- If `data.phase=define` and `data.status=complete`: Proceed with development
- If define not complete: Warn user, offer to proceed or run `/mp:define` first

---

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

You MUST execute these steps in order. Skipping steps is PROHIBITED.

### STEP 1: Validate Workspace (MANDATORY)

**Goal:** Ensure workspace is ready for development.

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json
```

**Response handling:**
- If `status=ok`: Proceed to STEP 2
- If `status=error`: Check `data.details`, guide user to fix issues
- If `status=blocked`: Follow `remediation` guidance

---

### STEP 2: Route Providers (MANDATORY)

**Goal:** Determine which AI providers to use for development.

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent develop --dir "$PWD" --json
```

**Display visual indicators:**
```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Multi-provider implementation mode
🛠️ Develop Phase: [Brief description of what's being built]

Providers:
🔴 Codex CLI: Code generation and patterns
🟡 Gemini CLI: Alternative approaches
🔵 Claude: Integration and quality gates
```

**Response handling:**
- Review `data.selected_providers` and `data.reason`
- Development typically uses coordinated multi-provider mode

---

### STEP 3: Update State (MANDATORY)

**Goal:** Record development phase start.

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"develop","stage":"executing"}' --dir "$PWD" --json
```

---

### STEP 4: Execute Development (MANDATORY - Quality Gates)

**4a. Implementation:**
- Use selected providers for code generation
- Follow TDD principles (red-green-refactor)
- Apply patterns from discovery/definition phases
- Write minimal code to pass tests

**4b. Quality Gates:**

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" test run --dir "$PWD" --json
"${CLAUDE_PLUGIN_ROOT}/bin/mp" coverage check --dir "$PWD" --json
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type no-shell --dir "$PWD" --json
```

**4c. Integration:**
- Merge perspectives from all providers
- Resolve conflicts between generated code
- Ensure consistency with project patterns

---

### STEP 5: Persist Development (MANDATORY)

**Goal:** Save development results for delivery phase.

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"develop","status":"complete","files_modified":"<list>","tests_passed":"<count>"}' --dir "$PWD" --json
```

---

### STEP 6: Present Results

**Output format:**

```
🛠️ **Development Complete**

**Implementation Summary:**
- [Brief description of what was built]

**Files Modified:**
- [List of files created/modified]

**Quality Gates:**
- Tests: [X] passed / [Y] total
- Coverage: [Z]%
- No-shell validation: [passed/failed]

**Next Steps:**
- Run `/mp:deliver` to validate and document
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Workspace invalid | Guide to `/mp:init` to initialize |
| Define not complete | Warn user, offer to proceed or run define |
| Tests failing | Block progression, fix before proceeding |
| Coverage too low | Block progression, add more tests |
| No-shell validation failed | Remove shell script references |
| No providers available | Check `mp status` for provider health |

---

## Phase Transitions

- **From:** Define (requirements and scope defined)
- **To:** Deliver (implementation complete, ready for validation)
