#!/usr/bin/env bash
# Main test runner for OpenCode plugin test suite
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================"
echo " OpenCode Plugin Test Suite"
echo "========================================"
echo ""
echo "Repository: $(cd ../.. && pwd)"
echo "Test time: $(date)"
echo ""

RUN_INTEGRATION=false
VERBOSE=false
SPECIFIC_TEST=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --integration|-i)
            RUN_INTEGRATION=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --test|-t)
            SPECIFIC_TEST="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --integration, -i  Run integration tests (requires OpenCode)"
            echo "  --verbose, -v      Show verbose output"
            echo "  --test, -t NAME    Run only the specified test"
            echo "  --help, -h         Show this help"
            echo ""
            echo "Core tests:"
            echo "  test-plugin-loading.sh"
            echo "  test-skills-core.sh"
            echo "  test-doctor-init.sh"
            echo "  test-onboarding-smoke.sh"
            echo "  test-context-quality.sh"
            echo "  test-ask-role-core.sh"
            echo "  test-config-priority.sh"
            echo "  test-ask-role-args.sh"
            echo "  test-connector-exit-code.sh"
            echo "  test-claude-connector.sh"
            echo "  test-prompt-preserve.sh"
            echo "  test-roles-schema.sh"
            echo "  test-track-workflow.sh"
            echo "  test-routing-lanes.sh"
            echo "  test-workflow-engine.sh"
            echo "  test-mcp-config.sh"
            echo "  test-governance-checks.sh"
            echo "  test-context-budget-priority.sh"
            echo "  test-plan-evidence.sh"
            echo "  test-routing-observability.sh"
            echo "  test-multipowers-route-workflow.sh"
            echo "  test-plan-status-update.sh"
            echo "  test-full-routing-flow.sh"
            echo "  test-router-run-command.sh"
            echo "  test-active-track-enforcement.sh"
            echo "  test-workflow-role-preflight.sh"
            echo "  test-governance-strictness.sh"
            echo "  test-governance-artifact.sh"
            echo "  test-track-complete-governance-gate.sh"
            echo "  test-plan-governance-evidence.sh"
            echo "  test-docs-sync-mapping.sh"
            echo "  test-observability-lifecycle.sh"
            echo "  test-update-command.sh"
            echo "  test-template-sync-candidates.sh"
            echo ""
            echo "Integration tests:"
            echo "  test-tools.sh"
            echo "  test-priority.sh"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Core tests: no external OpenCode runtime required
# Keep order deterministic and avoid destructive tests.
tests=(
    "test-plugin-loading.sh"
    "test-skills-core.sh"
    "test-doctor-init.sh"
    "test-onboarding-smoke.sh"
    "test-context-quality.sh"
    "test-ask-role-core.sh"
    "test-config-priority.sh"
    "test-ask-role-args.sh"
    "test-connector-exit-code.sh"
    "test-claude-connector.sh"
    "test-prompt-preserve.sh"
    "test-roles-schema.sh"
    "test-track-workflow.sh"
    "test-routing-lanes.sh"
    "test-workflow-engine.sh"
    "test-mcp-config.sh"
    "test-governance-checks.sh"
    "test-context-budget-priority.sh"
    "test-plan-evidence.sh"
    "test-routing-observability.sh"
    "test-multipowers-route-workflow.sh"
    "test-plan-status-update.sh"
    "test-full-routing-flow.sh"
    "test-router-run-command.sh"
    "test-active-track-enforcement.sh"
    "test-workflow-role-preflight.sh"
    "test-governance-strictness.sh"
    "test-governance-artifact.sh"
    "test-track-complete-governance-gate.sh"
    "test-plan-governance-evidence.sh"
    "test-docs-sync-mapping.sh"
    "test-observability-lifecycle.sh"
    "test-update-command.sh"
    "test-template-sync-candidates.sh"
)

integration_tests=(
    "test-tools.sh"
    "test-priority.sh"
)

if [ "$RUN_INTEGRATION" = true ]; then
    tests+=("${integration_tests[@]}")
fi

if [ -n "$SPECIFIC_TEST" ]; then
    tests=("$SPECIFIC_TEST")
fi

passed=0
failed=0
skipped=0

for test in "${tests[@]}"; do
    echo "----------------------------------------"
    echo "Running: $test"
    echo "----------------------------------------"

    test_path="$SCRIPT_DIR/$test"

    if [ ! -f "$test_path" ]; then
        echo "  [SKIP] Test file not found: $test"
        skipped=$((skipped + 1))
        continue
    fi

    if [ ! -x "$test_path" ]; then
        chmod +x "$test_path"
    fi

    start_time=$(date +%s)

    if [ "$VERBOSE" = true ]; then
        if bash "$test_path"; then
            end_time=$(date +%s)
            duration=$((end_time - start_time))
            echo ""
            echo "  [PASS] $test (${duration}s)"
            passed=$((passed + 1))
        else
            end_time=$(date +%s)
            duration=$((end_time - start_time))
            echo ""
            echo "  [FAIL] $test (${duration}s)"
            failed=$((failed + 1))
        fi
    else
        if output=$(bash "$test_path" 2>&1); then
            end_time=$(date +%s)
            duration=$((end_time - start_time))
            echo "  [PASS] (${duration}s)"
            passed=$((passed + 1))
        else
            end_time=$(date +%s)
            duration=$((end_time - start_time))
            echo "  [FAIL] (${duration}s)"
            echo ""
            echo "  Output:"
            echo "$output" | sed 's/^/    /'
            failed=$((failed + 1))
        fi
    fi

    echo ""
done

echo "========================================"
echo " Test Results Summary"
echo "========================================"
echo ""
echo "  Passed:  $passed"
echo "  Failed:  $failed"
echo "  Skipped: $skipped"
echo ""

if [ "$RUN_INTEGRATION" = false ] && [ ${#integration_tests[@]} -gt 0 ]; then
    echo "Note: Integration tests were not run."
    echo "Use --integration flag to run tests that require OpenCode."
    echo ""
fi

if [ $failed -gt 0 ]; then
    echo "STATUS: FAILED"
    exit 1
fi

echo "STATUS: PASSED"
exit 0
