#!/usr/bin/env bash
# TaskCompleted Hook Handler - Claude Code v2.1.33+
# Manages phase transitions when workflow tasks complete
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

# Increment completed task count
mapfile -t _session_vals < <(python3 - "$SESSION_FILE" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    d=json.load(f)
pt=d.get("phase_tasks",{})
print(pt.get("completed",0))
print(pt.get("total",0))
print(d.get("autonomy","supervised"))
PY
)
COMPLETED="${_session_vals[0]:-0}"
TOTAL="${_session_vals[1]:-0}"
AUTONOMY="${_session_vals[2]:-supervised}"
COMPLETED=$((COMPLETED + 1))

python3 - "$SESSION_FILE" "$COMPLETED" <<'PY'
import json, sys
path, completed = sys.argv[1], int(sys.argv[2])
with open(path, "r", encoding="utf-8") as f:
    d=json.load(f)
d.setdefault("phase_tasks", {})["completed"] = completed
with open(path + ".tmp", "w", encoding="utf-8") as f:
    json.dump(d, f, ensure_ascii=False, indent=2); f.write("\n")
PY
mv "${SESSION_FILE}.tmp" "$SESSION_FILE"

# Record metrics
METRICS_DIR="${CLAUDE_OCTOPUS_WORKSPACE:-${PWD}/.multipowers/temp}/metrics"
mkdir -p "$METRICS_DIR"
echo "{\"event\":\"task_completed\",\"phase\":\"$CURRENT_PHASE\",\"completed\":$COMPLETED,\"total\":$TOTAL,\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" \
    >> "${METRICS_DIR}/completion-events.jsonl"

# Check if phase is complete
if [[ "$COMPLETED" -ge "$TOTAL" ]] && [[ "$TOTAL" -gt 0 ]]; then
    # Determine next phase
    case "$CURRENT_PHASE" in
        probe)  NEXT_PHASE="grasp" ;;
        grasp)  NEXT_PHASE="tangle" ;;
        tangle) NEXT_PHASE="ink" ;;
        ink)    NEXT_PHASE="complete" ;;
        *)      NEXT_PHASE="unknown" ;;
    esac

    # Record phase completion
    RESULTS_DIR="${CLAUDE_OCTOPUS_WORKSPACE:-${PWD}/.multipowers/temp}/results"
    TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    echo "{\"phase\":\"$CURRENT_PHASE\",\"completed_at\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"tasks_completed\":$COMPLETED,\"next_phase\":\"$NEXT_PHASE\"}" \
        > "${RESULTS_DIR}/${CURRENT_PHASE}-complete-${TIMESTAMP}.json" 2>/dev/null || true

    if [[ "$NEXT_PHASE" == "complete" ]]; then
        echo "🐙 TaskCompleted: All workflow phases complete! ✅"
        python3 - "$SESSION_FILE" <<'PY'
import json, sys
path = sys.argv[1]
with open(path, "r", encoding="utf-8") as f:
    d=json.load(f)
d["phase"] = "complete"
d["workflow_status"] = "finished"
with open(path + ".tmp", "w", encoding="utf-8") as f:
    json.dump(d, f, ensure_ascii=False, indent=2); f.write("\n")
PY
        mv "${SESSION_FILE}.tmp" "$SESSION_FILE"
    else
        case "$AUTONOMY" in
            autonomous|semi-autonomous)
                # Auto-transition to next phase
                python3 - "$SESSION_FILE" "$NEXT_PHASE" <<'PY'
import json, sys
path, next_phase = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    d=json.load(f)
d["phase"] = next_phase
d["phase_tasks"] = {"total": 0, "completed": 0}
with open(path + ".tmp", "w", encoding="utf-8") as f:
    json.dump(d, f, ensure_ascii=False, indent=2); f.write("\n")
PY
                mv "${SESSION_FILE}.tmp" "$SESSION_FILE"
                echo "🐙 TaskCompleted: Phase '$CURRENT_PHASE' complete → Auto-transitioning to '$NEXT_PHASE'"
                ;;
            *)
                # Supervised mode - signal but don't auto-transition
                echo "🐙 TaskCompleted: Phase '$CURRENT_PHASE' complete ($COMPLETED/$TOTAL tasks)."
                echo "   Next phase: '$NEXT_PHASE' — awaiting user approval to proceed."
                ;;
        esac
    fi
else
    PROGRESS_PCT=$((COMPLETED * 100 / TOTAL))
    echo "🐙 TaskCompleted: Phase '$CURRENT_PHASE' progress: $COMPLETED/$TOTAL ($PROGRESS_PCT%)"
fi
