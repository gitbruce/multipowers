#!/usr/bin/env bash
# Agent Teams Bridge for Claude Octopus
# Unified task-ledger that bridges orchestrate.sh bash-spawned agents with Agent Teams state.

OCTOPUS_AGENT_TEAMS_BRIDGE="${OCTOPUS_AGENT_TEAMS_BRIDGE:-auto}"

_BRIDGE_DIR="${WORKSPACE_DIR:-${CLAUDE_OCTOPUS_WORKSPACE:-${PWD}/.multipowers/temp}}/bridge"
_BRIDGE_LEDGER="${_BRIDGE_DIR}/task-ledger.json"
_BRIDGE_LOCKFILE="${_BRIDGE_DIR}/.ledger.lock"

bridge_is_enabled() {
    case "$OCTOPUS_AGENT_TEAMS_BRIDGE" in
        enabled) return 0 ;;
        disabled) return 1 ;;
        auto) [[ "${SUPPORTS_AGENT_TEAMS_BRIDGE:-false}" == "true" ]] && return 0 || return 1 ;;
        *) return 1 ;;
    esac
}

bridge_ledger_update() {
    local action="$1"
    shift
    bridge_is_enabled || return 0
    [[ -f "$_BRIDGE_LEDGER" ]] || return 1
    python3 - "$_BRIDGE_LEDGER" "$action" "$@" <<'PY'
import datetime, json, sys
path, action, *args = sys.argv[1:]
with open(path, "r", encoding="utf-8") as f:
    d = json.load(f)
phases = d.setdefault("phases", {})
tasks = d.setdefault("tasks", {})
if action == "register_task":
    task_id, agent_type, phase, role = args
    tasks[task_id] = {
        "agent_type": agent_type,
        "phase": phase,
        "role": role,
        "status": "running",
        "registered_at": datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z",
        "completed_at": None,
    }
    p = phases.setdefault(phase, {})
    p["total_tasks"] = int(p.get("total_tasks", 0)) + 1
elif action == "complete_task":
    task_id, status = args
    t = tasks.get(task_id, {})
    t["status"] = status
    t["completed_at"] = datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z"
    tasks[task_id] = t
    phase = t.get("phase", "")
    if phase:
        p = phases.setdefault(phase, {})
        p["completed_tasks"] = int(p.get("completed_tasks", 0)) + 1
elif action == "inject_gate":
    phase, gate_type, threshold = args
    p = phases.setdefault(phase, {})
    p["gate"] = {
        "type": gate_type,
        "threshold": float(threshold),
        "status": "pending",
        "result": None,
    }
elif action == "set_gate_result":
    phase, passed, ratio = args
    p = phases.setdefault(phase, {})
    gate = p.setdefault("gate", {})
    gate["status"] = "evaluated"
    gate["result"] = {"passed": passed == "true", "completion_ratio": float(ratio)}
elif action == "enqueue_cross_provider":
    desc, source, target = args
    d.setdefault("cross_provider_queue", []).append({
        "description": desc,
        "source": source,
        "target": target,
        "queued_at": datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z",
        "status": "pending",
    })
elif action == "route_memory":
    phase, key, value = args
    warm = d.setdefault("memory", {}).setdefault("warm_start", {})
    warm[f"{phase}.{key}"] = value
elif action == "set_warm_start_file":
    phase, file = args
    warm = d.setdefault("memory", {}).setdefault("warm_start", {})
    warm[phase] = file
elif action == "set_phase_summary":
    phase, summary, synthesis_file = args
    ps = d.setdefault("memory", {}).setdefault("phase_summaries", {})
    ps[phase] = {
        "summary": summary,
        "synthesis_file": synthesis_file,
        "generated_at": datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z",
    }
elif action == "set_current_phase":
    phase = args[0]
    d["current_phase"] = phase
elif action == "complete_workflow":
    d["status"] = "completed"
    d["completed_at"] = datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z"
tmp = path + ".tmp"
with open(tmp, "w", encoding="utf-8") as f:
    json.dump(d, f, ensure_ascii=False, indent=2)
    f.write("\n")
PY
    mv "${_BRIDGE_LEDGER}.tmp" "$_BRIDGE_LEDGER"
}

bridge_init_ledger() {
    local workflow_name="${1:-embrace}"
    local task_group="${2:-$(date +%s)}"
    bridge_is_enabled || return 0
    mkdir -p "$_BRIDGE_DIR"
    python3 - "$_BRIDGE_LEDGER" "$workflow_name" "$task_group" <<'PY'
import datetime, json, sys
path, workflow, group = sys.argv[1:]
data = {
    "workflow": workflow,
    "task_group": group,
    "started_at": datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z",
    "status": "running",
    "current_phase": None,
    "phases": {},
    "tasks": {},
    "gate_results": {},
    "cross_provider_queue": [],
    "memory": {"warm_start": {}, "phase_summaries": {}},
}
with open(path, "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False, indent=2)
    f.write("\n")
PY
}

bridge_register_task() { bridge_ledger_update "register_task" "$1" "$2" "$3" "${4:-none}"; }
bridge_mark_task_complete() { bridge_ledger_update "complete_task" "$1" "${2:-completed}"; }

bridge_check_phase_complete() {
    local phase="$1"
    bridge_is_enabled || return 1
    [[ -f "$_BRIDGE_LEDGER" ]] || return 1
    python3 - "$_BRIDGE_LEDGER" "$phase" <<'PY'
import json, sys
path, phase = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    d = json.load(f)
p = d.get("phases", {}).get(phase, {})
total = int(p.get("total_tasks", 0))
done = int(p.get("completed_tasks", 0))
raise SystemExit(0 if total > 0 and done >= total else 1)
PY
}

bridge_inject_gate_task() { bridge_ledger_update "inject_gate" "$1" "${2:-quality}" "${3:-0.75}"; }

bridge_evaluate_gate() {
    local phase="$1"
    bridge_is_enabled || return 0
    [[ -f "$_BRIDGE_LEDGER" ]] || return 1
    local ratio_and_pass
    ratio_and_pass=$(python3 - "$_BRIDGE_LEDGER" "$phase" <<'PY'
import json, sys
path, phase = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    d = json.load(f)
p = d.get("phases", {}).get(phase, {})
g = p.get("gate", {})
threshold = float(g.get("threshold", 0.75))
total = int(p.get("total_tasks", 0))
done = int(p.get("completed_tasks", 0))
ratio = (done / total) if total > 0 else 0.0
passed = ratio >= threshold
print(f"{ratio:.2f}|{'true' if passed else 'false'}")
PY
)
    local ratio passed
    IFS='|' read -r ratio passed <<< "$ratio_and_pass"
    bridge_ledger_update "set_gate_result" "$phase" "$passed" "$ratio"
    [[ "$passed" == "true" ]]
}

bridge_get_idle_dispatch_target() {
    local preferred_provider="${1:-claude}"
    bridge_is_enabled || return 1
    [[ -f "$_BRIDGE_LEDGER" ]] || return 1
    python3 - "$_BRIDGE_LEDGER" "$preferred_provider" <<'PY'
import json, sys
path, preferred = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    d = json.load(f)
busy = set()
for t in d.get("tasks", {}).values():
    if t.get("status") == "running":
        busy.add(str(t.get("agent_type", "")).split("-")[0])
idle = [p for p in ("codex", "gemini", "claude") if p not in busy]
if preferred in idle:
    print(preferred)
elif idle:
    print(idle[0])
PY
}

bridge_enqueue_cross_provider_task() { bridge_ledger_update "enqueue_cross_provider" "$1" "$2" "${3:-}"; }
bridge_route_memory() { bridge_ledger_update "route_memory" "$1" "$2" "$3"; }

bridge_write_warm_start_memory() {
    local phase="$1"
    local content="$2"
    bridge_is_enabled || return 0
    local memory_file="${_BRIDGE_DIR}/warm-start-${phase}.md"
    mkdir -p "$_BRIDGE_DIR"
    printf '%s\n' "$content" > "$memory_file"
    bridge_ledger_update "set_warm_start_file" "$phase" "$memory_file"
}

bridge_generate_phase_summary() {
    local phase="$1"
    local synthesis_file="$2"
    bridge_is_enabled || return 0
    [[ -f "$synthesis_file" ]] || return 1
    local summary
    summary=$(head -c 2000 "$synthesis_file" 2>/dev/null || true)
    bridge_ledger_update "set_phase_summary" "$phase" "$summary" "$synthesis_file"
}

bridge_get_phase_status() {
    local phase="$1"
    bridge_is_enabled || return 1
    [[ -f "$_BRIDGE_LEDGER" ]] || return 1
    python3 - "$_BRIDGE_LEDGER" "$phase" <<'PY'
import json, sys
path, phase = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    d = json.load(f)
print(json.dumps(d.get("phases", {}).get(phase, {}), ensure_ascii=False))
PY
}

bridge_get_workflow_status() {
    bridge_is_enabled || return 1
    [[ -f "$_BRIDGE_LEDGER" ]] || return 1
    python3 - "$_BRIDGE_LEDGER" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    d = json.load(f)
tasks = d.get("tasks", {})
completed = sum(1 for t in tasks.values() if t.get("status") == "completed")
pending = sum(1 for q in d.get("cross_provider_queue", []) if q.get("status") == "pending")
out = {
    "workflow": d.get("workflow"),
    "status": d.get("status"),
    "current_phase": d.get("current_phase"),
    "total_tasks": len(tasks),
    "completed_tasks": completed,
    "pending_cross_provider": pending,
}
print(json.dumps(out, ensure_ascii=False))
PY
}

bridge_update_current_phase() { bridge_ledger_update "set_current_phase" "$1"; }
bridge_mark_workflow_complete() { bridge_ledger_update "complete_workflow"; }

bridge_cleanup() {
    bridge_is_enabled || return 0
    if [[ -f "$_BRIDGE_LEDGER" ]]; then
        local archive_dir="${_BRIDGE_DIR}/history"
        mkdir -p "$archive_dir"
        local ts
        ts=$(date +%Y%m%d-%H%M%S)
        cp "$_BRIDGE_LEDGER" "${archive_dir}/ledger-${ts}.json" 2>/dev/null || true
    fi
    rm -f "$_BRIDGE_LOCKFILE"
}
