#!/usr/bin/env python3
"""Wrapper for Gemini CLI with structured logging and error propagation."""

import os
import subprocess
import sys
import time
from typing import Optional

from utils import count_tokens, log_execution, log_structured


def _normalize_args(args: Optional[list[str]]) -> list[str]:
    role_args = list(args) if args else []

    if role_args.count("-p") > 1:
        print("[GEMINI.PY WARNING] Multiple '-p' args detected; using one.", file=sys.stderr)

    role_args = [item for item in role_args if item != "-p"]
    role_args.insert(0, "-p")
    return role_args


def _runtime_role(default: str) -> str:
    role = os.environ.get("MULTIPOWERS_ROLE", "").strip()
    return role or default


def _runtime_request_id() -> Optional[str]:
    request_id = os.environ.get("MULTIPOWERS_REQUEST_ID", "").strip()
    return request_id or None


def call_gemini(prompt: str, args: Optional[list[str]] = None) -> str:
    role_args = _normalize_args(args)
    tokens = count_tokens(prompt)
    role_name = _runtime_role("architect")
    request_id = _runtime_request_id()
    metadata = {"request_id": request_id} if request_id else None

    log_execution(role_name, "gemini", tokens)

    cmd = ["gemini", *role_args, prompt]
    masked_cmd = ["gemini", *role_args, "<prompt>"]

    start_time = time.time()
    stderr_summary = None

    try:
        result = subprocess.run(cmd, capture_output=True, text=True, check=False)
        duration_ms = int((time.time() - start_time) * 1000)

        if result.returncode != 0:
            print(f"[GEMINI.PY ERROR] Command failed with exit code {result.returncode}", file=sys.stderr)
            print(f"[GEMINI.PY ERROR] Command: {' '.join(masked_cmd)}", file=sys.stderr)
            if result.stderr:
                stderr_lines = [line for line in result.stderr.strip().split("\n") if line.strip()]
                if stderr_lines:
                    print("[GEMINI.PY ERROR] Error output (first 5 lines):", file=sys.stderr)
                    for line in stderr_lines[:5]:
                        print(f"  {line}", file=sys.stderr)
                    stderr_summary = " | ".join(stderr_lines[:3])

            log_structured(
                role_name,
                "gemini",
                result.returncode,
                duration_ms,
                tokens,
                stderr_summary,
                "command_failed",
                metadata=metadata,
            )
            sys.exit(result.returncode)

        log_structured(role_name, "gemini", 0, duration_ms, tokens, metadata=metadata)
        return result.stdout
    except Exception as exc:
        duration_ms = int((time.time() - start_time) * 1000)
        error_msg = str(exc)
        print(f"[GEMINI.PY ERROR] Failed to call gemini: {error_msg}", file=sys.stderr)
        log_structured(role_name, "gemini", 1, duration_ms, tokens, error_msg, "exception", metadata=metadata)
        sys.exit(1)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 gemini.py '<prompt>' [args...]", file=sys.stderr)
        print("Example: python3 gemini.py 'Brainstorm architecture for auth system' -p", file=sys.stderr)
        sys.exit(1)

    prompt = sys.argv[1]
    args = sys.argv[2:] if len(sys.argv) > 2 else []

    output = call_gemini(prompt, args)
    print(output, end="")
