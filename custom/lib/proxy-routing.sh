#!/usr/bin/env bash
set -euo pipefail

custom_proxy_url_for_provider() {
  local provider="${1:-}"
  local config_file
  config_file="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/config/proxy.json"
  [[ -f "$config_file" ]] || return 0
  python3 - "$config_file" "$provider" <<'PY'
import json
import socket
import subprocess
import sys

cfg_file, provider = sys.argv[1], sys.argv[2]
try:
    with open(cfg_file, "r", encoding="utf-8") as f:
        cfg = json.load(f)
except Exception:
    raise SystemExit(0)
if not cfg.get("enabled", False):
    raise SystemExit(0)
providers = cfg.get("providers", [])
if provider not in providers:
    raise SystemExit(0)

def detect_host() -> str:
    try:
        out = subprocess.check_output(
            ["ip", "route", "show", "default"],
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        if out:
            parts = out.split()
            if "via" in parts:
                idx = parts.index("via")
                if idx + 1 < len(parts):
                    return parts[idx + 1].strip()
    except Exception:
        pass

    try:
        with open("/proc/net/route", "r", encoding="utf-8") as f:
            next(f, None)
            for line in f:
                cols = line.split()
                if len(cols) < 3:
                    continue
                if cols[1] != "00000000":
                    continue
                gateway_hex = cols[2]
                gateway_bytes = bytes.fromhex(gateway_hex)
                if len(gateway_bytes) == 4:
                    return ".".join(str(b) for b in gateway_bytes)
    except Exception:
        pass

    for name in ("host.docker.internal",):
        try:
            return socket.gethostbyname(name)
        except Exception:
            continue
    return ""

host = str(cfg.get("host", "")).strip()
if host.lower() in {"", "auto", "detect", "dynamic"}:
    host = detect_host()

port = str(cfg.get("port", "")).strip()
if not host or not port:
    raise SystemExit(0)
print(f"http://{host}:{port}")
PY
}
