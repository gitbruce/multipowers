# flow-define

Define phase for requirements clarification and scope definition.

## Overview

This skill orchestrates the definition phase, using atomic mp commands to
clarify requirements, establish scope, and create structured specifications.

## Workflow Stages

### Stage 1: Load Discovery Context

**Goal:** Retrieve previous discovery phase results.

**Action:** Check current state
```bash
mp state get --dir . --json
```

**Branch:**
- If `data.state.phase=discover` and `data.state.status=complete`: Proceed to Stage 2
- If no discovery context: Optionally run discovery first or proceed with user input
- Use `data.state.findings` to inform definition

### Stage 2: Validate Workspace

**Goal:** Ensure workspace is ready for definition work.

**Action:** Run workspace validation
```bash
mp validate --type workspace --dir . --json
```

**Branch:**
- If `status=ok`: Proceed to Stage 3
- If `status=error`: Check `message` for reason

### Stage 3: Route Providers

**Goal:** Determine which AI providers to use for definition.

**Action:** Run provider routing
```bash
mp route --intent define --dir . --json
```

**Branch:**
- Review `data.selected_providers` and `data.reason`
- Definition typically uses single-provider mode for consistency

### Stage 4: Execute Definition

**Goal:** Create structured requirements and scope.

**Action:** Based on routing results:
- Use `data.selected_providers` to coordinate definition work
- Generate structured specification document
- Define acceptance criteria and constraints

**State Tracking:**
```bash
mp state update --data '{"phase":"define","stage":"executing"}' --dir . --json
```

### Stage 5: Persist Definition

**Goal:** Save definition results for development phase.

**Action:** Update state with definition results
```bash
mp state update --data '{"phase":"define","status":"complete","spec":"<path>"}' --dir . --json
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
| No discovery context | Ask user to proceed without or run discovery |
| Workspace invalid | Guide user to `/mp:init` |
| State update fails | Retry with simpler data |

## Example Usage

User: "/mp:define authentication system requirements"

1. Load any existing discovery context
2. Validate workspace readiness
3. Route to appropriate provider
4. Generate structured requirements
5. Save definition results to state
