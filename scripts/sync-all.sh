#!/usr/bin/env bash
set -euo pipefail

exec ./scripts/mp-devx -action sync-all "$@"
