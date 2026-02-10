#!/usr/bin/env bash
# Safe update state checker/apply helper.
set -euo pipefail

json_mode="false"
apply_mode="false"

usage() {
    cat <<'USAGE_EOF'
Usage: scripts/check_update_state.sh [--json] [--apply]

Options:
  --json   Emit JSON output
  --apply  Apply safe ff-only pull when possible
USAGE_EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --json)
            json_mode="true"
            shift
            ;;
        --apply)
            apply_mode="true"
            shift
            ;;
        --help|-h)
            usage
            exit 0
            ;;
        *)
            echo "[UPDATE] Unknown option: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

if ! command -v git >/dev/null 2>&1; then
    echo "[UPDATE] git not available" >&2
    exit 1
fi

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "[UPDATE] not inside a git repository" >&2
    exit 1
fi

branch=$(git branch --show-current 2>/dev/null || true)
if [[ -z "$branch" ]]; then
    branch="detached"
fi

dirty="false"
if ! git diff --quiet || ! git diff --cached --quiet; then
    dirty="true"
fi

upstream=""
has_upstream="false"
if upstream=$(git rev-parse --abbrev-ref --symbolic-full-name @{upstream} 2>/dev/null); then
    has_upstream="true"
fi

ahead=0
behind=0
if [[ "$has_upstream" == "true" ]]; then
    counts=$(git rev-list --left-right --count HEAD..."$upstream" 2>/dev/null || echo "0 0")
    ahead=$(echo "$counts" | awk '{print $1}')
    behind=$(echo "$counts" | awk '{print $2}')
fi

if [[ "$apply_mode" == "true" ]]; then
    if [[ "$dirty" == "true" ]]; then
        echo "[UPDATE] FAIL: working tree is dirty; aborting apply" >&2
        exit 1
    fi

    if [[ "$has_upstream" != "true" ]]; then
        echo "[UPDATE] FAIL: no upstream branch configured; cannot apply" >&2
        exit 1
    fi

    if [[ "$behind" -eq 0 ]]; then
        if [[ "$json_mode" == "true" ]]; then
            python3 - "$branch" "$dirty" "$has_upstream" "$upstream" "$ahead" "$behind" <<'PY'
import json
import sys
print(json.dumps({
    "branch": sys.argv[1],
    "dirty": sys.argv[2] == "true",
    "has_upstream": sys.argv[3] == "true",
    "upstream": sys.argv[4],
    "ahead": int(sys.argv[5]),
    "behind": int(sys.argv[6]),
    "applied": False,
    "message": "already up to date",
}, ensure_ascii=False, sort_keys=True))
PY
        else
            echo "[UPDATE] already up to date"
        fi
        exit 0
    fi

    git pull --ff-only

    # refresh state after pull
    counts=$(git rev-list --left-right --count HEAD..."$upstream" 2>/dev/null || echo "0 0")
    ahead=$(echo "$counts" | awk '{print $1}')
    behind=$(echo "$counts" | awk '{print $2}')

    if [[ "$json_mode" == "true" ]]; then
        python3 - "$branch" "$dirty" "$has_upstream" "$upstream" "$ahead" "$behind" <<'PY'
import json
import sys
print(json.dumps({
    "branch": sys.argv[1],
    "dirty": sys.argv[2] == "true",
    "has_upstream": sys.argv[3] == "true",
    "upstream": sys.argv[4],
    "ahead": int(sys.argv[5]),
    "behind": int(sys.argv[6]),
    "applied": True,
    "message": "ff-only pull completed",
}, ensure_ascii=False, sort_keys=True))
PY
    else
        echo "[UPDATE] ff-only pull completed"
    fi

    exit 0
fi

if [[ "$json_mode" == "true" ]]; then
    python3 - "$branch" "$dirty" "$has_upstream" "$upstream" "$ahead" "$behind" <<'PY'
import json
import sys
print(json.dumps({
    "branch": sys.argv[1],
    "dirty": sys.argv[2] == "true",
    "has_upstream": sys.argv[3] == "true",
    "upstream": sys.argv[4],
    "ahead": int(sys.argv[5]),
    "behind": int(sys.argv[6]),
}, ensure_ascii=False, sort_keys=True))
PY
else
    echo "[UPDATE] branch=$branch"
    echo "[UPDATE] dirty=$dirty"
    echo "[UPDATE] has_upstream=$has_upstream"
    echo "[UPDATE] upstream=${upstream:-none}"
    echo "[UPDATE] ahead=$ahead"
    echo "[UPDATE] behind=$behind"
fi
