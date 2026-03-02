# flow-deliver

Deliver phase for validation, review, and deployment preparation.

## Overview

This skill orchestrates the delivery phase, using atomic mp commands to
validate implementation, perform code review, and prepare for deployment.

## Workflow Stages

### Stage 1: Load Development Context

**Goal:** Retrieve previous development phase results.

**Action:** Check current state
```bash
mp state get --dir . --json
```

**Branch:**
- If `data.state.phase=develop` and `data.state.status=complete`: Proceed to Stage 2
- If no development context: Guide user to run `/mp:develop` first
- Use `data.state.tests_passed` and `data.state.coverage` to validate readiness

### Stage 2: Validate Workspace

**Goal:** Ensure workspace is ready for delivery.

**Action:** Run workspace validation
```bash
mp validate --type workspace --dir . --json
```

**Branch:**
- If `status=ok`: Proceed to Stage 3
- If `status=error`: Check `message` for reason

### Stage 3: Route Providers

**Goal:** Determine which AI providers to use for delivery.

**Action:** Run provider routing
```bash
mp route --intent deliver --dir . --json
```

**Branch:**
- Review `data.selected_providers` and `data.reason`
- Delivery may use multi-provider mode for comprehensive review

### Stage 4: Run Tests

**Goal:** Verify all tests pass before delivery.

**Action:** Run test suite
```bash
mp test run --dir . --json
```

**Branch:**
- If `data.status=passed`: Proceed to coverage check
- If `data.status=failed`: Review `data.failed_tests` and fix before proceeding

### Stage 5: Check Coverage

**Goal:** Ensure adequate test coverage.

**Action:** Run coverage check
```bash
mp coverage check --dir . --json
```

**Branch:**
- If `data.coverage_pct` meets threshold: Proceed to no-shell validation
- If coverage low: Return to development phase

### Stage 6: Validate No-Shell Runtime

**Goal:** Ensure no shell script references in codebase.

**Action:** Run no-shell validation
```bash
mp validate --type no-shell --dir . --json
```

**Branch:**
- If `status=ok`: Proceed to review
- If `status=blocked`: Review `data.violations` and remove shell references

### Stage 7: Execute Review

**Goal:** Perform comprehensive code review.

**Action:** Based on routing results:
- Use `data.selected_providers` to coordinate review
- Check for security issues, code quality, documentation
- Validate against acceptance criteria from definition phase

**State Tracking:**
```bash
mp state update --data '{"phase":"deliver","stage":"reviewing"}' --dir . --json
```

### Stage 8: Persist Delivery

**Goal:** Save delivery results and mark workflow complete.

**Action:** Update state with delivery results
```bash
mp state update --data '{"phase":"deliver","status":"complete","review_passed":true,"ready_for_deploy":true}' --dir . --json
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
| No development context | Guide user to `/mp:develop` |
| Tests fail | Return to development phase |
| Coverage too low | Add tests in development phase |
| No-shell violations | Remove shell script references |
| Review finds issues | Fix issues and re-run delivery |

## Example Usage

User: "/mp:deliver authentication system"

1. Load development results
2. Validate workspace
3. Route to appropriate providers
4. Run tests (verify pass)
5. Check coverage
6. Validate no-shell runtime
7. Execute code review
8. Save delivery results
