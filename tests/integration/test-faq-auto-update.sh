#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT

export PROJECT_ROOT="$TEST_ROOT"
source "$ROOT/custom/lib/faq-synthesizer.sh"

record_failure_event "debate" "grapple" "multi-provider" "provider quorum failure" "fewer than 2 providers" "check provider availability" 1
record_failure_event "debate" "grapple" "multi-provider" "provider quorum failure" "fewer than 2 providers" "check provider availability" 1
record_failure_event "debate" "grapple" "gemini" "timeout" "operation timed out" "increase timeout and retry" 124

FAQ_FILE="$TEST_ROOT/.multipowers/FAQ.md"
EVENTS_FILE="$TEST_ROOT/.multipowers/temp/events/failures.ndjson"

[[ -f "$FAQ_FILE" ]]
[[ -f "$EVENTS_FILE" ]]
rg -n '^## provider-capacity$|^## timeout$|^## model-unavailable$|^## missing-context$' "$FAQ_FILE" >/dev/null
rg -n 'fewer than 2 providers' "$FAQ_FILE" >/dev/null
rg -n 'Seen: 2' "$FAQ_FILE" >/dev/null

echo "PASS test-faq-auto-update"
