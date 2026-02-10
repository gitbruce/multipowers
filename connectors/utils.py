#!/usr/bin/env python3
"""Utility functions for multipowers connectors."""

from __future__ import annotations

import argparse
import json
import os
import sys
from datetime import datetime
from typing import Any, Optional


def count_tokens(text: str) -> int:
    # ~4 characters per token heuristic for basic monitoring
    return len(text) // 4


def sanitize_prompt(prompt: str) -> str:
    """Legacy helper kept for backward compatibility.

    Connectors now pass prompt directly when using subprocess argument arrays.
    """
    return prompt


def truncate_context(context: str, max_tokens: int = 128000) -> str:
    tokens = count_tokens(context)
    if tokens <= max_tokens:
        return context

    safe_chars = max_tokens * 4
    if len(context) > safe_chars:
        truncated = context[:safe_chars]
        return f"{truncated}\n\n[...CONTEXT TRUNCATED...]"
    return context


def log_execution(role: str, tool: str, token_count: int):
    print(f"[ASK-ROLE] Role: {role}, Tool: {tool}, Tokens: {token_count}", file=sys.stderr)


def _log_file_path() -> str:
    log_dir = "outputs/runs"
    os.makedirs(log_dir, exist_ok=True)
    today = datetime.now().strftime("%Y-%m-%d")
    return os.path.join(log_dir, f"{today}.jsonl")


def _write_jsonl_entry(entry: dict[str, Any]):
    log_file = _log_file_path()
    with open(log_file, "a", encoding="utf-8") as handle:
        handle.write(json.dumps(entry, ensure_ascii=False) + "\n")


def _resolve_request_id(metadata: Optional[dict[str, Any]]) -> Optional[str]:
    if metadata and metadata.get("request_id"):
        return str(metadata["request_id"])

    env_request_id = os.environ.get("MULTIPOWERS_REQUEST_ID", "").strip()
    if env_request_id:
        return env_request_id
    return None


def log_structured(
    role: str,
    tool: str,
    exit_code: int,
    duration_ms: int,
    token_estimate: int,
    error_summary: Optional[str] = None,
    error_class: Optional[str] = None,
    metadata: Optional[dict[str, Any]] = None,
):
    """Write structured log entry to outputs/runs/YYYY-MM-DD.jsonl."""
    entry: dict[str, Any] = {
        "timestamp": datetime.now().isoformat(),
        "role": role,
        "tool": tool,
        "exit_code": exit_code,
        "duration_ms": duration_ms,
        "token_estimate": token_estimate,
    }

    request_id = _resolve_request_id(metadata)
    if request_id:
        entry["request_id"] = request_id

    if error_summary:
        entry["error_summary"] = error_summary
    if error_class:
        entry["error_class"] = error_class
    if metadata:
        for key, value in metadata.items():
            if value is None or key == "request_id":
                continue
            entry[key] = value

    _write_jsonl_entry(entry)


def log_context_event(
    role: str,
    token_estimate: int,
    context_budget: int,
    context_files: list[str],
    truncated_files: list[str],
    truncated: bool,
    request_id: Optional[str] = None,
):
    metadata: dict[str, Any] = {
        "event": "context_prepared",
        "context_budget": context_budget,
        "context_file_count": len(context_files),
        "context_files": context_files,
        "truncated": truncated,
        "truncated_files": truncated_files,
    }
    if request_id:
        metadata["request_id"] = request_id

    log_structured(
        role=role,
        tool="ask-role",
        exit_code=0,
        duration_ms=0,
        token_estimate=token_estimate,
        error_class="context_prepared",
        metadata=metadata,
    )


def validate_role(role: str, config: dict) -> bool:
    if "roles" not in config:
        raise ValueError("Invalid config: missing 'roles' key")

    if role not in config["roles"]:
        return False

    return True


def get_role_config(role: str, config: dict) -> dict:
    if not validate_role(role, config):
        available = ", ".join(config["roles"].keys())
        raise ValueError(f"Role '{role}' not found. Available roles: {available}")

    return config["roles"][role]


def parse_json_config(config_path: str) -> Optional[dict]:
    try:
        with open(config_path, "r", encoding="utf-8") as handle:
            return json.load(handle)
    except FileNotFoundError:
        return None
    except json.JSONDecodeError as exc:
        print(f"[ASK-ROLE ERROR] Invalid JSON in {config_path}: {exc}", file=sys.stderr)
        sys.exit(1)


def _parse_csv(value: str) -> list[str]:
    if not value:
        return []
    return [item for item in value.split(",") if item]


def _parse_bool(value: str) -> bool:
    lowered = value.strip().lower()
    return lowered in {"1", "true", "yes", "y"}


def _run_context_log_cli(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description="Write ask-role context structured log")
    parser.add_argument("--role", required=True)
    parser.add_argument("--token-estimate", type=int, required=True)
    parser.add_argument("--budget", type=int, required=True)
    parser.add_argument("--context-files", default="")
    parser.add_argument("--truncated-files", default="")
    parser.add_argument("--truncated", default="false")
    parser.add_argument("--request-id", default="")
    args = parser.parse_args(argv)

    request_id = args.request_id.strip() or os.environ.get("MULTIPOWERS_REQUEST_ID", "").strip() or None

    log_context_event(
        role=args.role,
        token_estimate=args.token_estimate,
        context_budget=args.budget,
        context_files=_parse_csv(args.context_files),
        truncated_files=_parse_csv(args.truncated_files),
        truncated=_parse_bool(args.truncated),
        request_id=request_id,
    )
    return 0


def main(argv: list[str]) -> int:
    if len(argv) > 0 and argv[0] == "context-log":
        return _run_context_log_cli(argv[1:])

    test_prompt = "This is a test prompt with special chars: $HOME, `backtick`, \"quotes\", (parens)"
    print(f"Original: {test_prompt}")
    print(f"Sanitized: {sanitize_prompt(test_prompt)}")
    print(f"Token count (approx): {count_tokens(test_prompt)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
