#!/usr/bin/env bash
# MCP config validation and doctor integration tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

create_workspace() {
    local workspace
    workspace=$(mktemp -d)
    cp -R "$REPO_ROOT/bin" "$workspace/"
    cp -R "$REPO_ROOT/config" "$workspace/"
    cp -R "$REPO_ROOT/connectors" "$workspace/"
    cp -R "$REPO_ROOT/scripts" "$workspace/"
    cp -R "$REPO_ROOT/templates" "$workspace/"
    chmod +x "$workspace/bin/multipowers"
    echo "$workspace"
}

echo "Testing MCP config integration..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

# Ensure baseline context exists for doctor checks.
./bin/multipowers init --repair >/tmp/test_mcp_repair_out.txt 2>/tmp/test_mcp_repair_err.txt

# Test 1: default MCP config validates.
echo "[TEST 1] validate_mcp accepts default config"
if python3 scripts/validate_mcp.py --config config/mcp.default.json --quiet >/dev/null 2>&1; then
    echo "  [PASS] Default MCP config is valid"
else
    echo "  [FAIL] Default MCP config should validate"
    exit 1
fi

# Test 2: doctor uses project override and validates it.
echo "[TEST 2] doctor uses project MCP override"
mkdir -p conductor/config
cat > conductor/config/mcp.json <<'JSON_EOF'
{
  "mcpServers": {
    "filesystem": {
      "enabled": true,
      "command": "mcp-filesystem"
    }
  }
}
JSON_EOF

if output=$(./bin/multipowers doctor 2>&1); then
    if echo "$output" | grep -q "effective mcp config: conductor/config/mcp.json" && echo "$output" | grep -q "mcp config validates"; then
        echo "  [PASS] Doctor validates project MCP config"
    else
        echo "  [FAIL] Doctor output missing MCP override/validation details"
        echo "  Output: $output"
        exit 1
    fi
else
    echo "  [FAIL] Doctor should pass with valid project MCP config"
    echo "  Output: $output"
    exit 1
fi

# Test 3: invalid MCP config fails with actionable guidance.
echo "[TEST 3] invalid MCP config fails clearly"
cat > conductor/config/mcp.json <<'JSON_EOF'
{
  "mcpServers": {
    "filesystem": {
      "enabled": "yes"
    }
  }
}
JSON_EOF

set +e
output=$(./bin/multipowers doctor 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Doctor should fail with invalid MCP config"
    exit 1
fi
if echo "$output" | grep -q "mcp config validation failed" && echo "$output" | grep -q "Fix your MCP config"; then
    echo "  [PASS] Invalid MCP diagnostics are actionable"
else
    echo "  [FAIL] Missing actionable MCP validation diagnostics"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "MCP config tests PASSED"
