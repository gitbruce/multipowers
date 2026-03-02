# flow-discover

Discover phase for research and exploration using multi-AI perspectives.

## Overview

This skill orchestrates the discovery phase of the workflow, using atomic mp commands
to gather information and structure findings. It leverages multiple AI providers for
broad exploration and synthesis.

## Workflow Stages

### Stage 1: Validate Workspace

**Goal:** Ensure the workspace is ready for discovery work.

**Action:** Run workspace validation
```bash
mp validate --type workspace --dir . --json
```

**Branch:**
- If `status=ok`: Proceed to Stage 2
- If `status=error`: Check `message` for reason, guide user to run `/mp:init` if context missing
- If `status=blocked`: Review `data.missing` and collect required context

### Stage 2: Route Providers

**Goal:** Determine which AI providers to use for discovery.

**Action:** Run provider routing
```bash
mp route --intent discover --dir . --json
```

**Branch:**
- If `status=ok`: Review `data.selected_providers` and `data.reason`
- If `status=error`: Check `message` - may need to configure providers
- Use `data.available_providers` to understand what's available

### Stage 3: Execute Discovery

**Goal:** Perform multi-perspective research and exploration.

**Action:** Based on routing results, coordinate discovery across providers:
- Use `data.selected_providers` to determine which AIs to query
- Synthesize findings from multiple perspectives
- Structure results in a coherent format

**State Tracking:**
```bash
mp state update --data '{"phase":"discover","stage":"executing"}' --dir . --json
```

### Stage 4: Persist Findings

**Goal:** Save discovery results for downstream phases.

**Action:** Update state with discovery results
```bash
mp state update --data '{"phase":"discover","status":"complete","findings":"<summary>"}' --dir . --json
```

## Response Contract

All mp commands return a JSON response with:
- `status`: "ok" | "error" | "blocked"
- `action`: Recommended next action
- `message`: Human-readable status
- `error_code`: Error identifier if applicable
- `data`: Structured response data
- `remediation`: Suggested fix if blocked

## Error Handling

| Condition | Action |
|-----------|--------|
| Workspace invalid | Guide user to `/mp:init` |
| No providers available | Check provider configuration |
| Provider routing fails | Fall back to single-provider mode |
| State update fails | Retry with simpler data |

## Example Usage

User: "/mp:discover OAuth authentication patterns"

1. Validate workspace readiness
2. Route to available providers (e.g., claude, codex, gemini)
3. Execute parallel research queries
4. Synthesize findings into structured output
5. Update state with discovery completion
