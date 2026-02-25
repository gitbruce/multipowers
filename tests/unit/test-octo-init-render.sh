#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

TEMPLATES_DIR="$ROOT/custom/templates/conductor"
[[ -d "$TEMPLATES_DIR" ]] || { echo "missing templates dir: $TEMPLATES_DIR"; exit 1; }

expected="$(cat <<'EOF'
code_styleguides/cpp.md
code_styleguides/csharp.md
code_styleguides/dart.md
code_styleguides/general.md
code_styleguides/go.md
code_styleguides/html-css.md
code_styleguides/javascript.md
code_styleguides/python.md
code_styleguides/typescript.md
workflow.md
EOF
)"

actual="$(cd "$TEMPLATES_DIR" && find . -type f | sed 's|^\./||' | sort)"
[[ "$actual" == "$expected" ]] || {
  echo "template set mismatch"
  echo "expected:"
  printf "%s\n" "$expected"
  echo "actual:"
  printf "%s\n" "$actual"
  exit 1
}

# Ensure /octo:init still renders required artifacts when upstream templates omit them.
rg -n 'if \[\[ ! -f "\$croot/product\.md" \]\]' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n 'if \[\[ ! -f "\$croot/product-guidelines\.md" \]\]' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n 'if \[\[ ! -f "\$croot/tech-stack\.md" \]\]' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n 'if \[\[ ! -f "\$croot/tracks\.md" \]\]' "$ROOT/scripts/orchestrate.sh" >/dev/null

echo "PASS test-octo-init-render"
