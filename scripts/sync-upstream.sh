#!/usr/bin/env bash
# Safely sync local/fork main with upstream, then rebase custom branch on top.
# Usage:
#   ./scripts/sync-upstream.sh
#   ./scripts/sync-upstream.sh --main main --custom multipowers

set -euo pipefail

MAIN_BRANCH="${MAIN_BRANCH:-main}"
CUSTOM_BRANCH="${CUSTOM_BRANCH:-multipowers}"
UPSTREAM_REMOTE="${UPSTREAM_REMOTE:-upstream}"
ORIGIN_REMOTE="${ORIGIN_REMOTE:-origin}"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --main)
            MAIN_BRANCH="${2:-}"
            shift 2
            ;;
        --custom)
            CUSTOM_BRANCH="${2:-}"
            shift 2
            ;;
        --upstream)
            UPSTREAM_REMOTE="${2:-}"
            shift 2
            ;;
        --origin)
            ORIGIN_REMOTE="${2:-}"
            shift 2
            ;;
        -h|--help)
            cat <<'EOF'
Usage: scripts/sync-upstream.sh [options]

Options:
  --main <branch>      Mirror branch to sync from upstream (default: main)
  --custom <branch>    Custom branch to rebase on mirror branch (default: multipowers)
  --upstream <remote>  Upstream remote name (default: upstream)
  --origin <remote>    Fork remote name (default: origin)
EOF
            exit 0
            ;;
        *)
            echo "ERROR: Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "ERROR: Not inside a git repository." >&2
    exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
    echo "ERROR: Working tree is not clean. Commit or stash changes before sync." >&2
    exit 1
fi

if ! git remote get-url "$UPSTREAM_REMOTE" >/dev/null 2>&1; then
    echo "ERROR: Upstream remote '$UPSTREAM_REMOTE' is not configured." >&2
    exit 1
fi

if ! git remote get-url "$ORIGIN_REMOTE" >/dev/null 2>&1; then
    echo "ERROR: Origin remote '$ORIGIN_REMOTE' is not configured." >&2
    exit 1
fi

current_branch="$(git branch --show-current)"

echo "==> Fetching latest from $UPSTREAM_REMOTE"
git fetch "$UPSTREAM_REMOTE" --prune

if ! git show-ref --verify --quiet "refs/remotes/$UPSTREAM_REMOTE/$MAIN_BRANCH"; then
    echo "ERROR: Remote branch '$UPSTREAM_REMOTE/$MAIN_BRANCH' not found." >&2
    exit 1
fi

echo "==> Syncing $MAIN_BRANCH to $UPSTREAM_REMOTE/$MAIN_BRANCH"
git switch "$MAIN_BRANCH" >/dev/null
git reset --hard "$UPSTREAM_REMOTE/$MAIN_BRANCH"

echo "==> Pushing $MAIN_BRANCH to $ORIGIN_REMOTE (force-with-lease)"
git push "$ORIGIN_REMOTE" "$MAIN_BRANCH" --force-with-lease

if [[ "$CUSTOM_BRANCH" != "$MAIN_BRANCH" ]]; then
    echo "==> Rebasing $CUSTOM_BRANCH onto $MAIN_BRANCH"
    git switch "$CUSTOM_BRANCH" >/dev/null
    git rebase "$MAIN_BRANCH"

    echo "==> Pushing $CUSTOM_BRANCH to $ORIGIN_REMOTE (force-with-lease)"
    git push "$ORIGIN_REMOTE" "$CUSTOM_BRANCH" --force-with-lease
fi

echo "==> Restoring original branch: $current_branch"
git switch "$current_branch" >/dev/null

echo "Sync complete."

