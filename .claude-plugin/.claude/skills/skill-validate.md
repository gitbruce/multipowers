# skill-validate

Comprehensive validation skill using typed validation commands.

## Overview

This skill provides validation capabilities using atomic mp commands to
check workspace, no-shell runtime, TDD environment, test results, and coverage.

## Validation Types

### Workspace Validation

**Goal:** Ensure workspace has required context files.

**Action:**
```bash
mp validate --type workspace --dir . --json
```

**Response:**
- `data.valid=true`: Workspace is ready
- `data.valid=false`: Check `data.reason` for issues

### No-Shell Runtime Validation

**Goal:** Ensure no shell script references in runtime files.

**Action:**
```bash
mp validate --type no-shell --dir . --json
```

**Response:**
- `data.valid=true`: No shell references found
- `data.valid=false`: Check `data.violations` for specific issues

### TDD Environment Validation

**Goal:** Ensure TDD environment is properly configured.

**Action:**
```bash
mp validate --type tdd-env --dir . --json
```

**Response:**
- `data.valid=true`: TDD environment ready
- `data.details.test_framework`: Test framework info
- `data.details.coverage_tool`: Coverage tool info

### Test Run Validation

**Goal:** Verify all tests pass.

**Action:**
```bash
mp validate --type test-run --dir . --json
```

**Alternative:** Use test command directly
```bash
mp test run --dir . --json
```

**Response:**
- `data.valid=true` or `data.status=passed`: All tests pass
- `data.failed_tests`: List of failing tests if any

### Coverage Validation

**Goal:** Verify test coverage meets threshold.

**Action:**
```bash
mp validate --type coverage --dir . --json
```

**Alternative:** Use coverage command directly
```bash
mp coverage check --dir . --json
```

**Response:**
- `data.valid=true`: Coverage meets threshold
- `data.coverage_pct`: Actual coverage percentage
- `data.packages`: Per-package coverage breakdown

## Comprehensive Validation

Run all validation types in sequence:

```bash
# 1. Workspace
mp validate --type workspace --dir . --json

# 2. TDD Environment
mp validate --type tdd-env --dir . --json

# 3. Tests
mp test run --dir . --json

# 4. Coverage
mp coverage check --dir . --json

# 5. No-Shell
mp validate --type no-shell --dir . --json
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

| Validation Type | Failure Action |
|-----------------|----------------|
| workspace | Run `/mp:init` to create context |
| no-shell | Remove shell script references |
| tdd-env | Set up test framework |
| test-run | Fix failing tests |
| coverage | Add more tests |

## Example Usage

User: "/mp:validate workspace readiness"

1. Run workspace validation
2. Report results and remediation steps

User: "/mp:validate full"

1. Run all validation types
2. Report comprehensive status
3. Provide remediation for any failures
