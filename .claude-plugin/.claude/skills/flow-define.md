---
name: flow-define
aliases:
  - define
  - define-workflow
  - grasp
description: Multi-AI requirements scoping (Double Diamond Define phase)
agent: Plan
context: fork
task_management: true
execution_mode: enforced
validation_gates:
  - mp_command_executed
  - synthesis_complete
trigger: |
  AUTOMATICALLY ACTIVATE when user requests clarification or scoping:
  - "define the requirements for X"
  - "clarify the scope of Y"
  - "what exactly does X need to do"
  - "scope out the Z feature"
---

# flow-define

Multi-AI requirements scoping using the Double Diamond Define phase.
Builds consensus on scope, constraints, and success criteria.

## Pre-Definition: Phase Dependency Check

Before starting definition, check if discover phase is complete:

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state get --key "phase" --dir "$PWD" --json
```

**Branch:**
- If `data.phase=discover` and `data.status=complete`: Proceed with full context
- If discover not complete: Warn user but allow proceeding
- Use `data.findings` to inform definition decisions

---

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

This workflow MUST execute all steps in order. Skipping steps is PROHIBITED.

### STEP 1: Validate Workspace (MANDATORY)

Execute workspace validation:

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json
```

**Parse JSON response:**
- If `status=ok`: Proceed to STEP 2
- If `status=error`: Guide user to `/mp:init` to initialize workspace
- Check `remediation` field for specific guidance

---

### STEP 2: Route Providers (MANDATORY)

Determine which AI providers participate in definition:

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent define --dir "$PWD" --json
```

**Display visual indicators immediately after routing:**

```
🐙 **CLAUDE OCTOPUS ACTIVATED** - Multi-provider definition mode
🎯 Define Phase: [Brief description of what is being defined]

Providers:
🔴 Codex CLI: Technical requirements analysis
🟡 Gemini CLI: Business context and constraints
🔵 Claude: Consensus building and synthesis
```

---

### STEP 3: Update State (MANDATORY)

Mark definition phase as executing:

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"define","stage":"executing"}' --dir "$PWD" --json
```

---

### STEP 4: Execute Definition (MANDATORY - Consensus Building)

Gather perspectives from selected providers and build consensus:

**4a. Gather Multi-AI Perspectives:**

🔴 **Codex Analysis:**
- Technical requirements and constraints
- Implementation complexity assessment
- Technology-specific considerations
- Integration points and dependencies

🟡 **Gemini Analysis:**
- Business context and stakeholder needs
- Market constraints and opportunities
- User experience implications
- Risk factors and mitigation strategies

**4b. Build Consensus On:**

| Aspect | Questions to Answer |
|--------|---------------------|
| Scope Boundaries | What is IN scope? What is OUT of scope? |
| Success Criteria | How do we measure completion? |
| Constraints | What limitations must we work within? |
| Trade-offs | What compromises are acceptable? |
| Priority Ordering | What must be done first? |

**4c. Document Decisions:**
- Record each decision with rationale
- Note any disagreements and resolution approach
- Capture assumptions made during definition

---

### STEP 5: Persist Definition (MANDATORY)

Save definition results to state:

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"define","status":"complete","scope":"<scope_summary>","decisions":"<key_decisions>","criteria":"<success_criteria>"}' --dir "$PWD" --json
```

**Required fields in data:**
- `scope`: Brief scope summary
- `decisions`: Key decisions made
- `criteria`: Success criteria defined

---

### STEP 6: Present Results

Output the definition summary:

```
## Definition Complete

### Scope Summary
[Concise description of defined scope]

### Key Requirements
1. [Requirement 1]
2. [Requirement 2]
3. [Requirement 3]

### Success Criteria
- [ ] [Criterion 1]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

### Constraints
- [Constraint 1]
- [Constraint 2]

### Next Steps
Proceed to **Develop phase** with: `/mp:develop` or `octo build`
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Workspace invalid | Guide user to `/mp:init` to initialize |
| No providers available | Check `mp status` for provider health |
| Discover not complete | Warn user, offer to proceed or run discover first |
| Consensus fails | Document disagreements, escalate to user for decision |
| State update fails | Retry with simplified data structure |
| Validation blocked | Check `remediation` field, follow guidance |

---

## Response Contract

All mp commands return JSON with:

| Field | Description |
|-------|-------------|
| `status` | "ok" \| "error" \| "blocked" |
| `action` | Recommended next action |
| `message` | Human-readable status |
| `error_code` | Error identifier if applicable |
| `data` | Structured response data |
| `remediation` | Suggested fix if blocked |

---

## Phase Transitions

| Direction | Trigger |
|-----------|---------|
| **From Discover** | Research complete, ready to define scope |
| **To Develop** | Requirements defined, ready to implement |

**Transition command to Develop:**
```
/mp:develop
```

---

## Example Usage

**User:** "/mp:define authentication system requirements"

**Execution:**
1. Load discovery context (if exists)
2. Validate workspace readiness
3. Route to appropriate providers
4. Gather technical and business perspectives
5. Build consensus on scope and criteria
6. Persist definition to state
7. Present summary and next steps
