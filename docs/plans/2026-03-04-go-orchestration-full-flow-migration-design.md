# Go Orchestration Full-Flow Migration Design

## Goal

Migrate advanced workflow business flows from `main` shell orchestration to Go-native orchestration in the Go branch, with flow-level behavioral equivalence (not shell-level source equivalence), covering all major flows:

- `discover`
- `define`
- `develop`
- `deliver`
- `debate`
- `embrace`

This migration is a one-way cutover. No temporary env-flag fallback path is kept.

## Scope

### In Scope

1. Build a unified Go orchestration engine for multi-step flow execution.
2. Reproduce key high-level orchestration mechanisms from shell flow behavior:
- perspective decomposition
- parallel execution
- synchronization and progress tracking
- progressive synthesis
- final synthesis
- one-hop fallback-aware dispatch
3. Move global orchestration defaults into `config/orchestration.yaml`.
4. Keep flow-specific overrides in `config/workflows.yaml`.
5. Switch workflow entrypoints to the new orchestration runtime.

### Out of Scope

1. Source-level line-by-line shell migration.
2. Keeping legacy shell orchestration as runtime fallback.
3. Multi-hop fallback chains beyond current one-hop policy contract.

## Architecture

The Go orchestration runtime is organized into three core layers:

1. Planner layer
2. Executor layer
3. Synthesizer layer

### 1) Planner

Input:

- workflow name
- optional task id
- prompt
- project directory
- config from:
  - `config/workflows.yaml`
  - `config/orchestration.yaml`
  - `config/agents.yaml`
  - `config/providers.yaml`

Output:

- `ExecutionPlan`
  - resolved phases
  - perspectives for each phase
  - selected candidates/agents
  - concurrency controls
  - progressive/final synthesis controls
  - loop controls (from global defaults, optionally overridden per flow/task)

### 2) Executor

Responsibilities:

1. Execute plan steps concurrently using goroutines/worker pool.
2. Dispatch each step through policy-based `DispatchWithFallback`.
3. Collect structured step results:
- model/executor used
- fallback/degraded info
- duration
- output payload

4. Synchronize completion and propagate cancellation/timeout cleanly.

### 3) Synthesizer

Responsibilities:

1. Progressive synthesis:
- trigger when minimum completed steps and minimum valid output thresholds are met.
- produce intermediate synthesis artifacts to reduce waiting.

2. Final synthesis:
- aggregate valid step outputs after execution completes.
- produce final unified report and conflict/convergence summary.

## Configuration Strategy

### `config/orchestration.yaml` (Global defaults)

Holds reusable orchestration defaults across flows:

- `phase_defaults`
- `ralph_wiggum` loop defaults
- `skill_triggers`
- default perspective/parallel/synthesis thresholds

### `config/workflows.yaml` (Flow-specific overrides only)

Holds flow-level (and optional task-level) overrides:

- model
- fallback policy
- flow-specific phase/perspective/parallel/synthesis overrides

No requirement to include global defaults here.

### `config/providers.yaml`

Holds execution strategy:

- provider kind
- command template
- enforcement
- model patterns
- fallback policies

Workflow should primarily declare model intent; provider resolution is inferred through model patterns unless explicit executor mapping is provided.

### `config/agents.yaml`

Holds agent persona metadata and runtime-relevant agent attributes, while orchestration controls stay outside.

## Resolution Precedence

For orchestration plan building:

1. `workflow.tasks.<task>` override
2. `workflow.default` override
3. `orchestration.yaml` global defaults

For provider dispatch:

1. explicit profile (if present)
2. model-pattern provider inference from `providers.yaml`
3. deterministic fallback selection rule

## Flow Equivalence Target

Target is process equivalence with shell behavior for major orchestration mechanisms, not literal shell implementation parity.

Required equivalence points:

1. multi-perspective decomposition
2. parallel execution and synchronization
3. progressive synthesis trigger and intermediate result availability
4. final synthesis generation
5. robust failure handling and fallback-aware execution metadata

## Error Handling Model

Standardized error classes:

1. `config_error`
2. `plan_error`
3. `dispatch_error`
4. `synthesis_error`

Execution should return structured outcomes even under partial failures, with degraded/fallback traces preserved.

## Observability Contract

Each run emits structured orchestration metadata:

- `execution_plan_id`
- step lifecycle records
- fallback trace
- progressive synthesis events
- final synthesis status and report summary

## Cutover Policy

No backward-compatibility runtime toggle is kept.

Cutover is enforced by:

1. strict preflight config validation
2. workflow entrypoint migration to orchestration engine
3. full test verification before declaring completion

## Validation Matrix

### Planner

1. precedence resolution (`task > workflow > global`)
2. missing/invalid phase configuration behavior
3. deterministic plan generation

### Executor

1. concurrency limits respected
2. cancellation and timeout correctness
3. fallback/degraded metadata integrity

### Progressive Synthesis

1. trigger thresholds
2. no premature trigger
3. repeat trigger behavior correctness

### Final Synthesis

1. valid-result filtering
2. minimum-valid-input behavior
3. complete report structure

### End-to-End Flow Coverage

At least one success path and one failure/degraded path per flow:

- discover
- define
- develop
- deliver
- debate
- embrace

## Risks and Mitigations

1. Risk: behavior drift from shell flow expectations.
- Mitigation: flow-level equivalence tests and explicit perspective/synthesis contract tests.

2. Risk: increased orchestration complexity in one-step cutover.
- Mitigation: strict module boundaries (planner/executor/synthesizer) and hard validation gates.

3. Risk: config ambiguity between global defaults and flow overrides.
- Mitigation: strict precedence definition and traceable resolved plan output.
