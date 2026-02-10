#!/usr/bin/env python3
"""Check conductor/template sync candidate drift for maintainers."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Check template sync candidates")
    parser.add_argument(
        "--mapping-file",
        default="config/template-sync-rules.json",
        help="Mapping config JSON path",
    )
    parser.add_argument(
        "--threshold",
        type=int,
        default=0,
        help="Allowed drift count before non-zero exit",
    )
    parser.add_argument("--json", action="store_true", help="Emit JSON output")
    parser.add_argument("--quiet", action="store_true", help="Suppress success output")
    return parser.parse_args()


def load_mappings(path: Path) -> list[tuple[str, str]]:
    if not path.exists():
        return []

    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        raise SystemExit(f"[TEMPLATE-SYNC] invalid mapping file: {path}: {exc}")

    if not isinstance(payload, dict):
        return []

    raw = payload.get("mappings", [])
    mappings: list[tuple[str, str]] = []
    if not isinstance(raw, list):
        return mappings

    for item in raw:
        if not isinstance(item, dict):
            continue
        source = item.get("source")
        target = item.get("target")
        if not isinstance(source, str) or not isinstance(target, str):
            continue
        if not source.strip() or not target.strip():
            continue
        mappings.append((source.strip(), target.strip()))

    return mappings


def compare_files(source: Path, target: Path) -> tuple[bool, str]:
    source_exists = source.exists()
    target_exists = target.exists()

    if not source_exists and not target_exists:
        return False, "both missing"
    if source_exists and not target_exists:
        return True, "target missing"
    if not source_exists and target_exists:
        return True, "source missing"

    source_text = source.read_text(encoding="utf-8")
    target_text = target.read_text(encoding="utf-8")
    if source_text != target_text:
        return True, "content differs"

    return False, "in sync"


def main() -> int:
    args = parse_args()
    mappings = load_mappings(Path(args.mapping_file))

    drifts: list[dict[str, str]] = []
    for source_raw, target_raw in mappings:
        source = Path(source_raw)
        target = Path(target_raw)
        drifted, reason = compare_files(source, target)
        if not drifted:
            continue

        drifts.append(
            {
                "source": source_raw,
                "target": target_raw,
                "reason": reason,
            }
        )

    payload = {
        "mapping_file": args.mapping_file,
        "mapping_count": len(mappings),
        "drift_count": len(drifts),
        "threshold": args.threshold,
        "drifts": drifts,
    }

    if args.json:
        print(json.dumps(payload, ensure_ascii=False, sort_keys=True))
    elif drifts:
        print(f"[TEMPLATE-SYNC] drift count: {len(drifts)}")
        for item in drifts:
            print(
                f"[TEMPLATE-SYNC] drift: {item['source']} -> {item['target']} ({item['reason']})"
            )
    elif not args.quiet:
        print(f"[TEMPLATE-SYNC] PASS ({len(mappings)} mappings)")

    if len(drifts) > args.threshold:
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
