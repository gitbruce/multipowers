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
