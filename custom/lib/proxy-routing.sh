#!/usr/bin/env bash
set -euo pipefail

custom_proxy_url_for_provider() {
  local provider="$1"
  local config_file="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/config/proxy.json"
  node -e '
const fs = require("fs");
const cfg = JSON.parse(fs.readFileSync(process.argv[1], "utf8"));
const provider = process.argv[2];
if (!cfg.enabled) process.exit(0);
if (!Array.isArray(cfg.providers) || !cfg.providers.includes(provider)) process.exit(0);
process.stdout.write(`http://${cfg.host}:${cfg.port}`);
' "$config_file" "$provider"
}
