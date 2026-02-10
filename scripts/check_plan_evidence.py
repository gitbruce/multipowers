#!/usr/bin/env python3
"""Check that DONE tasks in plan docs are covered by structured evidence sections."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path
from typing import Iterable

TASK_HEADER_RE = re.compile(r"^###\s+(T[0-9A-Za-z-]+)")
STATUS_DONE_RE = re.compile(r"\*\*(Status|状态)\*\*\s*[：:]\s*`DONE`", re.IGNORECASE)
EVIDENCE_HEADING_RE = re.compile(r"^##\s+.*(证据|Evidence)")
COVERAGE_RE = re.compile(r"\*\*(Coverage Task IDs|覆盖任务ID)\*\*:\s*`([^`]+)`")
TASK_ID_IN_TEXT_RE = re.compile(r"T[0-9A-Za-z-]+")

REQUIRED_EVIDENCE_FIELDS = [
    "**Date**:",
    "**Verifier**:",
    "**Command(s)**:",
    "**Exit Code**:",
    "**Key Output**:",
]

GOVERNANCE_PROOF_PATTERNS = [
    re.compile(r"run_governance_checks\.sh"),
    re.compile(r"\bsemgrep\b", re.IGNORECASE),
    re.compile(r"\bbiome\b", re.IGNORECASE),
    re.compile(r"\bruff\b", re.IGNORECASE),
    re.compile(r"governance artifact", re.IGNORECASE),
    re.compile(r"outputs/governance", re.IGNORECASE),
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate DONE-task evidence coverage in plan markdown files")
    parser.add_argument("--require-governance-evidence", action="store_true", help="Require governance proof in evidence sections")
    parser.add_argument("files", nargs="*", help="Plan markdown files to validate")
    return parser.parse_args()


def default_plan_files() -> list[Path]:
    return sorted(Path("docs/plans").glob("gap_analysis_plan*.md"))


def iter_task_sections(lines: list[str]) -> Iterable[tuple[str, str]]:
    current_task_id: str | None = None
    current_lines: list[str] = []

    for line in lines:
        match = TASK_HEADER_RE.match(line)
        if match:
            if current_task_id is not None:
                yield current_task_id, "\n".join(current_lines)
            current_task_id = match.group(1)
            current_lines = [line]
            continue

        if current_task_id is not None:
            current_lines.append(line)

    if current_task_id is not None:
        yield current_task_id, "\n".join(current_lines)


def find_evidence_text(lines: list[str]) -> str:
    for index, line in enumerate(lines):
        if EVIDENCE_HEADING_RE.match(line.strip()):
            return "\n".join(lines[index:])
    return ""


def parse_coverage_ids(evidence_text: str) -> set[str]:
    coverage_ids: set[str] = set()
    for match in COVERAGE_RE.finditer(evidence_text):
        coverage_chunk = match.group(2)
        for task_id in TASK_ID_IN_TEXT_RE.findall(coverage_chunk):
            coverage_ids.add(task_id)
    return coverage_ids


def has_governance_proof(evidence_text: str) -> bool:
    return any(pattern.search(evidence_text) for pattern in GOVERNANCE_PROOF_PATTERNS)


def check_file(path: Path, require_governance_evidence: bool) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines()

    done_task_ids: list[str] = []
    for task_id, section in iter_task_sections(lines):
        if STATUS_DONE_RE.search(section):
            done_task_ids.append(task_id)

    if not done_task_ids:
        return errors

    evidence_text = find_evidence_text(lines)
    if not evidence_text:
        errors.append(f"{path}: DONE tasks exist but no evidence section found")
        return errors

    coverage_ids = parse_coverage_ids(evidence_text)
    if not coverage_ids:
        errors.append(f"{path}: evidence section missing '**Coverage Task IDs**' entries")
    else:
        missing_coverage = [task_id for task_id in done_task_ids if task_id not in coverage_ids]
        if missing_coverage:
            errors.append(
                f"{path}: DONE tasks missing coverage IDs: {', '.join(missing_coverage)}"
            )

    for required_field in REQUIRED_EVIDENCE_FIELDS:
        if required_field not in evidence_text:
            errors.append(f"{path}: evidence section missing required field '{required_field}'")

    if require_governance_evidence and not has_governance_proof(evidence_text):
        errors.append(
            f"{path}: governance evidence required but missing (expected governance command or artifact reference)"
        )

    return errors


def main() -> int:
    args = parse_args()

    if args.files:
        files = [Path(item) for item in args.files]
    else:
        files = default_plan_files()

    if not files:
        print("[PLAN-EVIDENCE] No plan files found", file=sys.stderr)
        return 1

    all_errors: list[str] = []
    for file_path in files:
        if not file_path.exists():
            all_errors.append(f"{file_path}: file not found")
            continue
        all_errors.extend(check_file(file_path, args.require_governance_evidence))

    if all_errors:
        for issue in all_errors:
            print(f"[PLAN-EVIDENCE] {issue}", file=sys.stderr)
        return 1

    print(f"[PLAN-EVIDENCE] PASS ({len(files)} files)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
