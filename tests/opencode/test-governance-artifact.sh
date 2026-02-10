#!/usr/bin/env bash
# Governance artifact output tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing governance artifact output..."

artifact=$(mktemp)
cleanup() {
    rm -f "$artifact"
}
trap cleanup EXIT

if bash scripts/run_governance_checks.sh \
    --mode advisory \
    --changed-file README.md \
    --request-id req-gov-art-001 \
    --track-id track-gov-art-001 \
    --artifact "$artifact" >/tmp/test_gov_art_out.txt 2>/tmp/test_gov_art_err.txt; then
    if python3 - "$artifact" <<'PY'
import json
import sys

with open(sys.argv[1], 'r', encoding='utf-8') as handle:
    payload = json.load(handle)

required = [
    'timestamp', 'mode', 'changed_files', 'tool_results', 'overall_exit_code', 'summary'
]
for key in required:
    if key not in payload:
        raise SystemExit(f"missing key: {key}")

if payload.get('mode') != 'advisory':
    raise SystemExit('mode mismatch')
if payload.get('request_id') != 'req-gov-art-001':
    raise SystemExit('request_id mismatch')
if payload.get('track_id') != 'track-gov-art-001':
    raise SystemExit('track_id mismatch')
if 'README.md' not in payload.get('changed_files', []):
    raise SystemExit('changed_files missing README.md')
if not isinstance(payload.get('tool_results'), list):
    raise SystemExit('tool_results not list')

print('ok')
PY
    then
        echo "  [PASS] governance artifact schema is valid"
    else
        echo "  [FAIL] governance artifact content invalid"
        cat /tmp/test_gov_art_err.txt
        exit 1
    fi
else
    echo "  [FAIL] governance run with artifact should succeed"
    cat /tmp/test_gov_art_err.txt
    exit 1
fi

echo ""
echo "Governance artifact tests PASSED"
