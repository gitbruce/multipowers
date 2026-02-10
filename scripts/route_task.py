#!/usr/bin/env python3
"""Deterministic task router for fast lane vs standard lane."""

from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import dataclass


STANDARD_KEYWORDS = {
    "architecture",
    "migrate",
    "migration",
    "refactor",
    "security",
    "review",
    "compliance",
    "workflow",
    "design",
    "database",
    "breaking",
    "governance",
    "major",
}

PLANNING_KEYWORDS = {
    "plan",
    "architecture",
    "design",
    "proposal",
    "spec",
}

REVIEW_KEYWORDS = {
    "review",
    "audit",
    "verify",
    "validation",
}


@dataclass(frozen=True)
class RouteDecision:
    lane: str
    reason: str
    suggested_workflow: str
    suggested_role: str

    def to_dict(self, task: str, risk_hint: str | None, force_lane: str | None) -> dict[str, str]:
        payload = {
            "task": task,
            "lane": self.lane,
            "reason": self.reason,
            "suggested_workflow": self.suggested_workflow,
            "suggested_role": self.suggested_role,
        }
        if risk_hint:
            payload["risk_hint"] = risk_hint
        if force_lane:
            payload["force_lane"] = force_lane
        return payload


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Route task to fast or standard lane")
    parser.add_argument("--task", required=True, help="Task description")
    parser.add_argument(
        "--risk-hint",
        choices=["low", "medium", "high", "critical"],
        help="Optional risk hint from caller",
    )
    parser.add_argument(
        "--force-lane",
        choices=["fast", "standard"],
        help="Force lane selection override",
    )
    parser.add_argument("--json", action="store_true", help="Emit JSON output")
    return parser.parse_args()


def _word_count(text: str) -> int:
    return len(re.findall(r"[A-Za-z0-9_]+", text))


def _suggest_standard_workflow(task_lower: str) -> tuple[str, str]:
    if any(keyword in task_lower for keyword in PLANNING_KEYWORDS):
        return "writing-plans", "architect"
    if any(keyword in task_lower for keyword in REVIEW_KEYWORDS):
        return "subagent-driven-development", "architect"
    return "subagent-driven-development", "coder"


def route_task(task: str, risk_hint: str | None, force_lane: str | None) -> RouteDecision:
    normalized_task = task.strip()
    if not normalized_task:
        raise ValueError("task must not be empty")

    task_lower = normalized_task.lower()

    if force_lane:
        if force_lane == "fast":
            return RouteDecision(
                lane="fast",
                reason="force_lane override requested",
                suggested_workflow="fast-execution",
                suggested_role="coder",
            )

        workflow, role = _suggest_standard_workflow(task_lower)
        return RouteDecision(
            lane="standard",
            reason="force_lane override requested",
            suggested_workflow=workflow,
            suggested_role=role,
        )

    score = 0
    reasons: list[str] = []

    if risk_hint in {"high", "critical"}:
        score += 3
        reasons.append(f"risk hint is {risk_hint}")
    elif risk_hint == "medium":
        score += 1
        reasons.append("risk hint is medium")

    matched_keywords = sorted(keyword for keyword in STANDARD_KEYWORDS if keyword in task_lower)
    if matched_keywords:
        score += 2
        reasons.append(f"contains standard keywords: {', '.join(matched_keywords)}")

    words = _word_count(task_lower)
    if words >= 28:
        score += 1
        reasons.append(f"task length is {words} words")

    if score >= 2:
        workflow, role = _suggest_standard_workflow(task_lower)
        return RouteDecision(
            lane="standard",
            reason="; ".join(reasons),
            suggested_workflow=workflow,
            suggested_role=role,
        )

    return RouteDecision(
        lane="fast",
        reason="bounded low-risk task",
        suggested_workflow="fast-execution",
        suggested_role="coder",
    )


def main() -> int:
    args = parse_args()

    try:
        decision = route_task(args.task, args.risk_hint, args.force_lane)
    except ValueError as exc:
        print(f"[ROUTE] {exc}", file=sys.stderr)
        return 1

    payload = decision.to_dict(args.task, args.risk_hint, args.force_lane)
    if args.json:
        print(json.dumps(payload, ensure_ascii=False, sort_keys=True))
        return 0

    print(f"lane={payload['lane']}")
    print(f"reason={payload['reason']}")
    print(f"suggested_workflow={payload['suggested_workflow']}")
    print(f"suggested_role={payload['suggested_role']}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
