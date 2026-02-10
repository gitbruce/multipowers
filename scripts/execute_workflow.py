#!/usr/bin/env python3
"""Execute configured workflows with node-level role dispatch."""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any


@dataclass
class WorkflowNode:
    node_id: str
    role: str
    prompt: str


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run a multipowers workflow")
    parser.add_argument("--config", required=True, help="Workflow config JSON path")
    parser.add_argument("--workflow", required=True, help="Workflow name")
    parser.add_argument("--task", required=True, help="Task input")
    parser.add_argument("--ask-role", default="bin/ask-role", help="Path to ask-role executable")
    parser.add_argument(
        "--roles-config",
        default="",
        help="Effective roles config path (defaults: conductor/config/roles.json, then config/roles.default.json)",
    )
    parser.add_argument("--request-id", default="", help="Request identifier")
    parser.add_argument("--track-id", default="", help="Track identifier")
    parser.add_argument("--json", action="store_true", help="Emit JSON summary")
    parser.add_argument("--dry-run", action="store_true", help="Resolve nodes without execution")
    return parser.parse_args()


def _load_json(path: Path) -> dict[str, Any]:
    try:
        with path.open("r", encoding="utf-8") as handle:
            loaded = json.load(handle)
    except FileNotFoundError:
        raise ValueError(f"Workflow config not found: {path}") from None
    except json.JSONDecodeError as exc:
        raise ValueError(f"Invalid workflow JSON in {path}: {exc}") from None

    if not isinstance(loaded, dict):
        raise ValueError(f"Workflow config must be a JSON object: {path}")
    return loaded


def _resolve_roles_config_path(explicit_path: str) -> Path:
    if explicit_path:
        return Path(explicit_path)

    preferred = Path("conductor/config/roles.json")
    if preferred.exists():
        return preferred

    fallback = Path("config/roles.default.json")
    return fallback


def _load_available_roles(roles_config_path: Path) -> set[str]:
    try:
        payload = _load_json(roles_config_path)
    except ValueError as exc:
        raise ValueError(f"Invalid roles config: {exc}") from None

    roles_obj = payload.get("roles")
    if not isinstance(roles_obj, dict) or not roles_obj:
        raise ValueError(f"Roles config missing non-empty object field 'roles': {roles_config_path}")

    roles = {role for role in roles_obj.keys() if isinstance(role, str) and role.strip()}
    if not roles:
        raise ValueError(f"Roles config contains no valid role names: {roles_config_path}")

    return roles


def _render_prompt(template: str, task: str, workflow_name: str, node_id: str, role: str) -> str:
    rendered = template
    rendered = rendered.replace("{task}", task)
    rendered = rendered.replace("{workflow}", workflow_name)
    rendered = rendered.replace("{node}", node_id)
    rendered = rendered.replace("{role}", role)
    return rendered


def _validate_workflow(
    config: dict[str, Any],
    workflow_name: str,
    task: str,
    available_roles: set[str],
) -> tuple[str, list[WorkflowNode]]:
    workflows = config.get("workflows")
    if not isinstance(workflows, dict):
        raise ValueError("Workflow config missing object field: workflows")

    workflow = workflows.get(workflow_name)
    if workflow is None:
        available = ", ".join(sorted(workflows.keys()))
        raise ValueError(f"Workflow '{workflow_name}' not found. Available workflows: {available}")
    if not isinstance(workflow, dict):
        raise ValueError(f"Workflow '{workflow_name}' must be an object")

    default_role = workflow.get("default_role")
    if not isinstance(default_role, str) or not default_role.strip():
        raise ValueError(f"Workflow '{workflow_name}' missing non-empty default_role")
    if default_role not in available_roles:
        raise ValueError(
            f"Workflow '{workflow_name}' unknown role '{default_role}' (node=default_role)"
        )

    nodes = workflow.get("nodes")
    if not isinstance(nodes, list) or not nodes:
        raise ValueError(f"Workflow '{workflow_name}' must define a non-empty nodes array")

    compiled_nodes: list[WorkflowNode] = []
    for index, raw_node in enumerate(nodes, start=1):
        if not isinstance(raw_node, dict):
            raise ValueError(f"Workflow '{workflow_name}' node[{index}] must be an object")

        node_id = raw_node.get("id")
        if not isinstance(node_id, str) or not node_id.strip():
            raise ValueError(f"Workflow '{workflow_name}' node[{index}] missing non-empty id")

        node_role = raw_node.get("role", default_role)
        if not isinstance(node_role, str) or not node_role.strip():
            raise ValueError(
                f"Workflow '{workflow_name}' node '{node_id}' has invalid role override"
            )
        if node_role not in available_roles:
            raise ValueError(
                f"Workflow '{workflow_name}' unknown role '{node_role}' (node='{node_id}')"
            )

        prompt_template = raw_node.get("prompt_template")
        if not isinstance(prompt_template, str) or not prompt_template.strip():
            raise ValueError(
                f"Workflow '{workflow_name}' node '{node_id}' missing non-empty prompt_template"
            )

        compiled_nodes.append(
            WorkflowNode(
                node_id=node_id,
                role=node_role,
                prompt=_render_prompt(prompt_template, task, workflow_name, node_id, node_role),
            )
        )

    return default_role, compiled_nodes


def _emit_event(
    event: str,
    role: str,
    request_id: str,
    track_id: str,
    workflow: str,
    node: str = "",
    reason: str = "",
    metadata: dict[str, Any] | None = None,
):
    utils_path = Path("connectors/utils.py")
    if not utils_path.exists():
        return

    cmd = [
        sys.executable,
        str(utils_path),
        "event-log",
        "--event",
        event,
        "--role",
        role,
        "--workflow",
        workflow,
    ]
    if request_id:
        cmd.extend(["--request-id", request_id])
    if track_id:
        cmd.extend(["--track-id", track_id])
    if node:
        cmd.extend(["--node", node])
    if reason:
        cmd.extend(["--reason", reason])
    if metadata:
        cmd.extend(["--metadata-json", json.dumps(metadata, ensure_ascii=False, sort_keys=True)])

    subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=False)


def _run_node(
    ask_role_path: str,
    node: WorkflowNode,
    request_id: str,
) -> tuple[int, str, str]:
    env = os.environ.copy()
    if request_id:
        env["MULTIPOWERS_REQUEST_ID"] = request_id

    result = subprocess.run(
        [ask_role_path, node.role, node.prompt],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    return result.returncode, result.stdout, result.stderr


def main() -> int:
    args = parse_args()

    try:
        config = _load_json(Path(args.config))
        roles_config_path = _resolve_roles_config_path(args.roles_config)
        available_roles = _load_available_roles(roles_config_path)
        _, nodes = _validate_workflow(config, args.workflow, args.task, available_roles)
    except ValueError as exc:
        print(f"[WORKFLOW] {exc}", file=sys.stderr)
        return 1

    ask_role_path = args.ask_role
    if not Path(ask_role_path).exists():
        print(f"[WORKFLOW] ask-role executable not found: {ask_role_path}", file=sys.stderr)
        return 1

    summary: dict[str, Any] = {
        "workflow": args.workflow,
        "request_id": args.request_id,
        "track_id": args.track_id,
        "dry_run": args.dry_run,
        "nodes": [],
    }

    for index, node in enumerate(nodes, start=1):
        node_payload: dict[str, Any] = {
            "index": index,
            "id": node.node_id,
            "role": node.role,
        }

        if args.dry_run:
            node_payload["status"] = "skipped"
            summary["nodes"].append(node_payload)
            continue

        exit_code, stdout_text, stderr_text = _run_node(ask_role_path, node, args.request_id)

        node_payload["status"] = "ok" if exit_code == 0 else "failed"
        node_payload["exit_code"] = exit_code
        if stdout_text.strip():
            node_payload["stdout"] = stdout_text.strip()
        if stderr_text.strip():
            node_payload["stderr"] = "\n".join(stderr_text.strip().splitlines()[:8])
        summary["nodes"].append(node_payload)

        _emit_event(
            event="workflow_node_executed",
            role=node.role,
            request_id=args.request_id,
            track_id=args.track_id,
            workflow=args.workflow,
            node=node.node_id,
            reason="node completed" if exit_code == 0 else "node failed",
            metadata={"index": index, "exit_code": exit_code},
        )

        if exit_code != 0:
            if args.json:
                print(json.dumps(summary, ensure_ascii=False, sort_keys=True))
            else:
                print(
                    f"[WORKFLOW] Node '{node.node_id}' failed (role={node.role}, exit_code={exit_code})",
                    file=sys.stderr,
                )
                if stderr_text.strip():
                    print(stderr_text, file=sys.stderr)
            return exit_code

    if args.json:
        print(json.dumps(summary, ensure_ascii=False, sort_keys=True))
        return 0

    print(f"Workflow: {args.workflow}")
    for node in summary["nodes"]:
        suffix = f"exit_code={node.get('exit_code', 0)}"
        print(f"- [{node['status']}] {node['id']} ({node['role']}) {suffix}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
