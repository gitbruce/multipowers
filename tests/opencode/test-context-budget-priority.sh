#!/usr/bin/env bash
# Context budget priority trimming tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Testing context budget priority trimming..."

workspace=$(mktemp -d)
tmp_bin_dir=$(mktemp -d)
prompt_capture=$(mktemp)
stderr_capture=$(mktemp)

cleanup() {
    rm -rf "$workspace" "$tmp_bin_dir"
    rm -f "$prompt_capture" "$stderr_capture"
}
trap cleanup EXIT

cp -R "$REPO_ROOT/bin" "$workspace/"
cp -R "$REPO_ROOT/config" "$workspace/"
cp -R "$REPO_ROOT/connectors" "$workspace/"
cp -R "$REPO_ROOT/scripts" "$workspace/"
chmod +x "$workspace/bin/ask-role"

mkdir -p "$workspace/conductor/context"
cat > "$workspace/conductor/context/product.md" <<'EOF_PRODUCT'
# Product
P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1 P1
EOF_PRODUCT

cat > "$workspace/conductor/context/product-guidelines.md" <<'EOF_GUIDE'
# Product Guidelines
G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1 G1
EOF_GUIDE

cat > "$workspace/conductor/context/workflow.md" <<'EOF_WORKFLOW'
# Workflow
W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1 W1
EOF_WORKFLOW

cat > "$workspace/conductor/context/tech-stack.md" <<'EOF_STACK'
# Tech Stack
T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1 T1
EOF_STACK

cat > "$tmp_bin_dir/gemini" <<'EOF_GEMINI'
#!/usr/bin/env python3
import os
import pathlib
import sys

pathlib.Path(os.environ["PROMPT_CAPTURE"]).write_text(sys.argv[-1], encoding="utf-8")
print("stub-gemini-ok")
EOF_GEMINI
chmod +x "$tmp_bin_dir/gemini"

log_file="$workspace/outputs/runs/$(date +%Y-%m-%d).jsonl"
before_lines=0
if [ -f "$log_file" ]; then
    before_lines=$(wc -l < "$log_file")
fi

echo "[TEST 1] Priority files are kept under tight budget"
(
    cd "$workspace"
    PATH="$tmp_bin_dir:$PATH" \
    PROMPT_CAPTURE="$prompt_capture" \
    MULTIPOWERS_CONTEXT_BUDGET=80 \
    ./bin/ask-role architect "budget-check" > /dev/null 2>"$stderr_capture"
)

captured_prompt=$(cat "$prompt_capture")

if echo "$captured_prompt" | grep -q -- "--- FILE: product.md ---" && \
   echo "$captured_prompt" | grep -q -- "--- FILE: product-guidelines.md ---"; then
    echo "  [PASS] High-priority files are preserved"
else
    echo "  [FAIL] High-priority files missing from prompt"
    exit 1
fi

if echo "$captured_prompt" | grep -q -- "--- FILE: tech-stack.md ---"; then
    echo "  [FAIL] Lowest-priority file should be trimmed under tight budget"
    exit 1
else
    echo "  [PASS] Low-priority file trimmed"
fi

echo "[TEST 2] Truncation decisions are traceable in stderr"
stderr_output=$(cat "$stderr_capture")
if echo "$stderr_output" | grep -q "Context exceeded budget" && \
   echo "$stderr_output" | grep -q "Truncated files:"; then
    echo "  [PASS] Truncation warning includes trimmed files"
else
    echo "  [FAIL] Missing truncation diagnostics"
    echo "  STDERR: $stderr_output"
    exit 1
fi

echo "[TEST 3] Structured logs include truncation and request_id linkage"
if [ ! -f "$log_file" ]; then
    echo "  [FAIL] Expected structured log file: $log_file"
    exit 1
fi

if ! python3 - "$log_file" "$before_lines" <<'PY'; then
import json
import sys

log_file = sys.argv[1]
before = int(sys.argv[2])

with open(log_file, "r", encoding="utf-8") as handle:
    lines = [line.strip() for line in handle if line.strip()]

new_lines = lines[before:]
entries = []
for line in new_lines:
    try:
        entries.append(json.loads(line))
    except json.JSONDecodeError:
        continue

context_entries = [
    entry
    for entry in entries
    if entry.get("event") == "context_prepared" and entry.get("role") == "architect"
]
connector_entries = [
    entry
    for entry in entries
    if entry.get("tool") == "gemini" and entry.get("role") == "architect"
]

if not context_entries:
    print("No context_prepared entry found", file=sys.stderr)
    raise SystemExit(1)
if not connector_entries:
    print("No connector entry found", file=sys.stderr)
    raise SystemExit(1)

context_latest = context_entries[-1]
connector_latest = connector_entries[-1]

truncated_files = context_latest.get("truncated_files") or []
if "tech-stack.md" not in truncated_files:
    print(f"Expected tech-stack.md in truncated_files, got: {truncated_files}", file=sys.stderr)
    raise SystemExit(1)

if context_latest.get("context_budget") != 80:
    print(f"Expected context_budget=80, got: {context_latest.get('context_budget')}", file=sys.stderr)
    raise SystemExit(1)

request_id_context = context_latest.get("request_id")
request_id_connector = connector_latest.get("request_id")
if not request_id_context or not request_id_connector:
    print("Missing request_id in context/connector log entries", file=sys.stderr)
    raise SystemExit(1)
if request_id_context != request_id_connector:
    print(
        f"request_id mismatch: context={request_id_context} connector={request_id_connector}",
        file=sys.stderr,
    )
    raise SystemExit(1)

print("ok")
PY
    echo "  [FAIL] Structured log missing expected fields"
    exit 1
fi

echo "  [PASS] Structured log contains truncation + request linkage"

echo ""
echo "Context budget priority tests PASSED"
