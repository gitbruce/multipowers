#!/usr/bin/env python3
"""Validate multipowers MCP configuration."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path
from typing import Any


ALLOWED_SERVER_FIELDS = {"enabled", "command", "args", "env", "description", "url"}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate MCP server configuration JSON")
    parser.add_argument("--config", required=True, help="Path to MCP config JSON")
    parser.add_argument("--quiet", action="store_true", help="Suppress success output")
    return parser.parse_args()


def _load_json(path: Path) -> dict[str, Any]:
    try:
        with path.open("r", encoding="utf-8") as handle:
            data = json.load(handle)
    except FileNotFoundError:
        print(f"[MCP-VALIDATION] Config file not found: {path}", file=sys.stderr)
        raise SystemExit(1)
    except json.JSONDecodeError as exc:
        print(f"[MCP-VALIDATION] Invalid JSON in {path}: {exc}", file=sys.stderr)
        raise SystemExit(1)

    if not isinstance(data, dict):
        print("[MCP-VALIDATION] Top-level JSON must be an object", file=sys.stderr)
        raise SystemExit(1)
    return data


def _validate_server(name: str, payload: Any) -> list[str]:
    errors: list[str] = []

    if not isinstance(payload, dict):
        errors.append(f"mcpServers.{name}: must be an object")
        return errors

    for key in payload:
        if key not in ALLOWED_SERVER_FIELDS:
            errors.append(f"mcpServers.{name}.{key}: unknown field")

    if "enabled" in payload and not isinstance(payload["enabled"], bool):
        errors.append(f"mcpServers.{name}.enabled: must be boolean")

    if "command" in payload and not isinstance(payload["command"], str):
        errors.append(f"mcpServers.{name}.command: must be string")

    args = payload.get("args")
    if args is not None:
        if not isinstance(args, list):
            errors.append(f"mcpServers.{name}.args: must be array")
        else:
            for idx, item in enumerate(args):
                if not isinstance(item, str):
                    errors.append(f"mcpServers.{name}.args[{idx}]: must be string")

    env = payload.get("env")
    if env is not None:
        if not isinstance(env, dict):
            errors.append(f"mcpServers.{name}.env: must be object")
        else:
            for env_key, env_value in env.items():
                if not isinstance(env_key, str):
                    errors.append(f"mcpServers.{name}.env key must be string")
                if not isinstance(env_value, str):
                    errors.append(f"mcpServers.{name}.env.{env_key}: must be string")

    return errors


def validate_mcp_config(data: dict[str, Any]) -> list[str]:
    errors: list[str] = []

    servers = data.get("mcpServers")
    if not isinstance(servers, dict):
        errors.append("mcpServers: missing or not an object")
        return errors

    if not servers:
        errors.append("mcpServers: must not be empty")
        return errors

    for server_name, server_payload in servers.items():
        if not isinstance(server_name, str) or not server_name.strip():
            errors.append("mcpServers key: must be non-empty string")
            continue
        errors.extend(_validate_server(server_name, server_payload))

    return errors


def main() -> int:
    args = parse_args()
    config_path = Path(args.config)
    data = _load_json(config_path)
    errors = validate_mcp_config(data)

    if errors:
        for error in errors:
            print(f"[MCP-VALIDATION] {error}", file=sys.stderr)
        print(
            "[MCP-VALIDATION] Fix your MCP config or run: cp config/mcp.default.json conductor/config/mcp.json",
            file=sys.stderr,
        )
        return 1

    if not args.quiet:
        print(f"[MCP-VALIDATION] PASS: {config_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
