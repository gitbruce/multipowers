# Proxy Routing

## What Changed From Upstream
Added overlay-driven proxy policy in `custom/config/proxy.json` for Codex/Gemini paths.

## Why This Exists
To route external provider traffic through controlled network egress.

## How To Use
Set host/port/providers in `custom/config/proxy.json`.

## Operational Impact
Only external CLI providers receive proxy env variables.

## Rollback Path
Set `enabled` to `false` in `custom/config/proxy.json`.
