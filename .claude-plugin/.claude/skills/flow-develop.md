# flow-develop

Develop phase for implementation with quality gates.

## Overview

This skill orchestrates the development phase, using atomic mp commands to
implement solutions with TDD discipline and quality gates.

## Workflow Stages

### Stage 1: Load Definition Context

**Goal:** Retrieve previous definition phase results.

**Action:** Check current state
```bash
mp state get --dir . --json
```

**Branch:**
- If `data.state.phase=define` and `data.state.status=complete`: Proceed to Stage 2
- If no definition context: Guide user to run `/mp:define` first
- Use `data.state.spec` to guide implementation

### Stage 2: Validate TDD Environment

**Goal:** Ensure TDD environment is ready.

**Action:** Run TDD environment validation
```bash
mp validate --type tdd-env --dir . --json
```

**Branch:**
- If `status=ok`: Proceed to Stage 3
- If `status=error`: Check `data.details` for missing requirements

### Stage 3: Route Providers

**Goal:** Determine which AI providers to use for development.

**Action:** Run provider routing
```bash
mp route --intent develop --dir . --json
```

**Branch:**
- Review `data.selected_providers` and `data.reason`
- Development typically uses single-provider mode for code consistency

### Stage 4: Run Tests (Red Phase)

**Goal:** Verify test status before implementation.

**Action:** Run test suite
```bash
mp test run --dir . --json
```

**Branch:**
- If `data.status=failed`: Good - tests are red, proceed to implement
- If `data.status=passed`: Tests exist - ensure testing the right thing
- Review `data.failed_tests` for expected failures

### Stage 5: Execute Development

**Goal:** Implement solution to make tests pass.

**Action:** Based on routing results:
- Use `data.selected_providers` to coordinate implementation
- Write minimal code to pass tests
- Follow coding standards and patterns

**State Tracking:**
```bash
mp state update --data '{"phase":"develop","stage":"implementing"}' --dir . --json
```

### Stage 6: Run Tests (Green Phase)

**Goal:** Verify implementation passes tests.

**Action:** Run test suite
```bash
mp test run --dir . --json
```

**Branch:**
- If `data.status=passed`: Proceed to coverage check
- If `data.status=failed`: Review failures and iterate

### Stage 7: Check Coverage

**Goal:** Ensure adequate test coverage.

**Action:** Run coverage check
```bash
mp coverage check --dir . --json
```

**Branch:**
- If `data.coverage_pct` meets threshold: Proceed to persist
- If coverage low: Add more tests and iterate

### Stage 8: Persist Development

**Goal:** Save development results for delivery phase.

**Action:** Update state with development results
```bash
mp state update --data '{"phase":"develop","status":"complete","tests_passed":true,"coverage":"<pct>"}' --dir . --json
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
| No definition context | Guide user to `/mp:define` |
| TDD environment invalid | Check test framework setup |
| Tests fail unexpectedly | Debug and fix implementation |
| Coverage too low | Add additional tests |

## Example Usage

User: "/mp:develop authentication system"

1. Load definition specification
2. Validate TDD environment
3. Route to appropriate provider
4. Run tests (verify red)
5. Implement solution
6. Run tests (verify green)
7. Check coverage
8. Save development results
