# Benchmark + Smart Routing Design

## Goal

Add a configuration-driven benchmark system for `/mp:*` commands that:

1. Detects code-related intent.
2. Forces fan-out execution to all configured available models when enabled.
3. Collects runtime and quality metrics asynchronously.
4. Persists cross-project benchmark history in local daily JSONL files.
5. Optionally overrides default routing using historical top-performing models for similar scenarios.

All benchmark/ranking/storage failures must not impact primary command execution.

## Confirmed Decisions

1. Trigger scope: all `/mp:*` commands, then classify code-related intent.
2. Code-intent decision: whitelist feature extraction plus LLM semantic judgment; LLM is final arbiter.
3. Storage: local JSONL store (not SQLite), cross-project path under user home.
4. Partitioning: daily files (`YYYY-MM-DD.*.jsonl`).
5. Fan-out scope: all available configured models from executor/model registry.
6. Judge: single configured `judge_model`.
7. Scoring: dimension scores use `1-5`, then weighted aggregate.
8. Smart routing override is controlled by a separate config toggle.
9. Override can happen only when smart routing toggle is on.
10. Minimum sample gate for override: `N=10` per candidate model per similar scenario.
11. Async-first pipeline; failures never block or fail main workflow.

## Architecture

### Main Path (Latency-Sensitive)

1. Parse `/mp:*` input and extract lightweight whitelist features.
2. Enqueue async classification request.
3. Resolve effective routing decision:
- if benchmark mode is off: normal routing.
- if benchmark mode is on and classified as code-related: override to all available models.
4. Execute main orchestration.
5. Return primary result as usual.

Main path must not wait for benchmark persistence or judge scoring.

### Async Workers (Best-Effort)

1. `intent_worker`
- Runs LLM semantic classification and emits finalized code-intent decision.
2. `metrics_worker`
- Records run metadata and model output metadata (duration, tokens, fallback, errors, tags).
3. `judge_worker`
- Uses configured `judge_model` to score all model outputs for the run.
4. `routing_worker`
- Recomputes best model by similar-scenario aggregates.
- Produces override decisions only when smart routing toggle is enabled.

Any worker failure is isolated and logged to error JSONL.

## Configuration Surface

```yaml
benchmark_mode:
  enabled: true
  async_enabled: true
  force_all_models_on_code: true
  judge_model: "claude-opus"
  code_intent:
    whitelist:
      task_types: []
      tech_features: []
      frameworks: []
      languages: []
    llm_semantic_judge: true
    llm_decision_priority: true
  storage:
    type: jsonl
    root: "~/.multipowers/metrics"
    partition: daily
  scoring:
    scale: "1-5"
    dimensions:
      - correctness
      - code_quality
      - performance
      - security
      - design_fitness
      - completeness
      - framework_alignment
      - testability
      - maintainability
      - clarity
    weights: {}
  fault_tolerance:
    never_block_main_flow: true
    retry_max: 3
    timeout_ms: 20000

smart_routing:
  enabled: false
  override_existing_routing_when_on: true
  strategy: performance_optimized
  min_samples_per_model: 10
  match_keys:
    - task_type
    - tech_features
    - framework
    - language
```

## JSONL Data Model

Root path: `~/.multipowers/metrics`

Daily files:

1. `runs.YYYY-MM-DD.jsonl`
2. `model_outputs.YYYY-MM-DD.jsonl`
3. `task_fingerprints.YYYY-MM-DD.jsonl`
4. `judge_scores.YYYY-MM-DD.jsonl`
5. `route_overrides.YYYY-MM-DD.jsonl`
6. `async_jobs.YYYY-MM-DD.jsonl`
7. `errors.YYYY-MM-DD.jsonl`

### Required Record Fields

`runs`:
- run_id
- timestamp_start
- timestamp_end
- command
- prompt_hash
- benchmark_mode_enabled
- smart_routing_enabled
- code_intent_final

`model_outputs`:
- run_id
- model
- provider
- duration_ms
- tokens_input
- tokens_output
- status
- fallback_used
- error_code

`task_fingerprints`:
- run_id
- task_type
- tech_features[]
- framework
- language
- whitelist_hits

`judge_scores`:
- run_id
- judged_model
- judge_model
- dimension_scores (1-5)
- weighted_score
- rationale_summary

`route_overrides`:
- run_id
- override_applied
- previous_model
- selected_model
- match_signature
- sample_count
- strategy

`async_jobs`:
- job_id
- job_type
- status
- attempts
- latency_ms

`errors`:
- job_id
- stage
- error_class
- message
- retryable

## Similar-Scenario Matching

Similarity key is a structured tag signature:

`task_type + tech_features + framework + language`

Routing strategy `performance_optimized`:

1. Filter historical judged records by exact/normalized tag signature.
2. Keep candidates with sample count >= `min_samples_per_model` (`10`).
3. Compute weighted average score per model.
4. Select highest score as recommended model.
5. Apply override only when `smart_routing.enabled=true`.

If no candidate satisfies sample gate, do not override.

## Failure Isolation and Degradation

1. Async path is best-effort; no error propagation to primary workflow result.
2. Queue saturation causes sampling/deferred writes instead of blocking.
3. Worker timeout or judge failure records error events only.
4. Storage write failure records internal error and retries with bounded attempts.
5. If retry budget exhausted, persist error metadata and drop event safely.

## Acceptance Criteria

1. With `benchmark_mode.enabled=true`, code-related `/mp:*` requests fan out to all configured available models.
2. Runtime records include duration, tokens, task tags (type/features/framework/language), and per-model outcomes.
3. Judge model emits structured multi-dimensional `1-5` scores and weighted aggregate.
4. JSONL files are written under cross-project global root with daily partitioning.
5. With `smart_routing.enabled=false`, no routing override occurs.
6. With `smart_routing.enabled=true`, override occurs only when sample count per model is at least `10`.
7. Any benchmark/scoring/storage failure does not fail or delay main command completion.

## Out of Scope

1. Cloud dashboard or remote telemetry services.
2. Manual human scoring workflows.
3. Multi-judge ensemble scoring in first version.
4. Online learning/bandit routing in first version.
