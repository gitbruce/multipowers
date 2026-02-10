#!/usr/bin/env python3
"""Validate Multipowers role configuration against roles schema."""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any


def _load_json(path: Path) -> Any:
    try:
        with path.open("r", encoding="utf-8") as handle:
            return json.load(handle)
    except FileNotFoundError:
        print(f"[ROLE-VALIDATION] Config file not found: {path}", file=sys.stderr)
        raise SystemExit(1)
    except json.JSONDecodeError as exc:
        print(f"[ROLE-VALIDATION] Invalid JSON in {path}: {exc}", file=sys.stderr)
        raise SystemExit(1)


def _type_name(value: Any) -> str:
    if isinstance(value, bool):
        return "boolean"
    if isinstance(value, int):
        return "integer"
    if isinstance(value, float):
        return "number"
    if isinstance(value, str):
        return "string"
    if isinstance(value, list):
        return "array"
    if isinstance(value, dict):
        return "object"
    return type(value).__name__


def _matches_type(value: Any, schema_type: str) -> bool:
    if schema_type == "object":
        return isinstance(value, dict)
    if schema_type == "array":
        return isinstance(value, list)
    if schema_type == "string":
        return isinstance(value, str)
    if schema_type == "number":
        return isinstance(value, (int, float)) and not isinstance(value, bool)
    if schema_type == "integer":
        return isinstance(value, int) and not isinstance(value, bool)
    if schema_type == "boolean":
        return isinstance(value, bool)
    return True


def _validate_role_map(config: dict[str, Any], schema: dict[str, Any], role_filter: str | None) -> list[str]:
    errors: list[str] = []

    roles = config.get("roles")
    if not isinstance(roles, dict):
        errors.append("'roles' must be an object")
        return errors

    roles_schema = schema.get("properties", {}).get("roles", {})
    pattern_props = roles_schema.get("patternProperties", {})
    pattern_str = next(iter(pattern_props.keys()), r"^[a-z-]+$")
    role_schema = next(iter(pattern_props.values()), {}) if pattern_props else {}

    pattern = re.compile(pattern_str)

    if role_filter and role_filter not in roles:
        errors.append(f"Role '{role_filter}' not found. Available roles: {', '.join(sorted(roles.keys()))}")
        return errors

    target_roles = {role_filter: roles[role_filter]} if role_filter else roles

    required_fields = role_schema.get("required", [])
    properties = role_schema.get("properties", {})
    allow_extra = role_schema.get("additionalProperties", True)

    for role_name, role_value in target_roles.items():
        if not pattern.match(role_name):
            errors.append(f"roles.{role_name}: invalid role key (must match /{pattern_str}/)")
            continue

        if not isinstance(role_value, dict):
            errors.append(f"roles.{role_name}: expected object, got {_type_name(role_value)}")
            continue

        for field in required_fields:
            if field not in role_value:
                errors.append(f"roles.{role_name}.{field}: missing required field")

        if not allow_extra:
            for field in role_value:
                if field not in properties:
                    errors.append(f"roles.{role_name}.{field}: unknown field")

        for field, value in role_value.items():
            field_schema = properties.get(field)
            if not field_schema:
                continue

            expected_type = field_schema.get("type")
            if expected_type and not _matches_type(value, expected_type):
                errors.append(
                    f"roles.{role_name}.{field}: expected {expected_type}, got {_type_name(value)}"
                )
                continue

            enum_values = field_schema.get("enum")
            if enum_values and value not in enum_values:
                errors.append(
                    f"roles.{role_name}.{field}: invalid value '{value}', expected one of {enum_values}"
                )

            if expected_type == "array":
                item_schema = field_schema.get("items", {})
                item_type = item_schema.get("type")
                min_len = item_schema.get("minLength")
                if isinstance(value, list):
                    for index, item in enumerate(value):
                        if item_type and not _matches_type(item, item_type):
                            errors.append(
                                f"roles.{role_name}.{field}[{index}]: expected {item_type}, got {_type_name(item)}"
                            )
                        if min_len is not None and isinstance(item, str) and len(item) < min_len:
                            errors.append(
                                f"roles.{role_name}.{field}[{index}]: must have minLength {min_len}"
                            )

            if expected_type == "string":
                min_len = field_schema.get("minLength")
                if min_len is not None and isinstance(value, str) and len(value) < min_len:
                    errors.append(f"roles.{role_name}.{field}: must have minLength {min_len}")

            if expected_type == "number":
                minimum = field_schema.get("minimum")
                maximum = field_schema.get("maximum")
                if minimum is not None and value < minimum:
                    errors.append(f"roles.{role_name}.{field}: must be >= {minimum}")
                if maximum is not None and value > maximum:
                    errors.append(f"roles.{role_name}.{field}: must be <= {maximum}")

    return errors


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate Multipowers roles config")
    parser.add_argument("--config", required=True, help="Path to roles config JSON")
    parser.add_argument("--schema", required=True, help="Path to roles schema JSON")
    parser.add_argument("--role", help="Optional role name to validate existence and structure")
    parser.add_argument("--quiet", action="store_true", help="Suppress success output")
    return parser.parse_args()


def main() -> int:
    args = parse_args()

    config_path = Path(args.config)
    schema_path = Path(args.schema)

    config = _load_json(config_path)
    schema = _load_json(schema_path)

    errors = _validate_role_map(config, schema, args.role)

    if errors:
        for issue in errors:
            print(f"[ROLE-VALIDATION] {issue}", file=sys.stderr)
        return 1

    if not args.quiet:
        print(f"[ROLE-VALIDATION] OK: {config_path}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
