# skill-tdd

Test-Driven Development skill with red-green-refactor discipline.

## Overview

This skill provides TDD workflow guidance using atomic mp commands to
enforce test-first development with quality gates.

## TDD Cycle

### Red Phase: Write Failing Test

**Goal:** Write a test that fails because the feature doesn't exist.

**Action:** Validate TDD environment first
```bash
mp validate --type tdd-env --dir . --json
```

**Branch:**
- If `status=ok`: Proceed to write failing test
- If `status=error`: Check `data.details` for missing requirements

**Write Test:**
- Write a minimal test that captures the requirement
- The test should fail for the right reason (feature not implemented)

**Verify Red:**
```bash
mp test run --dir . --json
```

**Expected:** `data.status=failed` with the new test in `data.failed_tests`

### Green Phase: Make Test Pass

**Goal:** Write minimal code to make the test pass.

**Action:** Implement the feature
- Write the simplest code that makes the test pass
- Don't worry about perfection - just make it work
- Avoid over-engineering

**Verify Green:**
```bash
mp test run --dir . --json
```

**Expected:** `data.status=passed` with all tests passing

### Refactor Phase: Improve Code Quality

**Goal:** Clean up the code while keeping tests green.

**Action:** Refactor implementation
- Remove duplication
- Improve naming
- Simplify logic
- Add documentation

**Verify Still Green:**
```bash
mp test run --dir . --json
```

**Expected:** `data.status=passed` - tests remain green after refactoring

### Coverage Check

**Goal:** Ensure adequate test coverage.

**Action:** Run coverage check
```bash
mp coverage check --dir . --json
```

**Branch:**
- If `data.coverage_pct` meets threshold: Cycle complete
- If coverage low: Add more tests before proceeding

## State Tracking

Track TDD progress in state:
```bash
mp state update --data '{"tdd_cycle":"red","current_feature":"<feature>"}' --dir . --json
```

Update as you progress:
```bash
mp state update --data '{"tdd_cycle":"green","current_feature":"<feature>"}' --dir . --json
mp state update --data '{"tdd_cycle":"refactor","current_feature":"<feature>"}' --dir . --json
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
| TDD environment invalid | Set up test framework first |
| Test passes unexpectedly | Ensure testing the right thing |
| Multiple tests fail | Fix one at a time |
| Refactor breaks tests | Revert and refactor more carefully |

## Example Usage

User: "/mp:tdd implement user registration"

1. Validate TDD environment
2. Write failing test for user registration
3. Verify test fails (red)
4. Implement minimal registration logic
5. Verify test passes (green)
6. Refactor for code quality
7. Verify tests still pass
8. Check coverage
9. Repeat for next feature
