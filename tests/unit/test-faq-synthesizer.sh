#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT

export PROJECT_ROOT="$TEST_ROOT"
source "$ROOT/custom/lib/faq-synthesizer.sh"

record_failure_event "plan" "context-guard" "none" "missing context" "required file missing" "run /octo:init" 1
record_failure_event "plan" "context-guard" "none" "missing context" "required file missing" "run /octo:init" 1

FAQ_FILE="$TEST_ROOT/.multipowers/FAQ.md"
[[ -f "$FAQ_FILE" ]]
rg -n '^## missing-context$' "$FAQ_FILE" >/dev/null
rg -n 'Seen: 2' "$FAQ_FILE" >/dev/null

echo "PASS test-faq-synthesizer"
