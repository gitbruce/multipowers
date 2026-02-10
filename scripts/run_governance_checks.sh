#!/usr/bin/env bash
# Run major-change governance checks based on changed files.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

changed_files=()
changed_files_file=""
allow_missing_tools="true"

usage() {
    cat <<'USAGE_EOF'
Usage: scripts/run_governance_checks.sh [options]

Options:
  --changed-file <path>       Add one changed file (repeatable)
  --changed-files-file <path> Newline-delimited changed file list
  --strict-tools              Fail if semgrep/biome/ruff are missing or fail to execute
  --help                      Show this help
USAGE_EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --changed-file)
            changed_files+=("$2")
            shift 2
            ;;
        --changed-files-file)
            changed_files_file="$2"
            shift 2
            ;;
        --strict-tools)
            allow_missing_tools="false"
            shift
            ;;
        --help|-h)
            usage
            exit 0
            ;;
        *)
            echo "[GOVERNANCE] Unknown option: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

if [[ -n "$changed_files_file" ]]; then
    if [[ ! -f "$changed_files_file" ]]; then
        echo "[GOVERNANCE] changed file list not found: $changed_files_file" >&2
        exit 1
    fi
    while IFS= read -r line; do
        line="${line#${line%%[![:space:]]*}}"
        line="${line%${line##*[![:space:]]}}"
        if [[ -n "$line" ]]; then
            changed_files+=("$line")
        fi
    done < "$changed_files_file"
fi

if [[ ${#changed_files[@]} -eq 0 ]]; then
    if command -v git >/dev/null 2>&1; then
        while IFS= read -r line; do
            [[ -n "$line" ]] && changed_files+=("$line")
        done < <(git diff --name-only)
    fi
fi

if [[ ${#changed_files[@]} -eq 0 ]]; then
    echo "[GOVERNANCE] No changed files detected; skipping checks"
    exit 0
fi

echo "[GOVERNANCE] Changed files: ${#changed_files[@]}"

docs_args=()
for file in "${changed_files[@]}"; do
    docs_args+=(--changed-file "$file")
done
python3 scripts/check_docs_sync.py --quiet "${docs_args[@]}"

py_files=()
js_files=()
for file in "${changed_files[@]}"; do
    case "$file" in
        *.py)
            py_files+=("$file")
            ;;
        *.js|*.jsx|*.ts|*.tsx|*.mjs|*.cjs)
            js_files+=("$file")
            ;;
    esac
done

run_or_warn_missing() {
    local tool_name="$1"
    shift

    if ! command -v "$tool_name" >/dev/null 2>&1; then
        if [[ "$allow_missing_tools" == "true" ]]; then
            echo "[GOVERNANCE] WARN: $tool_name is not installed; skipping" >&2
            return 0
        fi

        echo "[GOVERNANCE] FAIL: required tool missing: $tool_name" >&2
        return 1
    fi

    echo "[GOVERNANCE] Running: $tool_name"
    set +e
    "$@"
    local status=$?
    set -e

    if [[ $status -eq 0 ]]; then
        return 0
    fi

    if [[ "$allow_missing_tools" == "true" ]]; then
        echo "[GOVERNANCE] WARN: $tool_name failed to run cleanly (exit=$status); skipping in non-strict mode" >&2
        return 0
    fi

    echo "[GOVERNANCE] FAIL: $tool_name returned exit code $status" >&2
    return "$status"
}

# Semgrep runs when any non-doc file changes.
non_doc_files=()
for file in "${changed_files[@]}"; do
    if [[ "$file" == *.md ]] || [[ "$file" == *.txt ]] || [[ "$file" == *.log ]]; then
        continue
    fi
    non_doc_files+=("$file")
done

if [[ ${#non_doc_files[@]} -gt 0 ]]; then
    run_or_warn_missing semgrep semgrep --config auto --error "${non_doc_files[@]}"
fi

if [[ ${#js_files[@]} -gt 0 ]]; then
    run_or_warn_missing biome biome check "${js_files[@]}"
fi

if [[ ${#py_files[@]} -gt 0 ]]; then
    run_or_warn_missing ruff ruff check "${py_files[@]}"
fi

echo "[GOVERNANCE] PASS"
