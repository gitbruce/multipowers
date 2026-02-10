#!/usr/bin/env bash
# Connector exit code propagation tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing connector exit code propagation..."

run_and_capture() {
    local cmd=$1
    local output_file=$2
    set +e
    eval "$cmd" >"$output_file" 2>&1
    local status=$?
    set -e
    echo "$status"
}

# Test 1: codex connector returns non-zero when codex command is unavailable/fails

echo "[TEST 1] codex connector propagates failure"
out1=$(mktemp)
status=$(run_and_capture "python3 connectors/codex.py 'test prompt' exec --non-existent-flag" "$out1")
if [ "$status" -eq 0 ]; then
    echo "  [FAIL] Expected non-zero exit for failing codex call"
    cat "$out1"
    rm -f "$out1"
    exit 1
fi
if grep -q "Command failed with exit code\|Failed to call codex" "$out1"; then
    echo "  [PASS] codex connector failure propagated"
else
    echo "  [FAIL] Missing codex failure diagnostics"
    cat "$out1"
    rm -f "$out1"
    exit 1
fi
rm -f "$out1"

# Test 2: gemini connector returns non-zero when gemini command is unavailable/fails

echo "[TEST 2] gemini connector propagates failure"
out2=$(mktemp)
status=$(run_and_capture "python3 connectors/gemini.py 'test prompt' -p" "$out2")
if [ "$status" -eq 0 ]; then
    echo "  [FAIL] Expected non-zero exit for failing gemini call"
    cat "$out2"
    rm -f "$out2"
    exit 1
fi
if grep -q "Command failed with exit code\|Failed to call gemini" "$out2"; then
    echo "  [PASS] gemini connector failure propagated"
else
    echo "  [FAIL] Missing gemini failure diagnostics"
    cat "$out2"
    rm -f "$out2"
    exit 1
fi
rm -f "$out2"

# Test 3: error output does not leak full prompt string in command log

echo "[TEST 3] command log masks prompt"
out3=$(mktemp)
status=$(run_and_capture "python3 connectors/codex.py 'SENSITIVE_PROMPT_CONTENT_123' exec" "$out3")
if grep -q "SENSITIVE_PROMPT_CONTENT_123" "$out3"; then
    echo "  [FAIL] Sensitive prompt leaked in error output"
    cat "$out3"
    rm -f "$out3"
    exit 1
fi
echo "  [PASS] Prompt masked in error output"
rm -f "$out3"

# Test 4: ask-role sets runtime role for structured logs

echo "[TEST 4] ask-role propagates runtime role into structured logs"
workspace=$(mktemp -d)
tmp_bin_dir=$(mktemp -d)

cleanup() {
    rm -rf "$workspace" "$tmp_bin_dir"
}
trap cleanup EXIT

cp -R "$REPO_ROOT/bin" "$workspace/"
cp -R "$REPO_ROOT/config" "$workspace/"
cp -R "$REPO_ROOT/connectors" "$workspace/"
cp -R "$REPO_ROOT/scripts" "$workspace/"
mkdir -p "$workspace/conductor/config" "$workspace/conductor/context"
chmod +x "$workspace/bin/ask-role"

cat > "$workspace/conductor/config/roles.json" <<'JSON_EOF'
{
  "roles": {
    "oracle": {
      "description": "Review role",
      "tool": "gemini",
      "system_prompt": "Review",
      "args": ["-p"]
    }
  }
}
JSON_EOF

cat > "$workspace/conductor/context/product.md" <<'EOF_PRODUCT'
# Product
Role log test
EOF_PRODUCT
cat > "$workspace/conductor/context/product-guidelines.md" <<'EOF_GUIDE'
# Guidelines
Role log test
EOF_GUIDE
cat > "$workspace/conductor/context/workflow.md" <<'EOF_WORKFLOW'
# Workflow
Role log test
EOF_WORKFLOW
cat > "$workspace/conductor/context/tech-stack.md" <<'EOF_STACK'
# Tech
Role log test
EOF_STACK

cat > "$tmp_bin_dir/gemini" <<'PY'
#!/usr/bin/env python3
print("stub-gemini-ok")
PY
chmod +x "$tmp_bin_dir/gemini"

(
    cd "$workspace"
    PATH="$tmp_bin_dir:$PATH" ./bin/ask-role oracle "verify-role-log" >/tmp/test_connector_role_out.txt 2>/tmp/test_connector_role_err.txt
)

log_file="$workspace/outputs/runs/$(date +%Y-%m-%d).jsonl"
if [ ! -f "$log_file" ]; then
    echo "  [FAIL] Structured log file missing"
    exit 1
fi

if ! python3 - "$log_file" <<'PY'; then
import json
import sys

log_file = sys.argv[1]
with open(log_file, "r", encoding="utf-8") as handle:
    entries = [json.loads(line) for line in handle if line.strip()]

connector_entries = [entry for entry in entries if entry.get("tool") == "gemini"]
if not connector_entries:
    print("No gemini connector entry found", file=sys.stderr)
    raise SystemExit(1)

latest = connector_entries[-1]
if latest.get("role") != "oracle":
    print(f"Expected role=oracle, got role={latest.get('role')}", file=sys.stderr)
    raise SystemExit(1)

print("ok")
PY
    echo "  [FAIL] Runtime role not reflected in structured log"
    exit 1
fi

echo "  [PASS] Structured log role matches ask-role caller"

echo ""
echo "Connector exit code propagation tests PASSED"
