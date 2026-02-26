#!/usr/bin/env bash
# TeammateIdle Hook Handler - Claude Code v2.1.33+
# Dispatches queued work to idle agents during multi-agent workflows
# ═══════════════════════════════════════════════════════════════════════════════

set -euo pipefail

SESSION_FILE="${CLAUDE_OCTOPUS_WORKSPACE:-${PWD}/.multipowers/temp}/session.json"

# Only act if an active workflow session exists
if [[ ! -f "$SESSION_FILE" ]]; then
    exit 0
fi

CURRENT_PHASE=$(python3 - "$SESSION_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    d=json.load(f)
print(d.get("phase",""))
PY
)
if [[ -z "$CURRENT_PHASE" ]]; then
    exit 0
fi

QUEUE_LENGTH=$(python3 - "$SESSION_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    d=json.load(f)
print(len(d.get("agent_queue",[])))
PY
)

if [[ "$QUEUE_LENGTH" -gt 0 ]]; then
    # Dequeue next task
    mapfile -t _task_payload < <(python3 - "$SESSION_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    d=json.load(f)
q=d.get("agent_queue",[])
first=q[0] if q else {}
print(first.get("task","No task description"))
print(first.get("role","general"))
d["agent_queue"]=q[1:] if q else []
with open(sys.argv[1]+".tmp","w",encoding="utf-8") as f:
    json.dump(d,f,ensure_ascii=False,indent=2); f.write("\n")
PY
)
    NEXT_TASK="${_task_payload[0]:-No task description}"
    NEXT_ROLE="${_task_payload[1]:-general}"
    mv "${SESSION_FILE}.tmp" "$SESSION_FILE"

    # Track idle event in metrics
    METRICS_DIR="${CLAUDE_OCTOPUS_WORKSPACE:-${PWD}/.multipowers/temp}/metrics"
    mkdir -p "$METRICS_DIR"
    echo "{\"event\":\"teammate_idle\",\"phase\":\"$CURRENT_PHASE\",\"dispatched_task\":\"$NEXT_TASK\",\"dispatched_role\":\"$NEXT_ROLE\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" \
        >> "${METRICS_DIR}/idle-events.jsonl"

    # Output context for Claude to use
    echo "🐙 TeammateIdle: Dispatching queued task to idle agent"
    echo "Phase: $CURRENT_PHASE | Role: $NEXT_ROLE | Queue remaining: $((QUEUE_LENGTH - 1))"
else
    # No more work - check if phase should transition
    mapfile -t _phase_progress < <(python3 - "$SESSION_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    d=json.load(f)
pt=d.get("phase_tasks",{})
print(pt.get("completed",0))
print(pt.get("total",0))
PY
)
    COMPLETED="${_phase_progress[0]:-0}"
    TOTAL="${_phase_progress[1]:-0}"

    if [[ "$COMPLETED" -ge "$TOTAL" ]] && [[ "$TOTAL" -gt 0 ]]; then
        echo "🐙 TeammateIdle: All phase tasks complete. Ready for phase transition."
    fi
fi
