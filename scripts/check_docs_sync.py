#!/usr/bin/env python3
"""Check docs synchronization against changed files with component-aware mapping."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path


DOC_PATH_PREFIXES = (
    "docs/",
    "conductor/context/",
    "templates/conductor/context/",
)
DOC_SINGLE_FILES = {"README.md"}

DEFAULT_IGNORE_PREFIXES = (
    "tests/",
    "outputs/",
    ".worktrees/",
)

DEFAULT_IGNORE_SUFFIXES = (
    ".md",
    ".txt",
    ".log",
)

DEFAULT_RULES = [
    {"code_prefix": "bin/", "required_docs": ["README.md"]},
    {"code_prefix": "connectors/", "required_docs": ["README.md"]},
    {"code_prefix": "scripts/", "required_docs": ["README.md"]},
    {"code_prefix": "config/", "required_docs": ["README.md"]},
    {"code_prefix": "lib/", "required_docs": ["README.md"]},
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate docs sync against changed files")
    parser.add_argument("--changed-file", action="append", default=[], help="Changed file path")
    parser.add_argument(
        "--changed-files-file",
        help="Path to newline-delimited changed file list",
    )
    parser.add_argument(
        "--rules-file",
        default="config/docs-sync-rules.json",
        help="Rules file path",
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


def _load_rules(path: Path) -> tuple[list[dict[str, object]], tuple[str, ...], tuple[str, ...]]:
    if not path.exists():
        return DEFAULT_RULES, DEFAULT_IGNORE_PREFIXES, DEFAULT_IGNORE_SUFFIXES

    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        raise SystemExit(f"[DOCS-SYNC] invalid rules json: {path}: {exc}")

    if not isinstance(payload, dict):
        return DEFAULT_RULES, DEFAULT_IGNORE_PREFIXES, DEFAULT_IGNORE_SUFFIXES

    raw_rules = payload.get("rules", DEFAULT_RULES)
    rules: list[dict[str, object]] = []
    if isinstance(raw_rules, list):
        for item in raw_rules:
            if not isinstance(item, dict):
                continue
            code_prefix = item.get("code_prefix")
            required_docs = item.get("required_docs")
            if not isinstance(code_prefix, str) or not code_prefix:
                continue
            if not isinstance(required_docs, list) or not required_docs:
                continue
            docs = [doc for doc in required_docs if isinstance(doc, str) and doc.strip()]
            if not docs:
                continue
            rules.append({"code_prefix": code_prefix, "required_docs": docs})

    if not rules:
        rules = DEFAULT_RULES

    ignore_prefixes_raw = payload.get("ignore_prefixes", list(DEFAULT_IGNORE_PREFIXES))
    ignore_suffixes_raw = payload.get("ignore_suffixes", list(DEFAULT_IGNORE_SUFFIXES))

    ignore_prefixes = tuple(item for item in ignore_prefixes_raw if isinstance(item, str))
    ignore_suffixes = tuple(item for item in ignore_suffixes_raw if isinstance(item, str))

    return rules, ignore_prefixes or DEFAULT_IGNORE_PREFIXES, ignore_suffixes or DEFAULT_IGNORE_SUFFIXES


def _is_doc_file(path: str) -> bool:
    if path in DOC_SINGLE_FILES:
        return True
    return path.startswith(DOC_PATH_PREFIXES)


def _is_ignored_non_doc(path: str, ignore_prefixes: tuple[str, ...], ignore_suffixes: tuple[str, ...]) -> bool:
    if path.startswith(ignore_prefixes):
        return True
    return path.endswith(ignore_suffixes)


def _doc_change_satisfies(doc_changes: set[str], required_doc: str) -> bool:
    if required_doc.endswith("/"):
        return any(doc.startswith(required_doc) for doc in doc_changes)
    return required_doc in doc_changes


def _required_docs_for_file(path: str, rules: list[dict[str, object]]) -> list[str]:
    required: list[str] = []
    for rule in rules:
        code_prefix = str(rule["code_prefix"])
        if not path.startswith(code_prefix):
            continue
        docs = [item for item in rule["required_docs"] if isinstance(item, str)]
        required.extend(docs)

    # keep order, dedupe
    unique: list[str] = []
    seen: set[str] = set()
    for item in required:
        if item in seen:
            continue
        seen.add(item)
        unique.append(item)
    return unique


def validate_docs_sync(changed_files: list[str], rules_path: Path) -> tuple[bool, str]:
    rules, ignore_prefixes, ignore_suffixes = _load_rules(rules_path)

    if not changed_files:
        return True, "no changed files to evaluate"

    doc_changes = [item for item in changed_files if _is_doc_file(item)]
    non_doc_changes = [item for item in changed_files if not _is_doc_file(item)]

    behavior_changes = [
        item for item in non_doc_changes if not _is_ignored_non_doc(item, ignore_prefixes, ignore_suffixes)
    ]

    if not behavior_changes:
        return True, "no behavior/code changes to evaluate"

    doc_set = set(doc_changes)
    missing_mapped: list[tuple[str, list[str]]] = []

    for file_path in behavior_changes:
        required_docs = _required_docs_for_file(file_path, rules)
        if not required_docs:
            continue

        if any(_doc_change_satisfies(doc_set, item) for item in required_docs):
            continue

        missing_mapped.append((file_path, required_docs))

    if missing_mapped:
        snippets = []
        for file_path, required_docs in missing_mapped[:5]:
            snippets.append(f"{file_path} -> required docs: {', '.join(required_docs)}")
        return (
            False,
            "component docs sync required. " + " | ".join(snippets),
        )

    if not doc_changes:
        preview = ", ".join(behavior_changes[:5])
        return (
            False,
            "docs update required for behavior/code changes. "
            f"Changed code files (sample): {preview}. "
            "Update mapped docs (e.g. README.md).",
        )

    return True, "docs sync check passed"


def main() -> int:
    args = parse_args()
    changed_files = _load_changed_files(args)

    ok, message = validate_docs_sync(changed_files, Path(args.rules_file))
    if not ok:
        print(f"[DOCS-SYNC] FAIL: {message}", file=sys.stderr)
        return 1

    if not args.quiet:
        print(f"[DOCS-SYNC] PASS: {message}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
