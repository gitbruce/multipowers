#!/usr/bin/env python3
"""Write machine-readable governance artifact JSON."""

from __future__ import annotations

import argparse
import json
from datetime import datetime
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Record governance artifact JSON")
    parser.add_argument("--artifact", required=True, help="Artifact file path")
    parser.add_argument("--mode", required=True, choices=["strict", "advisory"], help="Governance mode")
    parser.add_argument("--overall-exit-code", required=True, type=int, help="Overall governance exit code")
    parser.add_argument("--request-id", default="", help="Optional request id")
    parser.add_argument("--track-id", default="", help="Optional track id")
    parser.add_argument("--changed-file", action="append", default=[], help="Changed file path")
    parser.add_argument(
        "--tool-result",
        action="append",
        default=[],
        help="Tool result in format: tool|status|exit_code|detail",
    )
    parser.add_argument("--summary", default="", help="Optional summary override")
    return parser.parse_args()


def parse_tool_result(raw: str) -> dict[str, object]:
    parts = raw.split("|", 3)
    if len(parts) < 4:
        return {
            "tool": raw,
            "status": "unknown",
            "exit_code": -1,
            "detail": "invalid tool-result payload",
        }

    tool, status, exit_code_raw, detail = parts
    try:
        exit_code = int(exit_code_raw)
    except ValueError:
        exit_code = -1

    return {
        "tool": tool,
        "status": status,
        "exit_code": exit_code,
        "detail": detail,
    }


def unique_keep_order(items: list[str]) -> list[str]:
    seen: set[str] = set()
    result: list[str] = []
    for item in items:
        value = item.strip()
        if not value or value in seen:
            continue
        seen.add(value)
        result.append(value)
    return result


def main() -> int:
    args = parse_args()

    changed_files = unique_keep_order(args.changed_file)
    tool_results = [parse_tool_result(item) for item in args.tool_result]

    passed = sum(1 for item in tool_results if item.get("status") == "passed")
    failed = sum(1 for item in tool_results if item.get("status") in {"failed", "missing"})
    skipped = sum(1 for item in tool_results if item.get("status") == "skipped")

    summary = args.summary.strip()
    if not summary:
        summary = f"passed={passed}, failed={failed}, skipped={skipped}, exit_code={args.overall_exit_code}"

    payload: dict[str, object] = {
        "timestamp": datetime.now().isoformat(timespec="seconds"),
        "mode": args.mode,
        "changed_files": changed_files,
        "tool_results": tool_results,
        "overall_exit_code": args.overall_exit_code,
        "summary": summary,
    }
    if args.request_id.strip():
        payload["request_id"] = args.request_id.strip()
    if args.track_id.strip():
        payload["track_id"] = args.track_id.strip()

    artifact_path = Path(args.artifact)
    artifact_path.parent.mkdir(parents=True, exist_ok=True)
    artifact_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    print(f"[GOVERNANCE] Artifact written: {artifact_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
