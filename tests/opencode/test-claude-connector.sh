#!/usr/bin/env bash
# Claude connector and dispatch tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing claude connector and ask-role dispatch..."

# Test 1: default roles config validates with claude tool enum.
echo "[TEST 1] roles schema accepts claude tool"
if python3 scripts/validate_roles.py --config config/roles.default.json --schema config/roles.schema.json --quiet >/dev/null 2>&1; then
    echo "  [PASS] Roles schema validation succeeds"
else
    echo "  [FAIL] Roles schema validation failed"
    exit 1
fi

# Test 2: ask-role dispatches claude role and logs structured output.
echo "[TEST 2] ask-role dispatches claude connector"
workspace=$(mktemp -d)
stub_bin=$(mktemp -d)

cleanup() {
    rm -rf "$workspace" "$stub_bin"
}
trap cleanup EXIT

cp -R bin "$workspace/"
cp -R config "$workspace/"
cp -R connectors "$workspace/"
cp -R scripts "$workspace/"
cp -R templates "$workspace/"
chmod +x "$workspace/bin/ask-role"

mkdir -p "$workspace/conductor/context"
cat > "$workspace/conductor/context/product.md" <<'CTX'
# Product
Claude connector test
CTX
cat > "$workspace/conductor/context/product-guidelines.md" <<'CTX'
# Guidelines
Claude connector test
CTX
cat > "$workspace/conductor/context/workflow.md" <<'CTX'
# Workflow
Claude connector test
CTX
cat > "$workspace/conductor/context/tech-stack.md" <<'CTX'
# Stack
Claude connector test
CTX

cat > "$stub_bin/claude" <<'STUB'
#!/usr/bin/env python3
import os
import sys

log_file = os.environ.get("CLAUDE_STUB_LOG")
if log_file:
    with open(log_file, "a", encoding="utf-8") as handle:
        handle.write("ARGS=" + " ".join(sys.argv[1:]) + "\n")

print("stub-claude-ok")
STUB
chmod +x "$stub_bin/claude"

log_file="$workspace/claude.stub.log"
out_file="$workspace/ask_role.out"
err_file="$workspace/ask_role.err"

(
    cd "$workspace"
    CLAUDE_STUB_LOG="$log_file" PATH="$stub_bin:$PATH" ./bin/ask-role reviewer-claude "Review this patch" >"$out_file" 2>"$err_file"
)

if grep -q "stub-claude-ok" "$out_file" && grep -q "Request ID" "$err_file"; then
    echo "  [PASS] ask-role reached claude connector"
else
    echo "  [FAIL] ask-role did not call claude connector as expected"
    echo "  stdout:"; cat "$out_file"
    echo "  stderr:"; cat "$err_file"
    exit 1
fi

if [ ! -f "$log_file" ] || ! grep -q "ARGS=-p" "$log_file"; then
    echo "  [FAIL] Claude CLI args were not normalized/passed"
    [ -f "$log_file" ] && cat "$log_file"
    exit 1
fi

structured_log="$workspace/outputs/runs/$(date +%Y-%m-%d).jsonl"
if [ ! -f "$structured_log" ]; then
    echo "  [FAIL] structured log file missing"
    exit 1
fi

if python3 - "$structured_log" <<'PY'
import json
import sys

entries = []
with open(sys.argv[1], "r", encoding="utf-8") as handle:
    for line in handle:
        line = line.strip()
        if not line:
            continue
        entries.append(json.loads(line))

claude_entries = [entry for entry in entries if entry.get("tool") == "claude"]
if not claude_entries:
    raise SystemExit("no claude log entries")

latest = claude_entries[-1]
if latest.get("role") != "reviewer-claude":
    raise SystemExit(f"expected role reviewer-claude, got {latest.get('role')}")

print("ok")
PY
then
    echo "  [PASS] Structured log metadata is correct"
else
    echo "  [FAIL] Structured log metadata invalid"
    exit 1
fi

echo ""
echo "Claude connector tests PASSED"
