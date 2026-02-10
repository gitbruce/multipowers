#!/usr/bin/env bash
# Context quality checker tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing context quality checks..."

# Test 1: Template context passes quality check

echo "[TEST 1] Template context quality passes"
if python3 scripts/check_context_quality.py --context-dir templates/conductor/context --quiet >/dev/null 2>&1; then
    echo "  [PASS] Template context passes"
else
    echo "  [FAIL] Template context should pass"
    exit 1
fi

# Test 2: Placeholder content fails

echo "[TEST 2] Placeholder content fails quality check"
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT
mkdir -p "$tmp_dir/context"

cat > "$tmp_dir/context/product.md" <<'EOF_PRODUCT'
# Product: [Project Name]

TODO: fill this file
EOF_PRODUCT
cat > "$tmp_dir/context/product-guidelines.md" <<'EOF_GUIDE'
# Product Guidelines
- Valid line 1
- Valid line 2
- Valid line 3
- Valid line 4
- Valid line 5
- Valid line 6
EOF_GUIDE
cat > "$tmp_dir/context/workflow.md" <<'EOF_WORKFLOW'
# Workflow
- Valid line 1
- Valid line 2
- Valid line 3
- Valid line 4
- Valid line 5
- Valid line 6
EOF_WORKFLOW
cat > "$tmp_dir/context/tech-stack.md" <<'EOF_STACK'
# Tech Stack
- Valid line 1
- Valid line 2
- Valid line 3
- Valid line 4
- Valid line 5
- Valid line 6
EOF_STACK

set +e
output=$(python3 scripts/check_context_quality.py --context-dir "$tmp_dir/context" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected placeholder context to fail"
    exit 1
fi
if echo "$output" | grep -q "contains placeholder/TODO content"; then
    echo "  [PASS] Placeholder content detected"
else
    echo "  [FAIL] Expected placeholder detection"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "Context quality tests PASSED"
