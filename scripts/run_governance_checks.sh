#!/usr/bin/env bash
# Run major-change governance checks based on changed files.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

changed_files=()
changed_files_file=""
mode="strict"
artifact_path=""
request_id=""
track_id=""

usage() {
    cat <<'USAGE_EOF'
Usage: scripts/run_governance_checks.sh [options]

Options:
  --changed-file <path>       Add one changed file (repeatable)
  --changed-files-file <path> Newline-delimited changed file list
  --mode <strict|advisory>    Governance mode (default: strict)
  --strict-tools              Backward-compatible alias for --mode strict
  --artifact <path>           Optional artifact output path (JSON)
  --request-id <id>           Optional request id for traceability
  --track-id <id>             Optional track id for traceability
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
        --mode)
            mode="$2"
            shift 2
            ;;
        --strict-tools)
            mode="strict"
            shift
            ;;
        --artifact)
            artifact_path="$2"
            shift 2
            ;;
        --request-id)
            request_id="$2"
            shift 2
            ;;
        --track-id)
            track_id="$2"
            shift 2
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

if [[ "$mode" != "strict" && "$mode" != "advisory" ]]; then
    echo "[GOVERNANCE] Invalid mode: $mode (expected strict|advisory)" >&2
    exit 1
fi

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

if [[ ${#changed_files[@]} -eq 0 ]] && command -v git >/dev/null 2>&1; then
    while IFS= read -r line; do
        [[ -n "$line" ]] && changed_files+=("$line")
    done < <(git diff --name-only)
    while IFS= read -r line; do
        [[ -n "$line" ]] && changed_files+=("$line")
    done < <(git diff --name-only --cached)
fi

# unique keep order
if [[ ${#changed_files[@]} -gt 0 ]]; then
    deduped=()
    seen_keys=""
    for file in "${changed_files[@]}"; do
        [[ -z "$file" ]] && continue
        if [[ ",$seen_keys," == *",$file,"* ]]; then
            continue
        fi
        deduped+=("$file")
        seen_keys+="${file},"
    done
    changed_files=("${deduped[@]}")
fi

declare -a tool_results=()

tool_result_add() {
    local tool_name="$1"
    local status="$2"
    local exit_code="$3"
    local detail="$4"
    detail=${detail//|//}
    tool_results+=("${tool_name}|${status}|${exit_code}|${detail}")
}

emit_artifact_if_needed() {
    local overall_exit_code="$1"

    if [[ -z "$artifact_path" ]]; then
        return 0
    fi

    if [[ ! -f "scripts/record_governance_artifact.py" ]]; then
        echo "[GOVERNANCE] WARN: artifact script missing, cannot write artifact: scripts/record_governance_artifact.py" >&2
        return 0
    fi

    artifact_cmd=(
        python3 scripts/record_governance_artifact.py
        --artifact "$artifact_path"
        --mode "$mode"
        --overall-exit-code "$overall_exit_code"
    )

    [[ -n "$request_id" ]] && artifact_cmd+=(--request-id "$request_id")
    [[ -n "$track_id" ]] && artifact_cmd+=(--track-id "$track_id")

    for file in "${changed_files[@]}"; do
        artifact_cmd+=(--changed-file "$file")
    done
    for result in "${tool_results[@]}"; do
        artifact_cmd+=(--tool-result "$result")
    done

    "${artifact_cmd[@]}"
}

run_optional_tool() {
    local tool_name="$1"
    local fail_count_ref_name="$2"
    shift 2

    if ! command -v "$tool_name" >/dev/null 2>&1; then
        if [[ "$mode" == "strict" ]]; then
            echo "[GOVERNANCE] FAIL: required tool missing: $tool_name" >&2
            tool_result_add "$tool_name" "missing" 127 "required tool missing"
            eval "$fail_count_ref_name=$(( $fail_count_ref_name + 1 ))"
            return 0
        fi

        echo "[GOVERNANCE] WARN: $tool_name not installed; skipped in advisory mode" >&2
        tool_result_add "$tool_name" "missing" 127 "missing (advisory mode)"
        return 0
    fi

    set +e
    "$@"
    local status=$?
    set -e

    if [[ $status -eq 0 ]]; then
        tool_result_add "$tool_name" "passed" 0 "ok"
        return 0
    fi

    if [[ "$mode" == "strict" ]]; then
        echo "[GOVERNANCE] FAIL: $tool_name returned exit code $status" >&2
        tool_result_add "$tool_name" "failed" "$status" "tool execution failed"
        eval "$fail_count_ref_name=$(( $fail_count_ref_name + 1 ))"
        return 0
    fi

    echo "[GOVERNANCE] WARN: $tool_name failed with exit=$status in advisory mode" >&2
    tool_result_add "$tool_name" "failed" "$status" "tool failure ignored (advisory mode)"
}

if [[ ${#changed_files[@]} -eq 0 ]]; then
    echo "[GOVERNANCE] No changed files detected; skipping checks"
    tool_result_add "changed-files" "skipped" 0 "no changed files"
    emit_artifact_if_needed 0
    echo "[GOVERNANCE] PASS"
    exit 0
fi

echo "[GOVERNANCE] Mode: $mode"
echo "[GOVERNANCE] Changed files: ${#changed_files[@]}"

fail_count=0

# 1) docs sync is mandatory in both modes

docs_args=()
for file in "${changed_files[@]}"; do
    docs_args+=(--changed-file "$file")
done

set +e
python3 scripts/check_docs_sync.py --quiet "${docs_args[@]}"
docs_status=$?
set -e
if [[ $docs_status -ne 0 ]]; then
    tool_result_add "docs-sync" "failed" "$docs_status" "docs sync validation failed"
    fail_count=$((fail_count + 1))
else
    tool_result_add "docs-sync" "passed" 0 "docs sync ok"
fi

py_files=()
js_files=()
non_doc_files=()
for file in "${changed_files[@]}"; do
    case "$file" in
        *.md|*.txt|*.log)
            ;;
        *)
            non_doc_files+=("$file")
            ;;
    esac

    case "$file" in
        *.py)
            py_files+=("$file")
            ;;
        *.js|*.jsx|*.ts|*.tsx|*.mjs|*.cjs)
            js_files+=("$file")
            ;;
    esac
done

if [[ ${#non_doc_files[@]} -gt 0 ]]; then
    run_optional_tool semgrep fail_count semgrep --config auto --error "${non_doc_files[@]}"
else
    tool_result_add "semgrep" "skipped" 0 "no applicable files"
fi

if [[ ${#js_files[@]} -gt 0 ]]; then
    run_optional_tool biome fail_count biome check "${js_files[@]}"
else
    tool_result_add "biome" "skipped" 0 "no JS/TS files"
fi

if [[ ${#py_files[@]} -gt 0 ]]; then
    run_optional_tool ruff fail_count ruff check "${py_files[@]}"
else
    tool_result_add "ruff" "skipped" 0 "no Python files"
fi

overall_exit_code=0
if [[ $fail_count -gt 0 ]]; then
    overall_exit_code=1
fi

emit_artifact_if_needed "$overall_exit_code"

if [[ $overall_exit_code -eq 0 ]]; then
    echo "[GOVERNANCE] PASS"
else
    echo "[GOVERNANCE] FAIL"
fi

exit "$overall_exit_code"
