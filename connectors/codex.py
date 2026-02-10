#!/usr/bin/env python3
"""Wrapper for Codex CLI with structured logging and error propagation."""

import os
import subprocess
import sys
import time
from typing import Optional

from utils import count_tokens, log_execution, log_structured


def _normalize_args(args: list[str]) -> list[str]:
    normalized = list(args)

    if "exec" not in normalized:
        normalized.insert(0, "exec")

    if normalized.count("exec") > 1:
        print("[CODEX.PY WARNING] Multiple 'exec' args detected; using first occurrence.", file=sys.stderr)
        first_exec_index = normalized.index("exec")
        normalized = [item for idx, item in enumerate(normalized) if item != "exec" or idx == first_exec_index]

    return normalized


def _runtime_role(default: str) -> str:
    role = os.environ.get("MULTIPOWERS_ROLE", "").strip()
    return role or default


def _runtime_request_id() -> Optional[str]:
    request_id = os.environ.get("MULTIPOWERS_REQUEST_ID", "").strip()
    return request_id or None


def call_codex(prompt: str, args: list[str]) -> str:
    """Call Codex CLI and propagate failures with safe logging."""
    role_args = _normalize_args(args)
    tokens = count_tokens(prompt)
    role_name = _runtime_role("coder")
    request_id = _runtime_request_id()
    metadata = {"request_id": request_id} if request_id else None

    log_execution(role_name, "codex", tokens)

    cmd = ["codex", *role_args, prompt]
    masked_cmd = ["codex", *role_args, "<prompt>"]

    start_time = time.time()
    stderr_summary = None

    try:
        result = subprocess.run(cmd, capture_output=True, text=True, check=False)
        duration_ms = int((time.time() - start_time) * 1000)

        if result.returncode != 0:
            print(f"[CODEX.PY ERROR] Command failed with exit code {result.returncode}", file=sys.stderr)
            print(f"[CODEX.PY ERROR] Command: {' '.join(masked_cmd)}", file=sys.stderr)
            if result.stderr:
                stderr_lines = [line for line in result.stderr.strip().split("\n") if line.strip()]
                if stderr_lines:
                    print("[CODEX.PY ERROR] Error output (first 5 lines):", file=sys.stderr)
                    for line in stderr_lines[:5]:
                        print(f"  {line}", file=sys.stderr)
                    stderr_summary = " | ".join(stderr_lines[:3])

            log_structured(
                role_name,
                "codex",
                result.returncode,
                duration_ms,
                tokens,
                stderr_summary,
                "command_failed",
                metadata=metadata,
            )
            sys.exit(result.returncode)

        log_structured(role_name, "codex", 0, duration_ms, tokens, metadata=metadata)
        return result.stdout
    except Exception as exc:
        duration_ms = int((time.time() - start_time) * 1000)
        error_msg = str(exc)
        print(f"[CODEX.PY ERROR] Failed to call codex: {error_msg}", file=sys.stderr)
        log_structured(role_name, "codex", 1, duration_ms, tokens, error_msg, "exception", metadata=metadata)
        sys.exit(1)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 codex.py '<prompt>' [args...]", file=sys.stderr)
        print("Example: python3 codex.py 'Implement TDD for user auth' exec --skip-git-repo-check", file=sys.stderr)
        sys.exit(1)

    prompt = sys.argv[1]
    args = sys.argv[2:]

    output = call_codex(prompt, args)
    print(output, end="")
