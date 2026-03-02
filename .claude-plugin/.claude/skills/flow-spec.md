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
trigger: |
  AUTOMATICALLY ACTIVATE when user requests specification:
  - "create a specification for X"
  - "write NLSpec for Y"
  - "specify the requirements for Z"
  - "generate a spec document"
---

# flow-spec

Generate structured specifications (NLSpec) from multi-AI research.

## EXECUTION CONTRACT (MANDATORY - CANNOT SKIP)

### STEP 1: Clarifying Questions (MANDATORY)

**Ask via AskUserQuestion:**

```javascript
AskUserQuestion({
  questions: [
    {
      question: "What should be specified?",
      header: "Subject",
      multiSelect: false,
      options: [
        {label: "Use provided description", description: "Use the description from the command"},
        {label: "Describe now", description: "I'll provide the system description"}
      ]
    },
    {
      question: "Who interacts with this system?",
      header: "Actors",
      multiSelect: true,
      options: [
        {label: "End Users", description: "People using the system"},
        {label: "Developers", description: "Engineers building/maintaining"},
        {label: "Admins", description: "System administrators"},
        {label: "External Services", description: "Other systems/APIs"}
      ]
    },
    {
      question: "What constraints matter most?",
      header: "Constraints",
      multiSelect: true,
      options: [
        {label: "Performance", description: "Speed and responsiveness"},
        {label: "Security", description: "Data protection and access control"},
        {label: "Compatibility", description: "Integration requirements"},
        {label: "Scale", description: "Handling growth and volume"}
      ]
    },
    {
      question: "How complex is this system?",
      header: "Complexity",
      multiSelect: false,
      options: [
        {label: "Clear", description: "Well-understood, straightforward"},
        {label: "Complicated", description: "Multiple parts, but knowable"},
        {label: "Complex", description: "Emergent behavior, unknowns"}
      ]
    }
  ]
})
```

**DO NOT PROCEED TO STEP 2 until questions answered.**

---

### STEP 2: Visual Indicators (MANDATORY)

Display the specification banner:

```
🐙 **CLAUDE OCTOPUS ACTIVATED** - NLSpec Authoring Mode
Spec Phase: Generating specification for [project/feature]

Providers:
🔴 Codex CLI - Structure and patterns
🟡 Gemini CLI - Ecosystem context
🔵 Claude - Synthesis & NLSpec generation

Estimated Time: 3-7 minutes
```

---

### STEP 3: Validate & Route (MANDATORY)

**Validate workspace:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" validate --type workspace --dir "$PWD" --json
```

**Route providers:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" route --intent define --dir "$PWD" --json
```

**Update state:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"spec","stage":"researching"}' --dir "$PWD" --json
```

---

### STEP 4: Multi-AI Specification Research

**Query each provider:**

🔴 **Codex Analysis:**
- Architectural patterns for similar systems
- Technical best practices
- Implementation considerations

🟡 **Gemini Analysis:**
- Ecosystem standards
- Community approaches
- Integration patterns

🔵 **Claude Synthesis:**
- Combine perspectives
- Identify gaps
- Structure specification

---

### STEP 5: Generate NLSpec Structure

Create structured specification:

```markdown
# [System Name] Specification

## Overview
[2-3 sentence description of the system]

## Actors
| Actor | Description | Primary Interactions |
|-------|-------------|---------------------|
| [Actor 1] | [Who they are] | [What they do] |

## Functional Requirements
1. [Requirement 1]
2. [Requirement 2]
3. [Requirement 3]

## Non-Functional Requirements
| Category | Requirement |
|----------|-------------|
| Performance | [Performance requirement] |
| Security | [Security requirement] |
| Compatibility | [Compatibility requirement] |
| Scale | [Scale requirement] |

## Constraints
- [Constraint 1]
- [Constraint 2]

## Success Criteria
- [ ] [Criterion 1]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

## Out of Scope
- [Explicitly excluded items]

## Open Questions
- [Question 1]
- [Question 2]
```

---

### STEP 6: Validate Completeness (MANDATORY)

**Check all required sections present:**

```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"spec","stage":"validating"}' --dir "$PWD" --json
```

**Completeness checklist:**
- [ ] Overview complete
- [ ] Actors identified
- [ ] Functional requirements listed
- [ ] Non-functional requirements defined
- [ ] Constraints documented
- [ ] Success criteria defined
- [ ] Out of scope items listed

**If incomplete:**
- Note missing sections
- Generate placeholder content
- Flag for user review

---

### STEP 7: Persist Specification (MANDATORY)

**Save specification:**
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" state update --data '{"phase":"spec","status":"complete","spec_path":"specs/[name]-spec.md","sections_complete":"<count>"}' --dir "$PWD" --json
```

**Write spec file:**
```bash
mkdir -p specs
cat > specs/[name]-spec.md << 'EOF'
[Specification content]
EOF
```

---

### STEP 8: Present Results

**Output format:**

```
## Specification Generated

**File:** specs/[name]-spec.md

### Summary
- **Actors:** [count] identified
- **Requirements:** [count] functional, [count] non-functional
- **Constraints:** [count] defined
- **Success Criteria:** [count] criteria

### Sections Complete
- [x] Overview
- [x] Actors
- [x] Functional Requirements
- [x] Non-Functional Requirements
- [x] Constraints
- [x] Success Criteria
- [x] Out of Scope
- [ ] Open Questions (if any)

### Next Steps
1. Review the specification
2. Answer open questions
3. Proceed to `/mp:define` for detailed requirements
```

---

## Error Handling

| Condition | Action |
|-----------|--------|
| Workspace invalid | Guide user to `/mp:init` |
| No providers | Proceed with Claude-only synthesis |
| Incomplete spec | Generate placeholders, flag for review |
| Validation fails | List missing sections, prompt user |

## NLSpec Format Reference

NLSpec (Natural Language Specification) is a structured format for capturing system requirements in natural language while maintaining clear organization.

**Key principles:**
- Human-readable first
- Structured sections
- Clear actors and interactions
- Measurable success criteria
- Explicit scope boundaries
