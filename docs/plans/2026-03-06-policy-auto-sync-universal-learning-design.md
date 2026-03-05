# Policy Auto Sync Universal Learning Design

## Goal

Build a Go-native, tool-agnostic policy auto-sync architecture that enables invisible intelligent learning for users across any target project, while keeping safety-critical controls explicit and auditable.

## Problem Statement

The current plugin runtime has strong static policy and hook guardrails, but lacks a closed-loop learning system for:

1. repetitive manual reminders across tools,
2. dynamic project-specific rule evolution,
3. cold-start policy quality before runtime events exist,
4. unified event ownership and context continuity,
5. safe implicit learning,
6. persistent context injection into every run,
7. privacy-safe cross-project reuse,
8. storage-vs-accuracy trade-offs at scale.

## Final Product Decisions

1. Core architecture is `Policy Auto Sync` in Go; LLM only interprets/explains.
2. Event model is raw-fact-first (`no conclusion in ingest layer`).
3. Rule layering is:
   - `plugin defaults < project rules < session override < auto-learned overlays`
4. Non-safety high-confidence proposals can auto-activate.
5. Safety-critical rules never auto-activate and always require explicit confirmation.
6. Users can implicitly refine policies via normal input; accepted changes persist directly to `auto-learned overlay`.
7. User-negative feedback triggers ask-question confirmation with two actions: `delete` or `skip-this-session`.
8. If user confirms `delete`, policy is revoked immediately, related learning data is deleted, cooldown is applied, and relearning restarts from zero after cooldown.
9. If user confirms `skip-this-session`, rule is suppressed only for current session prompt injection.
10. Every plugin execution injects active policy context into prompt paths, including external vibe coding tool calls.
11. Cross-project reuse stores only desensitized preference patterns, guarded by similarity threshold, cooldown, and rollback.
12. Storage uses short-lived raw logs + long-lived aggregates with hard quota and adaptive degradation.

## Architecture Overview

Unified pipeline:

`Raw Events -> Detector Registry -> Signal Graph -> Scoring/Proposal -> Policy Sync -> Overlay Activation -> Prompt Injection -> Doctor/Hook Governance`

Execution ownership:

1. Go runtime decides collection, scoring, activation, rollback, cleanup, and injection.
2. LLM consumes deterministic policy snapshot/context and provides natural-language explanation only.

## Core Components

### 1. Event Ingestion

- Multi-entry event capture:
  - hooks (`PreToolUse`, `PostToolUse`, `SessionStart`, `Stop`, worktree events),
  - `mp` and `mp-devx` command pre/post,
  - policy operations (sync/apply/ignore/rollback/revoke),
  - optional git sampler.
- EventSink is append-only JSONL.
- Facts only: timestamps, actor/tool/action/result, execution contract, scope, evidence references.

### 2. Detector Registry

- Pluggable detectors produce universal dimensions, not project-type labels.
- Initial detector classes:
  - branching,
  - workspace,
  - command_contract,
  - risk_profile,
  - test discipline,
  - review discipline,
  - external tool behavior,
  - monorepo partitioning.
- Language/domain modules (Go/Node/Python/Java/Docs/Data/Infra) are plugins behind a stable detector interface.

### 3. Scoring And Proposal Engine

- Proposal evidence metrics:
  - support,
  - conflict rate,
  - cross-session stability,
  - time decay,
  - confidence.
- Proposal state machine:
  - `observed -> advisory -> shadow -> auto-candidate -> auto-applied`
  - terminal/alternate states:
    - `manual-required` (safety critical),
    - `ignored`,
    - `revoked`,
    - `rolled-back`,
    - `expired`.

### 4. Policy Sync Engine

- Default mode: dry-run analysis plus transparent auto-apply for qualified non-safety proposals.
- Explicit controls: apply, ignore, rollback, revoke.
- Activation writes to `auto-learned overlay` atomically.
- Drift/conflict can demote active policy back to shadow before rollback.

### 5. Prompt Injection Engine

- Every plugin run injects policy context into prompt path.
- Injection scope:
  - `/mp:*` workflows,
  - `mp` command execution path,
  - orchestration phase prompts,
  - external vibe coding tool call prompts.
- Injection payload is deterministic `PolicyContext` generated in Go:
  - active rules,
  - guardrails,
  - session-suppressed exclusions,
  - revoked/cooldown exclusions,
  - rule references (`rule_id`) and precedence.
- If external tool lacks direct prompt channel, adapter fallback uses context file/env mapping.
- Injection failures degrade gracefully (non-blocking), with audit and doctor warning.

### 6. Dynamic Context Snapshot

- Go-side deterministic snapshot generation; LLM does not author policy data.
- Snapshot includes:
  - active overlays,
  - recent drift signals,
  - pending safety confirmations,
  - recent user-feedback policy changes.

### 7. Cross-Project Semantic Fingerprint

- Data model stores preference patterns only (not project private raw values).
- Default local-only, desensitized, TTL-based lifecycle, one-click purge.
- Migration guardrail:
  - similarity threshold,
  - cooldown window,
  - rollback-ready activation path,
  - shadow-first where confidence is marginal.

## Cold Start Design

### Static Fingerprint Bootstrap

`init-fingerprint` extracts capability vectors from deterministic evidence:

1. config/manifests,
2. repository layout,
3. docs and markdown evidence.

Mandatory doc probe includes:

- `README.md`,
- `CLAUDE.md`,
- `AGENTS.md`,
- `PRODUCT.md`,
- common docs such as `ARCHITECTURE.md`, `CONTRIBUTING.md`, `docs/**/product*.md`, `docs/**/tech-stack*.md`, `docs/**/getting-started*.md`.

Output format:

- capability vector (`vcs_model`, `build_tool`, `test_harness`, `ci_provider`, `repo_shape`, `risk_profile`),
- per-dimension confidence,
- `evidence_map` (which source produced each dimension).

Evidence precedence:

`runtime events > config/manifests > markdown docs`

## User-Driven Policy Refinement

### Positive/Corrective Feedback

- LLM extracts structured policy-intent from user input.
- Go validator performs schema/safety/conflict checks.
- On pass, mutation persists directly to `auto-learned overlay`.

### Negative Feedback (Revoke)

When user denies a policy, runtime must first ask:

- `delete policy`
- `skip this session only`

#### If user chooses `delete policy`

1. remove active rule immediately from overlay,
2. mark proposal as `revoked_by_user`,
3. delete policy-related learning data:
   - proposal history entries,
   - shadow/active state links,
   - related samples,
   - aggregate contribution for that rule,
4. write minimal audit line (timestamp, rule_id, reason),
5. apply cooldown (`revoked_until`) during which auto-rebuild is blocked,
6. after cooldown, relearn from zero with no inherited support/confidence.

#### If user chooses `skip this session only`

1. keep persisted overlay and learning data unchanged,
2. write session-local suppression state,
3. exclude suppressed rule from prompt injection for current session only,
4. clear suppression automatically at session end.

## Monorepo Partition Strategy

- Learn and apply rules per subdirectory/service partition.
- Root-level rules remain minimal shared baseline.
- Resolution order:
  - partition overlay,
  - root overlay,
  - lower rule layers.

## Storage And Retention Strategy

### Files

- `.multipowers/policy/autosync/events.raw.YYYY-MM-DD.jsonl`
- `.multipowers/policy/autosync/proposals.jsonl`
- `.multipowers/policy/autosync/applied.jsonl`
- `.multipowers/policy/autosync/overlays.auto.json`
- `.multipowers/policy/autosync/daily_stats.json`
- `.multipowers/policy/autosync/signal_samples.jsonl`
- `.multipowers/policy/autosync/context.snapshot.json`
- `.multipowers/policy/autosync/project.fingerprint.json`
- `~/.multipowers/policy/autosync/global.semantic.json` (cross-project, desensitized)

### Retention And Compression

1. raw logs: daily split, 5MB rotate, gzip, retain 7-14 days,
2. aggregates: retain 90-180 days,
3. dedup: same `event_key` merged in 10-minute window,
4. evidence cap: up to 200 samples per signal.

### Hard Quotas

1. per project: 50MB,
2. global: 500MB.

When over global quota, cleanup uses `LRU + cumulative reference count`:

- first delete least-recently-used and low-reference assets,
- keep high-reference assets longer even if short-term inactive,
- tie-break order: oldest raw -> low-value samples -> rebuildable caches.

At 80% quota, switch to aggregate-first degradation mode.

## Activation Thresholds

Default auto-activation gate for non-safety rules:

- `support >= 8`
- `sessions >= 3`
- `conflict_rate < 15%`
- `confidence >= 0.95`
- stable behavior across at least 3-5 implicit validations.

High conflict rules remain in shadow only.

## Governance Integration

1. hooks: hard-stop on mandatory guardrails.
2. doctor: warn for
   - high-confidence unapplied proposals,
   - rule drift,
   - recurrence concentration,
   - negative migration signals.
3. all apply/rollback/revoke actions are append-only auditable.

## Runtime Commands

- `mp policy sync`
- `mp policy sync --apply`
- `mp policy sync --ignore <id>`
- `mp policy sync --rollback <id>`
- `mp policy stats`
- `mp policy gc`
- `mp policy tune --mode balanced|accuracy|storage`

User-facing default remains `/mp:*` and `mp`; `mp-devx` is maintainer/CI-focused.

## Failure Isolation

1. learning path failures never block main workflow execution,
2. detector failure is isolated per detector,
3. injection failure degrades with audit warning,
4. storage pressure triggers deterministic downshift before hard failure,
5. rollback path is always available for activated overlays.

## Acceptance Criteria

1. Non-safety high-confidence rules auto-activate without user friction.
2. Safety-critical rules never auto-activate.
3. Every plugin execution path (including external vibe coding tool) receives policy injection.
4. User-selected `delete` revokes rule immediately and resets relearning baseline after cooldown.
5. User-selected `skip-this-session` suppresses only current-session injection and does not delete persisted learning data.
6. Cross-project reuse remains local, desensitized, and guarded against negative transfer.
7. Storage stays within quota using `LRU + cumulative reference count` cleanup.
