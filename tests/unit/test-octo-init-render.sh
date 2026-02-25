#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

for f in \
  "$ROOT/custom/templates/conductor/product.md" \
  "$ROOT/custom/templates/conductor/product-guidelines.md" \
  "$ROOT/custom/templates/conductor/tech-stack.md" \
  "$ROOT/custom/templates/conductor/workflow.md" \
  "$ROOT/custom/templates/conductor/tracks.md"; do
  [[ -f "$f" ]] || { echo "missing template: $f"; exit 1; }
done

rg -n "\{\{PROJECT_NAME\}\}|\{\{PRODUCT_SUMMARY\}\}" "$ROOT/custom/templates/conductor" >/dev/null

echo "PASS test-octo-init-render"
