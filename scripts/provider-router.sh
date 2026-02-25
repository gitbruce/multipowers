#!/usr/bin/env bash
# Provider Router for Claude Octopus
# Latency-based provider routing with round-robin, fastest, and cheapest strategies.

OCTOPUS_ROUTING_MODE="${OCTOPUS_ROUTING_MODE:-round-robin}"

_ROUTER_STATE_FILE="${WORKSPACE_DIR:-${HOME}/.claude-octopus}/.router-state"
_ROUTER_STATS_FILE="${WORKSPACE_DIR:-${HOME}/.claude-octopus}/.provider-stats.json"

build_provider_stats() {
    local metrics_dir="${WORKSPACE_DIR:-${HOME}/.claude-octopus}"
    local metrics_file="${metrics_dir}/metrics-session.json"
    local stats_file="$_ROUTER_STATS_FILE"

    [[ -f "$metrics_file" ]] || return 1
    mkdir -p "$metrics_dir"

    python3 - "$metrics_file" "$stats_file" <<'PY'
import datetime, json, sys
metrics_file, stats_file = sys.argv[1], sys.argv[2]
with open(metrics_file, "r", encoding="utf-8") as f:
    data = json.load(f)
providers = {}
for p in data.get("phases", []):
    agent = str(p.get("agent", ""))
    if not agent:
        continue
    base = agent.split("-")[0]
    slot = providers.setdefault(base, {"lat": [], "cost": []})
    slot["lat"].append(float(p.get("duration_seconds", 0)) * 1000.0)
    slot["cost"].append(float(p.get("estimated_cost_usd", 0)))
result = {"providers": {}, "updated_at": datetime.datetime.utcnow().replace(microsecond=0).isoformat() + "Z"}
for name, slot in providers.items():
    cnt = max(1, len(slot["lat"]))
    result["providers"][name] = {
        "avg_latency_ms": (sum(slot["lat"]) / cnt),
        "call_count": len(slot["lat"]),
        "avg_cost_usd": (sum(slot["cost"]) / cnt),
    }
with open(stats_file, "w", encoding="utf-8") as f:
    json.dump(result, f, ensure_ascii=False, indent=2)
    f.write("\n")
PY
}

read_provider_metric() {
    local stats_file="$1"
    local provider="$2"
    local metric="$3"
    local fallback="$4"
    [[ -f "$stats_file" ]] || { echo "$fallback"; return 0; }
    python3 - "$stats_file" "$provider" "$metric" "$fallback" <<'PY'
import json, sys
path, provider, metric, fallback = sys.argv[1:]
try:
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)
    print(data.get("providers", {}).get(provider, {}).get(metric, fallback))
except Exception:
    print(fallback)
PY
}

select_fastest_provider() {
    local stats_file="$_ROUTER_STATS_FILE"
    local candidates=("$@")
    [[ ${#candidates[@]} -gt 0 ]] || return 1

    case "$OCTOPUS_ROUTING_MODE" in
        round-robin)
            local idx=0
            [[ -f "$_ROUTER_STATE_FILE" ]] && idx=$(cat "$_ROUTER_STATE_FILE" 2>/dev/null || echo "0")
            local selected="${candidates[$((idx % ${#candidates[@]}))]}"
            echo $(( (idx + 1) % ${#candidates[@]} )) > "$_ROUTER_STATE_FILE"
            echo "$selected"
            ;;
        fastest)
            local best="" best_latency=999999
            local candidate base_provider latency
            for candidate in "${candidates[@]}"; do
                base_provider="${candidate%%-*}"
                latency=$(read_provider_metric "$stats_file" "$base_provider" "avg_latency_ms" "999999")
                if awk -v a="$latency" -v b="$best_latency" 'BEGIN { exit !(a < b) }'; then
                    best="$candidate"; best_latency="$latency"
                fi
            done
            echo "${best:-${candidates[0]}}"
            ;;
        cheapest)
            local best="" best_cost=999999
            local candidate base_provider cost
            for candidate in "${candidates[@]}"; do
                base_provider="${candidate%%-*}"
                cost=$(read_provider_metric "$stats_file" "$base_provider" "avg_cost_usd" "999999")
                if awk -v a="$cost" -v b="$best_cost" 'BEGIN { exit !(a < b) }'; then
                    best="$candidate"; best_cost="$cost"
                fi
            done
            echo "${best:-${candidates[0]}}"
            ;;
        *)
            echo "${candidates[0]}"
            ;;
    esac
}

refresh_provider_stats() {
    build_provider_stats 2>/dev/null || true
}
