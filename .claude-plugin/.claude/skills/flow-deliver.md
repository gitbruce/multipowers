---
name: flow-deliver
aliases:
  - deliver
  - deliver-workflow
  - ink
description: Multi-AI validation and review (Double Diamond Deliver phase)
agent: general-purpose
context: fork
task_management: true
execution_mode: enforced
validation_gates:
  - mp_command_executed
  - validation_complete
trigger: |
  AUTOMATICALLY ACTIVATE when user requests validation or review:
  - "review X" or "validate Y"
  - "check if X works correctly"
  - "quality check for X"
---

# flow-deliver

Deliver phase for validation, review, and deployment preparation.

## Pre-Delivery: Phase Dependency Check

Before starting delivery, check if develop phase is complete:

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state get --key "phase" --dir "$PWD" --json
```

**Branch:**
- If `data.phase=develop` and `data.status=complete`: Proceed with full context
- If develop not complete: Warn user but allow proceeding

---

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

### STEP 1: Validate Workspace (MANDATORY)

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json
```

**Validation Gate:** Check `status=ok` before proceeding.

---

### STEP 2: Route Providers (MANDATORY)

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent deliver --dir "$PWD" --json
```

**Display visual indicators:**
```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Multi-provider validation mode
✅ Deliver Phase: [Brief description of what's being validated]

Providers:
🔴 Codex CLI: Code quality analysis
🟡 Gemini CLI: Security and edge cases
🔵 Claude: Synthesis and recommendations
```

**Validation Gate:** Verify `data.selected_providers` is populated.

---

### STEP 3: Update State (MANDATORY)

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"deliver","stage":"executing"}' --dir "$PWD" --json
```

---

### STEP 4: Execute Validation (MANDATORY - Multi-Perspective Review)

**4a. Code Quality Review (🔴 Codex):**
- Code structure and patterns
- Test coverage analysis
- Performance considerations

**4b. Security Review (🟡 Gemini):**
- Security vulnerabilities
- Edge case handling
- Error handling completeness

**4c. Synthesis (🔵 Claude):**
- Integrate findings from all providers
- Prioritize issues by severity
- Generate actionable recommendations

**Format output with provider indicators:**
```
🔴 **Codex Analysis:**
[Codex code quality findings...]

🟡 **Gemini Analysis:**
[Gemini security and edge case findings...]

🔵 **Claude Synthesis:**
[Integrated recommendations and priority actions...]
```

---

### STEP 5: Persist Validation (MANDATORY)

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"deliver","status":"complete","validation_result":"<pass/fail>","issues_found":"<count>"}' --dir "$PWD" --json
```

---

### STEP 6: Present Results

**Output format:**
```
## Validation Summary

**Result:** PASS / FAIL
**Issues Found:** [count] ([critical] critical, [high] high, [medium] medium)

### Recommendations
1. [Priority recommendation]
2. [Secondary recommendation]

### Next Steps
- [ ] Address critical issues before deployment
- [ ] Review high-priority items
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Workspace invalid | Guide user to `/mp:init` |
| No providers available | Check `mp status` |
| Develop not complete | Warn but allow proceeding |
| Validation fails | Document issues, suggest fixes |
| Critical issues found | Block deployment until resolved |

## Phase Transitions

- **From:** Develop (implementation complete)
- **To:** Ship (ready for deployment)
