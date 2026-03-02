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
