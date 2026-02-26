#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$ROOT/bin/octo"
LEGACY="$ROOT/scripts/orchestrate_legacy.sh"

if [[ "${OCTO_RUNTIME:-go}" == "legacy" ]]; then
  exec "$LEGACY" "$@"
fi

if [[ ! -x "$BIN" ]]; then
  (cd "$ROOT" && go build -o bin/octo ./cmd/octo)
fi

# legacy alias mapping
cmd="${1:-}"
shift || true
case "$cmd" in
  probe) cmd="discover" ;;
  grasp) cmd="define" ;;
  tangle) cmd="develop" ;;
  ink) cmd="deliver" ;;
  grapple) cmd="debate" ;;
  auto) cmd="develop" ;;
  "") echo "usage: orchestrate.sh <command> [prompt]"; exit 2 ;;
esac

prompt="${*:-}"
exec "$BIN" "$cmd" --dir "$PWD" --prompt "$prompt" --json
