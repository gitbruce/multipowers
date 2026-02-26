#!/usr/bin/env bash
set -euo pipefail
CUSTOM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Optional custom libraries; source only if present.
[[ -f "$CUSTOM_ROOT/lib/model-routing.sh" ]] && source "$CUSTOM_ROOT/lib/model-routing.sh"
[[ -f "$CUSTOM_ROOT/lib/proxy-routing.sh" ]] && source "$CUSTOM_ROOT/lib/proxy-routing.sh"
[[ -f "$CUSTOM_ROOT/lib/conductor-context.sh" ]] && source "$CUSTOM_ROOT/lib/conductor-context.sh"
[[ -f "$CUSTOM_ROOT/lib/faq-synthesizer.sh" ]] && source "$CUSTOM_ROOT/lib/faq-synthesizer.sh"
