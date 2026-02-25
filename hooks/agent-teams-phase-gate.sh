#!/bin/bash
# Claude Octopus Agent Teams Phase Gate Hook
# TaskCompleted hook that checks bridge ledger for phase completion.
set -euo pipefail

BRIDGE_DIR="${HOME}/.claude-octopus/bridge"
BRIDGE_LEDGER="${BRIDGE_DIR}/task-ledger.json"
hook_input=$(cat 2>/dev/null || true)

if [[ "${OCTOPUS_AGENT_TEAMS_BRIDGE:-auto}" == "disabled" ]]; then
    echo '{"decision": "continue"}'
    exit 0
fi
[[ -f "$BRIDGE_LEDGER" ]] || { echo '{"decision": "continue"}'; exit 0; }

task_id=""
if [[ -n "$hook_input" ]]; then
    task_id=$(python3 - <<'PY' "$hook_input"
import json, sys
try:
    print(json.loads(sys.argv[1]).get("task_id",""))
except Exception:
    print("")
PY
)
fi

mapfile -t _phase_meta < <(python3 - "$BRIDGE_LEDGER" "$task_id" <<'PY'
import datetime, json, sys
path, task_id = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    d = json.load(f)
current_phase = d.get("current_phase", "")
if task_id and task_id in d.get("tasks", {}):
    t = d["tasks"][task_id]
    t["status"] = "completed"
    t["completed_at"] = datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z"
    phase = t.get("phase", current_phase)
    if phase:
        p = d.setdefault("phases", {}).setdefault(phase, {})
        p["completed_tasks"] = int(p.get("completed_tasks", 0)) + 1
    with open(path + ".tmp", "w", encoding="utf-8") as out:
        json.dump(d, out, ensure_ascii=False, indent=2); out.write("\n")
total = int(d.get("phases", {}).get(current_phase, {}).get("total_tasks", 0))
completed = int(d.get("phases", {}).get(current_phase, {}).get("completed_tasks", 0))
threshold = float(d.get("phases", {}).get(current_phase, {}).get("gate", {}).get("threshold", 0.75))
print(current_phase)
print(total)
print(completed)
print(threshold)
PY
)
[[ -f "${BRIDGE_LEDGER}.tmp" ]] && mv "${BRIDGE_LEDGER}.tmp" "$BRIDGE_LEDGER"

current_phase="${_phase_meta[0]:-}"
total_tasks="${_phase_meta[1]:-0}"
completed_tasks="${_phase_meta[2]:-0}"
gate_threshold="${_phase_meta[3]:-0.75}"

if [[ -z "$current_phase" || "$total_tasks" -eq 0 ]]; then
    echo '{"decision": "continue"}'
    exit 0
fi

if [[ "$completed_tasks" -lt "$total_tasks" ]]; then
    remaining=$((total_tasks - completed_tasks))
    echo "{\"decision\": \"continue\", \"reason\": \"Phase $current_phase: $completed_tasks/$total_tasks tasks complete ($remaining remaining)\"}"
    exit 0
fi

completion_ratio=$(awk -v c="$completed_tasks" -v t="$total_tasks" 'BEGIN { printf "%.2f", c / t }')
gate_passed="false"
if awk -v r="$completion_ratio" -v t="$gate_threshold" 'BEGIN { exit !(r >= t) }'; then
    gate_passed="true"
fi

python3 - "$BRIDGE_LEDGER" "$current_phase" "$gate_passed" "$completion_ratio" <<'PY'
import json, sys
path, phase, passed, ratio = sys.argv[1:]
with open(path, "r", encoding="utf-8") as f:
    d=json.load(f)
gate = d.setdefault("phases", {}).setdefault(phase, {}).setdefault("gate", {})
gate["status"] = "evaluated"
gate["result"] = {"passed": passed == "true", "completion_ratio": float(ratio)}
with open(path + ".tmp", "w", encoding="utf-8") as out:
    json.dump(d, out, ensure_ascii=False, indent=2); out.write("\n")
PY
mv "${BRIDGE_LEDGER}.tmp" "$BRIDGE_LEDGER"

session_file="${HOME}/.claude-octopus/session.json"
if [[ -f "$session_file" ]]; then
    python3 - "$session_file" "$completed_tasks" <<'PY'
import json, sys
path, completed = sys.argv[1], int(sys.argv[2])
with open(path, "r", encoding="utf-8") as f:
    d=json.load(f)
d["phase_status"] = "completed"
d.setdefault("phase_tasks", {})["completed"] = completed
with open(path + ".tmp", "w", encoding="utf-8") as out:
    json.dump(d, out, ensure_ascii=False, indent=2); out.write("\n")
PY
    mv "${session_file}.tmp" "$session_file"
fi

if [[ "$gate_passed" == "true" ]]; then
    echo "{\"decision\": \"continue\", \"reason\": \"Phase $current_phase complete: quality gate passed (${completion_ratio} >= ${gate_threshold})\"}"
else
    echo "{\"decision\": \"block\", \"reason\": \"Phase $current_phase: quality gate failed (${completion_ratio} < ${gate_threshold}). $completed_tasks/$total_tasks tasks completed.\"}"
fi
