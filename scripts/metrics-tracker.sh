#!/usr/bin/env bash
# Metrics Tracker for Claude Octopus
# Tracks resource usage (tokens, duration, costs) for multi-AI operations.

get_metrics_base() {
    echo "${METRICS_BASE:-${WORKSPACE_DIR:-${HOME}/.claude-octopus}}"
}

init_metrics_tracking() {
    local base
    base=$(get_metrics_base)
    local metrics_file="${base}/metrics-session.json"
    local metrics_dir="${base}/metrics-history"
    mkdir -p "$metrics_dir"
    local session_id="${SESSION_ID:-$(date +%Y%m%d-%H%M%S)}"
    local started_at
    started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    python3 - "$metrics_file" "$session_id" "$started_at" <<'PY'
import json, sys
path, session_id, started_at = sys.argv[1:]
data = {
    "session_id": session_id,
    "started_at": started_at,
    "phases": [],
    "totals": {
        "duration_seconds": 0,
        "estimated_tokens": 0,
        "native_tokens": 0,
        "tool_uses": 0,
        "estimated_cost_usd": 0,
        "agent_calls": 0,
        "native_metrics_available": False,
    },
}
with open(path, "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False, indent=2)
    f.write("\n")
PY
}

record_agent_start() {
    local agent_type="$1"
    local model="$2"
    local prompt="$3"
    local phase="${4:-unknown}"
    : "${agent_type}" "${model}" "${prompt}" "${phase}"
    local agent_id="${agent_type}-$(date +%s)-$$"
    local start_time
    start_time=$(date +%s)
    local base
    base=$(get_metrics_base)
    echo "$start_time" > "${base}/.agent-start-${agent_id}"
    echo "$agent_id"
}

get_model_cost() {
    local model="$1"
    case "$model" in
        claude-opus-4-6)        echo "5.00" ;;
        claude-opus-4-5)        echo "15.00" ;;
        claude-sonnet-4-5)      echo "3.00" ;;
        claude-sonnet-4)        echo "3.00" ;;
        claude-haiku-*)         echo "0.25" ;;
        gpt-5.3-codex)          echo "4.00" ;;
        gpt-5*)                 echo "3.00" ;;
        gpt-4*)                 echo "3.00" ;;
        gemini-2.0-pro*)        echo "2.50" ;;
        gemini-2.0-flash*)      echo "0.30" ;;
        gemini-3-pro*)          echo "3.00" ;;
        gemini-3-flash*)        echo "0.25" ;;
        *)                      echo "1.00" ;;
    esac
}

record_agent_complete() {
    local agent_id="$1"
    local agent_type="$2"
    local model="$3"
    local output="$4"
    local phase="${5:-unknown}"
    local native_token_count="${6:-}"
    local native_tool_uses="${7:-}"
    local native_duration_ms="${8:-}"

    local base
    base=$(get_metrics_base)
    local metrics_file="${base}/metrics-session.json"
    local start_file="${base}/.agent-start-${agent_id}"
    [[ -f "$start_file" ]] || return 1

    local start_time end_time duration
    start_time=$(cat "$start_file")
    end_time=$(date +%s)
    duration=$((end_time - start_time))

    local has_native="false"
    local token_count=0
    local tool_use_count=0
    if [[ -n "$native_token_count" && "$native_token_count" =~ ^[0-9]+$ ]]; then
        has_native="true"
        token_count="$native_token_count"
        tool_use_count="${native_tool_uses:-0}"
        if [[ -n "$native_duration_ms" && "$native_duration_ms" =~ ^[0-9]+$ ]]; then
            duration=$(( native_duration_ms / 1000 ))
        fi
    fi

    local output_length estimated_output_tokens estimated_total_tokens cost_basis_tokens
    output_length=${#output}
    estimated_output_tokens=$((output_length / 4))
    estimated_total_tokens=$((estimated_output_tokens + 100))
    cost_basis_tokens="$estimated_total_tokens"
    [[ "$has_native" == "true" ]] && cost_basis_tokens="$token_count"

    local cost_per_1k estimated_cost
    cost_per_1k=$(get_model_cost "$model")
    estimated_cost=$(awk "BEGIN {printf \"%.4f\", ($cost_basis_tokens / 1000.0) * $cost_per_1k}")
    local timestamp
    timestamp="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

    python3 - "$metrics_file" "$agent_type" "$model" "$phase" "$duration" "$estimated_total_tokens" "$estimated_cost" "$timestamp" "$has_native" "$token_count" "$tool_use_count" <<'PY'
import json, sys
(
    path, agent_type, model, phase, duration, est_tokens, est_cost, timestamp,
    has_native, native_tokens, tool_uses
) = sys.argv[1:]
duration = int(duration)
est_tokens = int(est_tokens)
est_cost = float(est_cost)
native_tokens = int(native_tokens)
tool_uses = int(tool_uses)
with open(path, "r", encoding="utf-8") as f:
    data = json.load(f)
entry = {
    "agent": agent_type,
    "model": model,
    "phase": phase,
    "duration_seconds": duration,
    "estimated_tokens": est_tokens,
    "estimated_cost_usd": est_cost,
    "timestamp": timestamp,
}
if has_native == "true":
    entry["native_token_count"] = native_tokens
    entry["native_tool_uses"] = tool_uses
    entry["metrics_source"] = "native"
else:
    entry["metrics_source"] = "estimated"
data.setdefault("phases", []).append(entry)
tot = data.setdefault("totals", {})
tot["duration_seconds"] = int(tot.get("duration_seconds", 0)) + duration
tot["estimated_tokens"] = int(tot.get("estimated_tokens", 0)) + est_tokens
tot["estimated_cost_usd"] = float(tot.get("estimated_cost_usd", 0)) + est_cost
tot["agent_calls"] = int(tot.get("agent_calls", 0)) + 1
tot.setdefault("native_tokens", 0)
tot.setdefault("tool_uses", 0)
tot.setdefault("native_metrics_available", False)
if has_native == "true":
    tot["native_tokens"] = int(tot.get("native_tokens", 0)) + native_tokens
    tot["tool_uses"] = int(tot.get("tool_uses", 0)) + tool_uses
    tot["native_metrics_available"] = True
with open(path + ".tmp", "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False, indent=2)
    f.write("\n")
PY
    mv "${metrics_file}.tmp" "$metrics_file"
    rm -f "$start_file"
}

display_phase_metrics() {
    local phase="$1"
    local base metrics_file
    base=$(get_metrics_base)
    metrics_file="${base}/metrics-session.json"
    [[ -f "$metrics_file" ]] || return 0
    python3 - "$metrics_file" "$phase" <<'PY'
import json, sys
path, phase = sys.argv[1], sys.argv[2]
with open(path, "r", encoding="utf-8") as f:
    data = json.load(f)
items = [p for p in data.get("phases", []) if p.get("phase") == phase]
if not items:
    raise SystemExit(0)
dur = sum(int(p.get("duration_seconds", 0)) for p in items)
tok = sum(int(p.get("estimated_tokens", 0)) for p in items)
cost = sum(float(p.get("estimated_cost_usd", 0)) for p in items)
native_tok = sum(int(p.get("native_token_count", 0)) for p in items)
tool_uses = sum(int(p.get("native_tool_uses", 0)) for p in items)
has_native = any(p.get("metrics_source") == "native" for p in items)
print("")
print(f"📊 Phase Metrics ({phase}):")
print(f"  ⏱️  Duration: {dur}s")
if has_native:
    print(f"  📝 Tokens: {native_tok} (native) / {tok} (est.)")
    print(f"  🔧 Tool Uses: {tool_uses}")
else:
    print(f"  📝 Est. Tokens: {tok}")
print(f"  💰 Est. Cost: ${cost:.4f}")
print(f"  🤖 Agents: {len(items)}")
PY
}

display_session_metrics() {
    local base metrics_file metrics_dir
    base=$(get_metrics_base)
    metrics_file="${base}/metrics-session.json"
    metrics_dir="${base}/metrics-history"
    [[ -f "$metrics_file" ]] || return 0
    python3 - "$metrics_file" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
t = data.get("totals", {})
duration = int(t.get("duration_seconds", 0))
tokens = int(t.get("estimated_tokens", 0))
cost = float(t.get("estimated_cost_usd", 0))
calls = int(t.get("agent_calls", 0))
native_tokens = int(t.get("native_tokens", 0))
tool_uses = int(t.get("tool_uses", 0))
has_native = bool(t.get("native_metrics_available", False))
print("")
print("═══════════════════════════════════════════")
print("📊 Session Totals")
print("═══════════════════════════════════════════")
print(f"  ⏱️  Total Duration: {duration}s ({duration/60:.1f}m)")
if has_native:
    print(f"  📝 Tokens: {native_tokens} (native) / {tokens} (est.)")
    print(f"  🔧 Tool Uses: {tool_uses}")
else:
    print(f"  📝 Est. Tokens: {tokens}")
print(f"  💰 Est. Cost: ${cost:.4f}")
print(f"  🤖 Agent Calls: {calls}")
print("═══════════════════════════════════════════")
PY
    mkdir -p "$metrics_dir"
    cp "$metrics_file" "${metrics_dir}/session-$(date +%Y%m%d-%H%M%S).json"
}

display_provider_breakdown() {
    local base metrics_file
    base=$(get_metrics_base)
    metrics_file="${base}/metrics-session.json"
    [[ -f "$metrics_file" ]] || return 0
    python3 - "$metrics_file" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
groups = {"codex": [0,0.0], "gemini": [0,0.0], "claude": [0,0.0]}
for p in data.get("phases", []):
    agent = str(p.get("agent", ""))
    for key in groups:
        if agent.startswith(key):
            groups[key][0] += int(p.get("estimated_tokens", 0))
            groups[key][1] += float(p.get("estimated_cost_usd", 0))
print("")
print("Provider Breakdown:")
if groups["codex"][0] > 0:
    print(f"  🔴 Codex:  {groups['codex'][0]} tokens (${groups['codex'][1]:.4f})")
if groups["gemini"][0] > 0:
    print(f"  🟡 Gemini: {groups['gemini'][0]} tokens (${groups['gemini'][1]:.4f})")
if groups["claude"][0] > 0:
    print(f"  🔵 Claude: {groups['claude'][0]} tokens (${groups['claude'][1]:.4f})")
PY
}

record_agents_batch_complete() {
    local phase="$1"
    local task_group="$2"
    local base metrics_map
    base=$(get_metrics_base)
    metrics_map="${base}/.metrics-map"
    [[ -f "$metrics_map" ]] || return 0

    local result filename task_id metrics_line metrics_id agent_type model output
    for result in "$RESULTS_DIR"/${phase}-${task_group}-*.md; do
        [[ -f "$result" ]] || continue
        filename=$(basename "$result" .md)
        task_id="${filename#${phase}-${task_group}-}"
        metrics_line=$(grep "^${task_group}-${task_id}:" "$metrics_map" 2>/dev/null || true)
        [[ -n "$metrics_line" ]] || continue
        IFS=':' read -r _ metrics_id agent_type model <<< "$metrics_line"
        output=$(cat "$result" 2>/dev/null || echo "")
        local native_tokens="" native_tools="" native_duration=""
        if declare -f parse_task_metrics &>/dev/null; then
            parse_task_metrics "$output"
            native_tokens="$_PARSED_TOKENS"
            native_tools="$_PARSED_TOOL_USES"
            native_duration="$_PARSED_DURATION_MS"
        fi
        record_agent_complete "$metrics_id" "$agent_type" "$model" "$output" "$phase" \
            "$native_tokens" "$native_tools" "$native_duration"
        sed -i.bak "/^${task_group}-${task_id}:/d" "$metrics_map" 2>/dev/null || true
    done
}

display_per_phase_cost_table() {
    local base metrics_file
    base=$(get_metrics_base)
    metrics_file="${base}/metrics-session.json"
    [[ -f "$metrics_file" ]] || return 0
    python3 - "$metrics_file" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
rows = data.get("phases", [])
if not rows:
    raise SystemExit(0)
print("")
print("Per-Phase Cost Breakdown:")
print("┌──────────┬──────────────┬────────────┬──────────┬──────────┐")
print("│ Phase    │ Provider     │ Tokens     │ Cost     │ Duration │")
print("├──────────┼──────────────┼────────────┼──────────┼──────────┤")
for r in rows:
    phase = str(r.get("phase", ""))[:8]
    agent = str(r.get("agent", ""))
    if agent.startswith("codex"):
        provider = "🔴 codex"
    elif agent.startswith("gemini"):
        provider = "🟡 gemini"
    elif agent.startswith("claude"):
        provider = "🔵 claude"
    else:
        provider = agent[:12]
    tokens = int(r.get("estimated_tokens", 0))
    cost = float(r.get("estimated_cost_usd", 0))
    dur = int(r.get("duration_seconds", 0))
    print(f"│ {phase:<8} │ {provider:<12} │ {tokens:>10} │ ${cost:<7.3f} │ {dur:>7}s │")
print("└──────────┴──────────────┴────────────┴──────────┴──────────┘")
PY
}

export -f get_metrics_base
export -f init_metrics_tracking
export -f record_agent_start
export -f record_agent_complete
export -f record_agents_batch_complete
export -f get_model_cost
export -f display_phase_metrics
export -f display_session_metrics
export -f display_provider_breakdown
export -f display_per_phase_cost_table
