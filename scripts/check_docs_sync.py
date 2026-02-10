#!/usr/bin/env python3
"""Check docs synchronization against changed files."""

from __future__ import annotations

import argparse
import sys
from pathlib import Path


DOC_PATH_PREFIXES = (
    "docs/",
    "conductor/context/",
    "templates/conductor/context/",
)
DOC_SINGLE_FILES = {"README.md"}

IGNORE_PREFIXES = (
    "tests/",
    "outputs/",
    ".worktrees/",
)

IGNORE_SUFFIXES = (
    ".md",
    ".txt",
    ".log",
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate docs sync against changed files")
    parser.add_argument("--changed-file", action="append", default=[], help="Changed file path")
    parser.add_argument(
        "--changed-files-file",
        help="Path to newline-delimited changed file list",
    )
    parser.add_argument("--quiet", action="store_true", help="Suppress success output")
    return parser.parse_args()


def _load_changed_files(args: argparse.Namespace) -> list[str]:
    changed = [item.strip() for item in args.changed_file if item.strip()]

    if args.changed_files_file:
        path = Path(args.changed_files_file)
        if not path.exists():
            print(f"[DOCS-SYNC] changed file list not found: {path}", file=sys.stderr)
            raise SystemExit(1)
        for line in path.read_text(encoding="utf-8").splitlines():
            value = line.strip()
            if value:
                changed.append(value)

    unique: list[str] = []
    seen: set[str] = set()
    for item in changed:
        if item in seen:
            continue
        seen.add(item)
        unique.append(item)
    return unique


def _is_doc_file(path: str) -> bool:
    if path in DOC_SINGLE_FILES:
        return True
    return path.startswith(DOC_PATH_PREFIXES)


def _is_ignored_non_doc(path: str) -> bool:
    if path.startswith(IGNORE_PREFIXES):
        return True
    return path.endswith(IGNORE_SUFFIXES)


def validate_docs_sync(changed_files: list[str]) -> tuple[bool, str]:
    if not changed_files:
        return True, "no changed files to evaluate"

    doc_changes = [item for item in changed_files if _is_doc_file(item)]
    non_doc_changes = [item for item in changed_files if not _is_doc_file(item)]

    behavior_changes = [item for item in non_doc_changes if not _is_ignored_non_doc(item)]

    if behavior_changes and not doc_changes:
        preview = ", ".join(behavior_changes[:5])
        return (
            False,
            "docs update required for behavior/code changes. "
            f"Changed code files (sample): {preview}. "
            "Update README.md or docs/context files to describe the change.",
        )

    return True, "docs sync check passed"


def main() -> int:
    args = parse_args()
    changed_files = _load_changed_files(args)

    ok, message = validate_docs_sync(changed_files)
    if not ok:
        print(f"[DOCS-SYNC] FAIL: {message}", file=sys.stderr)
        return 1

    if not args.quiet:
        print(f"[DOCS-SYNC] PASS: {message}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
