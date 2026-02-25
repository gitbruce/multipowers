#!/usr/bin/env bash
set -euo pipefail

resolve_custom_model_for_role() {
  local role="$1"
  local config_file="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/config/models.json"
  node -e '
const fs = require("fs");
const cfg = JSON.parse(fs.readFileSync(process.argv[1], "utf8"));
const role = process.argv[2];
const lane = (cfg.role_routing && cfg.role_routing[role]) || cfg.fallback_lane;
const model = cfg.providers && cfg.providers[lane] ? cfg.providers[lane] : "";
process.stdout.write(model);
' "$config_file" "$role"
}
