#!/usr/bin/env bash
# State management utilities for Claude Octopus
# JSON-tool-free implementation for portability.

set -euo pipefail

if [[ -n "${CLAUDE_OCTOPUS_WORKSPACE:-}" && "${CLAUDE_OCTOPUS_WORKSPACE:0:1}" == "/" ]]; then
    STATE_DIR="${CLAUDE_OCTOPUS_WORKSPACE}"
else
    STATE_DIR="${HOME}/.claude-octopus"
fi
STATE_FILE="$STATE_DIR/state.json"
BACKUP_FILE="$STATE_DIR/state.json.backup"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

error() {
    echo -e "${RED}ERROR:${NC} $1" >&2
    exit 1
}

warning() {
    echo -e "${YELLOW}WARNING:${NC} $1" >&2
}

success() {
    echo -e "${GREEN}SUCCESS:${NC} $1"
}

json_valid_file() {
    local file="$1"
    python3 - "$file" <<'PY' >/dev/null 2>&1
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    json.load(f)
PY
}

json_get() {
    local path="$1"
    local default_value="${2:-null}"
    [[ -f "$STATE_FILE" ]] || { printf '%s\n' "$default_value"; return 0; }
    python3 - "$STATE_FILE" "$path" "$default_value" <<'PY'
import json, sys
fp, path, default = sys.argv[1], sys.argv[2], sys.argv[3]
try:
    with open(fp, "r", encoding="utf-8") as f:
        data = json.load(f)
    cur = data
    for part in path.split("."):
        if not part:
            continue
        if isinstance(cur, dict):
            if part not in cur:
                print(default); raise SystemExit(0)
            cur = cur[part]
        elif isinstance(cur, list) and part.isdigit():
            idx = int(part)
            if idx < 0 or idx >= len(cur):
                print(default); raise SystemExit(0)
            cur = cur[idx]
        else:
            print(default); raise SystemExit(0)
    if cur is None:
        print(default)
    elif isinstance(cur, bool):
        print("true" if cur else "false")
    elif isinstance(cur, (dict, list)):
        print(json.dumps(cur, ensure_ascii=False))
    else:
        print(str(cur))
except Exception:
    print(default)
PY
}

atomic_write() {
    local content="$1"
    local temp_file="${STATE_FILE}.tmp.$$"
    printf '%s\n' "$content" > "$temp_file"
    if ! json_valid_file "$temp_file"; then
        rm -f "$temp_file"
        error "Generated invalid JSON, write aborted"
    fi
    [[ -f "$STATE_FILE" ]] && cp "$STATE_FILE" "$BACKUP_FILE"
    mv "$temp_file" "$STATE_FILE"
}

init_state() {
    mkdir -p "$STATE_DIR" "$STATE_DIR/context" "$STATE_DIR/summaries" "$STATE_DIR/quick"

    if [[ -f "$STATE_FILE" ]]; then
        if json_valid_file "$STATE_FILE"; then
            success "State file already exists and is valid"
            return 0
        fi
        warning "State file is corrupted, backing up and recreating"
        mv "$STATE_FILE" "${STATE_FILE}.corrupt.$(date +%s)"
    fi

    local project_id
    project_id=$(git config --get remote.origin.url 2>/dev/null | md5sum | cut -d' ' -f1 || echo "$(basename "$PWD")" | md5sum | cut -d' ' -f1)

    cat > "$STATE_FILE" <<EOF
{
  "version": "1.0.0",
  "project_id": "$project_id",
  "current_workflow": null,
  "current_phase": null,
  "session_start": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "decisions": [],
  "blockers": [],
  "context": {
    "discover": null,
    "define": null,
    "develop": null,
    "deliver": null
  },
  "metrics": {
    "phases_completed": 0,
    "total_execution_time_minutes": 0,
    "provider_usage": {
      "codex": 0,
      "gemini": 0,
      "claude": 0
    }
  }
}
EOF

    success "Initialized state file at $STATE_FILE"
}

read_state() {
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    json_valid_file "$STATE_FILE" || error "State file is corrupted"
    cat "$STATE_FILE"
}

get_current_phase() {
    json_get "current_phase" "null"
}

get_current_workflow() {
    json_get "current_workflow" "null"
}

set_current_workflow() {
    local workflow="$1"
    local phase="${2:-null}"
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    local updated
    updated=$(python3 - "$STATE_FILE" "$workflow" "$phase" <<'PY'
import json, sys
fp, workflow, phase = sys.argv[1], sys.argv[2], sys.argv[3]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
data["current_workflow"] = workflow
data["current_phase"] = None if phase in ("", "null") else phase
print(json.dumps(data, ensure_ascii=False, indent=2))
PY
)
    atomic_write "$updated"
    success "Set current workflow to '$workflow', phase to '$phase'"
}

write_decision() {
    local phase="$1"
    local decision="$2"
    local rationale="$3"
    local commit="${4:-$(git rev-parse HEAD 2>/dev/null || echo 'none')}"
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    local updated
    updated=$(python3 - "$STATE_FILE" "$phase" "$decision" "$rationale" "$commit" "$(date -u +%Y-%m-%d)" <<'PY'
import json, sys
fp, phase, decision, rationale, commit, day = sys.argv[1:]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
data.setdefault("decisions", []).append({
    "phase": phase, "decision": decision, "rationale": rationale, "date": day, "commit": commit
})
print(json.dumps(data, ensure_ascii=False, indent=2))
PY
)
    atomic_write "$updated"
    success "Recorded decision for phase '$phase'"
}

write_blocker() {
    local description="$1"
    local phase="$2"
    local status="${3:-active}"
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    local updated
    updated=$(python3 - "$STATE_FILE" "$description" "$phase" "$status" "$(date -u +%Y-%m-%d)" <<'PY'
import json, sys
fp, description, phase, status, day = sys.argv[1:]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
data.setdefault("blockers", []).append({
    "description": description, "phase": phase, "status": status, "created": day
})
print(json.dumps(data, ensure_ascii=False, indent=2))
PY
)
    atomic_write "$updated"
    success "Recorded blocker for phase '$phase'"
}

update_blocker_status() {
    local description="$1"
    local new_status="$2"
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    local updated
    updated=$(python3 - "$STATE_FILE" "$description" "$new_status" <<'PY'
import json, sys
fp, description, new_status = sys.argv[1:]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
for b in data.get("blockers", []):
    if b.get("description") == description:
        b["status"] = new_status
print(json.dumps(data, ensure_ascii=False, indent=2))
PY
)
    atomic_write "$updated"
    success "Updated blocker status to '$new_status'"
}

update_context() {
    local phase="$1"
    local context="$2"
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    local updated
    updated=$(python3 - "$STATE_FILE" "$phase" "$context" <<'PY'
import json, sys
fp, phase, context = sys.argv[1:]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
data.setdefault("context", {})[phase] = context
print(json.dumps(data, ensure_ascii=False, indent=2))
PY
)
    atomic_write "$updated"
    success "Updated context for phase '$phase'"
}

get_context() {
    local phase="$1"
    json_get "context.${phase}" "null"
}

update_metrics() {
    local metric_type="$1"
    local value="$2"
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    local updated
    updated=$(python3 - "$STATE_FILE" "$metric_type" "$value" <<'PY'
import json, sys
fp, metric_type, value = sys.argv[1:]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
m = data.setdefault("metrics", {})
if metric_type == "phases_completed":
    m["phases_completed"] = int(m.get("phases_completed", 0)) + 1
elif metric_type == "execution_time":
    m["total_execution_time_minutes"] = float(m.get("total_execution_time_minutes", 0)) + float(value)
elif metric_type == "provider":
    pu = m.setdefault("provider_usage", {})
    pu[value] = int(pu.get(value, 0)) + 1
else:
    print("__ERROR__Unknown metric type", end="")
    raise SystemExit(3)
print(json.dumps(data, ensure_ascii=False, indent=2))
PY
) || error "Unknown metric type: $metric_type"
    atomic_write "$updated"
    success "Updated metric '$metric_type'"
}

get_decisions() {
    local phase="${1:-all}"
    [[ -f "$STATE_FILE" ]] || { echo "[]"; return; }
    if [[ "$phase" == "all" ]]; then
        json_get "decisions" "[]"
        return
    fi
    python3 - "$STATE_FILE" "$phase" <<'PY'
import json, sys
fp, phase = sys.argv[1], sys.argv[2]
with open(fp, "r", encoding="utf-8") as f:
    data = json.load(f)
for d in data.get("decisions", []):
    if d.get("phase") == phase:
        print(json.dumps(d, ensure_ascii=False))
PY
}

get_active_blockers() {
    [[ -f "$STATE_FILE" ]] || { echo "[]"; return; }
    python3 - "$STATE_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
for b in data.get("blockers", []):
    if b.get("status") == "active":
        print(json.dumps(b, ensure_ascii=False))
PY
}

show_summary() {
    [[ -f "$STATE_FILE" ]] || error "State file not found. Run 'init_state' first."
    echo "=== Claude Octopus State Summary ==="
    echo ""
    echo "Project ID: $(json_get "project_id" "")"
    echo "Session Start: $(json_get "session_start" "")"
    echo "Current Workflow: $(json_get "current_workflow" "none")"
    echo "Current Phase: $(json_get "current_phase" "none")"
    echo ""
    echo "Metrics:"
    echo "  Phases Completed: $(json_get "metrics.phases_completed" "0")"
    echo "  Execution Time: $(json_get "metrics.total_execution_time_minutes" "0") minutes"
    echo "  Provider Usage:"
    echo "    - Codex: $(json_get "metrics.provider_usage.codex" "0")"
    echo "    - Gemini: $(json_get "metrics.provider_usage.gemini" "0")"
    echo "    - Claude: $(json_get "metrics.provider_usage.claude" "0")"
    echo ""
    echo "Decisions: $(python3 - "$STATE_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    print(len(json.load(f).get("decisions", [])))
PY
)"
    echo "Active Blockers: $(python3 - "$STATE_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
print(sum(1 for b in data.get("blockers", []) if b.get("status") == "active"))
PY
)"
}

main() {
    local command="${1:-help}"
    shift || true
    case "$command" in
        init_state) init_state ;;
        read_state) read_state ;;
        get_current_phase) get_current_phase ;;
        get_current_workflow) get_current_workflow ;;
        set_current_workflow) set_current_workflow "$@" ;;
        write_decision) write_decision "$@" ;;
        write_blocker) write_blocker "$@" ;;
        update_blocker_status) update_blocker_status "$@" ;;
        update_context) update_context "$@" ;;
        get_context) get_context "$@" ;;
        update_metrics) update_metrics "$@" ;;
        get_decisions) get_decisions "$@" ;;
        get_active_blockers) get_active_blockers ;;
        show_summary) show_summary ;;
        help)
            cat <<EOF
Claude Octopus State Manager

Usage: state-manager.sh <command> [args]

Commands:
  init_state                                  Initialize state file
  read_state                                  Display current state (JSON)
  get_current_phase                           Get current phase
  get_current_workflow                        Get current workflow
  set_current_workflow <workflow> [phase]     Set current workflow and phase
  write_decision <phase> <decision> <rationale> [commit]
                                              Record a decision
  write_blocker <description> <phase> [status]
                                              Record a blocker
  update_blocker_status <description> <status>
                                              Update blocker status
  update_context <phase> <context>            Update phase context
  get_context <phase>                         Get phase context
  update_metrics <type> <value>               Update metrics
                                              Types: phases_completed, execution_time, provider
  get_decisions [phase]                       Get decisions (all or by phase)
  get_active_blockers                         Get active blockers
  show_summary                                Display state summary
  help                                        Show this help
EOF
            ;;
        *)
            error "Unknown command: $command. Run 'state-manager.sh help' for usage."
            ;;
    esac
}

if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi
