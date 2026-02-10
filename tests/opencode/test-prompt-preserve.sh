#!/usr/bin/env bash
# Prompt preservation tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing prompt preservation in connectors..."

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

cat > "$tmp_dir/codex" <<'PY'
#!/usr/bin/env python3
import json
import os
import sys

with open(os.environ["ARGS_FILE"], "w", encoding="utf-8") as handle:
    json.dump(sys.argv[1:], handle)
print("stub-codex-ok")
PY

cat > "$tmp_dir/gemini" <<'PY'
#!/usr/bin/env python3
import json
import os
import sys

with open(os.environ["ARGS_FILE"], "w", encoding="utf-8") as handle:
    json.dump(sys.argv[1:], handle)
print("stub-gemini-ok")
PY

chmod +x "$tmp_dir/codex" "$tmp_dir/gemini"

special_prompt='literal $HOME "quoted" (paren) `backtick` and spaces'

# Test codex connector
echo "[TEST 1] codex connector preserves prompt"
export PATH="$tmp_dir:$PATH"
export ARGS_FILE="$tmp_dir/codex_args.json"
python3 connectors/codex.py "$special_prompt" exec >/dev/null

captured=$(ARGS_FILE="$tmp_dir/codex_args.json" python3 - <<'PY'
import json
import os

with open(os.environ["ARGS_FILE"], "r", encoding="utf-8") as handle:
    args = json.load(handle)
print(args[-1])
PY
)

if [ "$captured" = "$special_prompt" ]; then
    echo "  [PASS] codex prompt preserved"
else
    echo "  [FAIL] codex prompt mutated"
    echo "  Expected: $special_prompt"
    echo "  Actual:   $captured"
    exit 1
fi

# Test gemini connector
echo "[TEST 2] gemini connector preserves prompt"
export ARGS_FILE="$tmp_dir/gemini_args.json"
python3 connectors/gemini.py "$special_prompt" -p >/dev/null

captured=$(ARGS_FILE="$tmp_dir/gemini_args.json" python3 - <<'PY'
import json
import os

with open(os.environ["ARGS_FILE"], "r", encoding="utf-8") as handle:
    args = json.load(handle)
print(args[-1])
PY
)

if [ "$captured" = "$special_prompt" ]; then
    echo "  [PASS] gemini prompt preserved"
else
    echo "  [FAIL] gemini prompt mutated"
    echo "  Expected: $special_prompt"
    echo "  Actual:   $captured"
    exit 1
fi

echo ""
echo "Prompt preservation tests PASSED"
