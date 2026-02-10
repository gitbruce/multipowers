#!/usr/bin/env python3
"""Check quality of conductor context files."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

DEFAULT_REQUIRED = [
    "product.md",
    "product-guidelines.md",
    "workflow.md",
    "tech-stack.md",
]

PLACEHOLDER_PATTERNS = [
    re.compile(r"\[(Project Name|Your Name|feature-name|primary users|core pain point|measurable value delivered)\]", re.IGNORECASE),
    re.compile(r"\[\.\.\.\]"),
    re.compile(r"TODO", re.IGNORECASE),
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate context file completeness and quality")
    parser.add_argument("--context-dir", default="conductor/context", help="Context directory path")
    parser.add_argument(
        "--required",
        nargs="*",
        default=DEFAULT_REQUIRED,
        help="Required context markdown filenames",
    )
    parser.add_argument("--quiet", action="store_true", help="Suppress success output")
    return parser.parse_args()


def non_empty_line_count(content: str) -> int:
    return sum(1 for line in content.splitlines() if line.strip())


def contains_placeholder(content: str) -> bool:
    if re.search(r"^-\s*\[\s*\]", content, re.MULTILINE):
        return True
    for pattern in PLACEHOLDER_PATTERNS:
        if pattern.search(content):
            return True
    return False


def validate_file(file_path: Path) -> list[str]:
    errors: list[str] = []
    try:
        content = file_path.read_text(encoding="utf-8")
    except FileNotFoundError:
        errors.append(f"missing file: {file_path}")
        return errors

    if non_empty_line_count(content) < 6:
        errors.append(f"{file_path}: content too short (need at least 6 non-empty lines)")

    if contains_placeholder(content):
        errors.append(f"{file_path}: contains placeholder/TODO content")

    return errors


def main() -> int:
    args = parse_args()
    context_dir = Path(args.context_dir)

    if not context_dir.exists():
        print(f"[CONTEXT-QUALITY] context directory missing: {context_dir}", file=sys.stderr)
        return 1

    if not context_dir.is_dir():
        print(f"[CONTEXT-QUALITY] context path is not a directory: {context_dir}", file=sys.stderr)
        return 1

    all_errors: list[str] = []
    for filename in args.required:
        all_errors.extend(validate_file(context_dir / filename))

    if all_errors:
        for issue in all_errors:
            print(f"[CONTEXT-QUALITY] {issue}", file=sys.stderr)
        return 1

    if not args.quiet:
        print(f"[CONTEXT-QUALITY] PASS: {context_dir}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
