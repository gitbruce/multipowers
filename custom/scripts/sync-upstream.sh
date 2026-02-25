#!/usr/bin/env bash
set -euo pipefail

UPSTREAM_REMOTE="${UPSTREAM_REMOTE:-upstream}"
MAIN_BRANCH="${MAIN_BRANCH:-main}"
TARGET_BRANCH="${TARGET_BRANCH:-multipowers}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$REPO_ROOT"

current_branch="$(git branch --show-current || true)"

if [[ -d .git/rebase-merge || -d .git/rebase-apply ]]; then
  echo "ERROR: rebase is in progress. Finish or abort it before sync." >&2
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "ERROR: working tree is not clean. Commit or stash changes before sync." >&2
  exit 1
fi

git fetch "$UPSTREAM_REMOTE"
git fetch origin

git switch "$MAIN_BRANCH"
git merge --ff-only "$UPSTREAM_REMOTE/$MAIN_BRANCH"

git switch "$TARGET_BRANCH"
if ! git merge "$MAIN_BRANCH" -m "chore(sync): merge main into $TARGET_BRANCH"; then
  echo "ERROR: merge conflict while syncing $TARGET_BRANCH with $MAIN_BRANCH." >&2
  echo "Resolve conflicts, then rerun overlay script manually:" >&2
  echo "  ./custom/scripts/apply-custom-overlay.sh" >&2
  exit 1
fi

"$SCRIPT_DIR/apply-custom-overlay.sh"

if rg -n "^(<<<<<<<|=======|>>>>>>>)" scripts/orchestrate.sh >/dev/null 2>&1; then
  echo "ERROR: conflict markers detected in scripts/orchestrate.sh after sync." >&2
  exit 1
fi

if ! bash -n scripts/orchestrate.sh; then
  echo "ERROR: scripts/orchestrate.sh has shell syntax errors after sync." >&2
  exit 1
fi

echo "Sync complete for $TARGET_BRANCH"

if [[ -n "$current_branch" && "$current_branch" != "$TARGET_BRANCH" ]]; then
  git switch "$current_branch"
fi
