#!/usr/bin/env python3
"""Update task status fields in a plan markdown file.

Updates both:
1) task section lines: **Status** / **状态**
2) task board table row: | T?-??? | `STATUS` | ...
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


ALLOWED_STATUS = {"TODO", "IN_PROGRESS", "BLOCKED", "DONE"}
TASK_HEADER_PATTERN = re.compile(r"^###\s+(T[0-9A-Za-z-]+)")
EN_STATUS_PATTERN = re.compile(r"^(\s*-\s*\*\*Status\*\*:\s*)`[^`]*`(\s*)$")
ZH_STATUS_PATTERN = re.compile(r"^(\s*-\s*\*\*状态\*\*[：:]\s*)`[^`]*`(\s*)$")
TABLE_ROW_PATTERN = re.compile(r"^(\|\s*)(T[0-9A-Za-z-]+)(\s*\|\s*)`([^`]*)`(\s*\|.*)$")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Update Status/状态 fields for a plan task")
    parser.add_argument("--file", required=True, help="Plan markdown file")
    parser.add_argument("--task-id", required=True, help="Task identifier, e.g. T6-001")
    parser.add_argument("--status", required=True, choices=sorted(ALLOWED_STATUS), help="Target status")
    return parser.parse_args()


def find_task_bounds(lines: list[str], task_id: str) -> tuple[int, int] | None:
    start = -1
    for idx, line in enumerate(lines):
        match = TASK_HEADER_PATTERN.match(line)
        if not match:
            continue
        if match.group(1) == task_id:
            start = idx
            break

    if start < 0:
        return None

    end = len(lines)
    for idx in range(start + 1, len(lines)):
        if lines[idx].startswith("### "):
            end = idx
            break

    return start, end


def update_task_section(lines: list[str], start: int, end: int, status: str) -> tuple[list[str], bool]:
    updated = list(lines)
    status_line = f"- **Status**: `{status}`"
    status_zh_line = f"- **状态**：`{status}`"

    found_en = False
    found_zh = False

    for idx in range(start, end):
        line = updated[idx]
        en_match = EN_STATUS_PATTERN.match(line)
        if en_match:
            updated[idx] = f"{en_match.group(1)}`{status}`{en_match.group(2)}"
            found_en = True
            continue

        zh_match = ZH_STATUS_PATTERN.match(line)
        if zh_match:
            updated[idx] = f"{zh_match.group(1)}`{status}`{zh_match.group(2)}"
            found_zh = True

    insert_at = start + 1
    while insert_at < end and updated[insert_at].strip() == "":
        insert_at += 1

    if not found_en:
        updated.insert(insert_at, status_line)
        insert_at += 1
        end += 1
    if not found_zh:
        updated.insert(insert_at, status_zh_line)

    changed = updated != lines
    return updated, changed


def update_task_board(lines: list[str], task_id: str, status: str) -> tuple[list[str], bool]:
    updated = list(lines)
    changed = False

    for idx, line in enumerate(updated):
        match = TABLE_ROW_PATTERN.match(line)
        if not match:
            continue

        row_task_id = match.group(2)
        if row_task_id != task_id:
            continue

        new_line = f"{match.group(1)}{row_task_id}{match.group(3)}`{status}`{match.group(5)}"
        if new_line != line:
            updated[idx] = new_line
            changed = True

    return updated, changed


def main() -> int:
    args = parse_args()
    plan_file = Path(args.file)

    if not plan_file.exists():
        print(f"[PLAN-STATUS] file not found: {plan_file}", file=sys.stderr)
        return 1

    text = plan_file.read_text(encoding="utf-8")
    lines = text.splitlines()

    bounds = find_task_bounds(lines, args.task_id)
    if bounds is None:
        print(f"[PLAN-STATUS] task not found: {args.task_id}", file=sys.stderr)
        return 1

    start, end = bounds
    updated_lines, changed_section = update_task_section(lines, start, end, args.status)
    updated_lines, changed_table = update_task_board(updated_lines, args.task_id, args.status)

    if changed_section or changed_table:
        plan_file.write_text("\n".join(updated_lines) + "\n", encoding="utf-8")

    print(f"[PLAN-STATUS] {args.task_id} => {args.status}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
