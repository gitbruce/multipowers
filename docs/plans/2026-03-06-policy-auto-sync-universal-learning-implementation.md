# Policy Auto Sync Universal Learning Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Deliver a Go-native Policy Auto Sync system that performs invisible policy learning, policy prompt injection (including external vibe coding tools), safe auto-activation, user-driven delete-or-skip confirmation on policy denial, and cross-project preference reuse for any target project.

**Architecture:** Introduce a new `internal/autosync` domain composed of ingestion, detectors, scoring, overlay activation, feedback confirmation (`delete` vs `skip-this-session`), fingerprint bootstrap, semantic reuse, and prompt injection. Integrate it with existing hooks/CLI/policy/doctor layers while preserving runtime safety boundaries and append-only governance logs.

**Tech Stack:** Go (`os`, `filepath`, `encoding/json`, `sync`, `time`), existing `internal/hooks`, `internal/cli`, `internal/policy`, `internal/doctor`, JSONL append-only logs under `.multipowers`, `go test`.

---

## Status Bar

**Overall:** `0/12 tasks complete (0%)`

| Task ID | Task Name | Status | Last Updated | Notes |
|---|---|---|---|---|
| T1 | Autosync contracts and paths | pending | - | - |
| T2 | Raw EventSink + dedup + rotation | pending | - | - |
| T3 | Multi-entry event emission wiring | pending | - | - |
| T4 | Detector Registry + generic signals | pending | - | - |
| T5 | Scoring engine + proposal lifecycle | pending | - | - |
| T6 | Overlay activation + deny-confirmation + revoke/cooldown/reset | pending | - | - |
| T7 | `init-fingerprint` doc-aware bootstrap | pending | - | - |
| T8 | Cross-project semantic preference reuse | pending | - | - |
| T9 | PolicyContext snapshot + prompt injection | pending | - | - |
| T10 | `mp policy` command surface | pending | - | - |
| T11 | Hook/Doctor governance integration | pending | - | - |
| T12 | E2E verification + docs hardening | pending | - | - |

### Status Update Rules

1. When starting a task, set its status to `in_progress` and update `Last Updated`.
2. When all task acceptance checks pass, set status to `done`, update `Overall`, and add commit SHA in `Notes`.
3. If blocked, set status to `blocked` and add blocker + owner in `Notes`.
4. Never mark `done` without command output evidence in the task section.

### Required Execution Skills

- Use `@superpowers:test-driven-development` before each code change.
- Use `@superpowers:verification-before-completion` before marking any task `done`.
- Use `@superpowers:systematic-debugging` if any test fails unexpectedly.

---

### Task 1: Establish Autosync Domain Contracts and Storage Paths

**Why:** A stable schema and path contract is needed so every later component writes/reads compatible artifacts.

**What:** Create `internal/autosync` base types, constants, and path helpers for all autosync artifacts.

**How:** Define typed models for raw events, signals, proposals, overlays, cooldown metadata, and project/global storage roots.

**Key Design:** Paths are deterministic and centralized; no path literals outside the path helper package.

**Files:**
- Create: `internal/autosync/types.go`
- Create: `internal/autosync/paths.go`
- Create: `internal/autosync/types_test.go`

**Step 1: Write the failing test**

```go
func TestDefaultPaths_AreStable(t *testing.T) {
    got := DefaultPaths("/tmp/p")
    if got.ProjectRoot != "/tmp/p/.multipowers/policy/autosync" {
        t.Fatalf("project root mismatch: %s", got.ProjectRoot)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/autosync -run TestDefaultPaths_AreStable -v`  
Expected: FAIL (package/files missing)

**Step 3: Write minimal implementation**

```go
type Paths struct{ ProjectRoot string }
func DefaultPaths(projectDir string) Paths {
    return Paths{ProjectRoot: filepath.Join(projectDir, ".multipowers", "policy", "autosync")}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/autosync -run TestDefaultPaths_AreStable -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/types.go internal/autosync/paths.go internal/autosync/types_test.go
git commit -m "feat(autosync): add base contracts and storage path helpers"
```

---

### Task 2: Implement Append-Only Raw EventSink with Dedup/Rotation Contract

**Why:** Learning quality and storage control depend on raw event durability plus bounded volume.

**What:** Build JSONL append sink with event-key dedup window and file rotation metadata.

**How:** Add sink writer, dedup cache (`event_key + 10m`), and daily file naming contract.

**Key Design:** Ingestion is fact-only append; dedup compacts repeated events by `count++` without losing first/last timestamps.

**Files:**
- Create: `internal/autosync/store/raw_sink.go`
- Create: `internal/autosync/store/dedup.go`
- Create: `internal/autosync/store/raw_sink_test.go`

**Step 1: Write the failing tests**

```go
func TestRawSink_AppendsJSONL(t *testing.T) {}
func TestRawSink_DedupWindowMergesEventKey(t *testing.T) {}
func TestRawSink_DailyFileNaming(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/store -run TestRawSink_ -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- `AppendRawEvent(event)` -> append JSON line to `events.raw.YYYY-MM-DD.jsonl`
- `DedupWindow.Apply(event)` -> merge by key within 10m
- return file path + dedup info for observability

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/store -run TestRawSink_ -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/store/raw_sink.go internal/autosync/store/dedup.go internal/autosync/store/raw_sink_test.go
git commit -m "feat(autosync): add append-only raw event sink with dedup window"
```

---

### Task 3: Wire Multi-Entry Event Emission (Hooks, mp/mp-devx, Policy Ops)

**Why:** Autosync cannot learn without broad event coverage and consistent ownership.

**What:** Emit normalized raw events from all required runtime entry points.

**How:** Integrate event emitter into hook handlers and CLI command lifecycle.

**Key Design:** Emission failures are non-blocking for runtime commands but always logged for doctor visibility.

**Files:**
- Create: `internal/autosync/emitter.go`
- Modify: `internal/hooks/handler.go`
- Modify: `internal/hooks/pre_tool_use.go`
- Modify: `internal/hooks/post_tool_use.go`
- Modify: `internal/hooks/stop.go`
- Modify: `internal/cli/root.go`
- Modify: `cmd/mp-devx/main.go`
- Create: `internal/autosync/emitter_test.go`

**Step 1: Write the failing tests**

```go
func TestHookEvents_ArePersistedToRawStream(t *testing.T) {}
func TestCLICommandEvents_ArePersistedToRawStream(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync ./internal/hooks ./internal/cli -run Test.*RawStream -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- `EmitRawEvent(projectDir, source, action, payload)`
- Call on hook enter/exit and CLI pre/post
- Persist policy operation events (`sync/apply/ignore/rollback/revoke`)

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync ./internal/hooks ./internal/cli -run Test.*RawStream -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/emitter.go internal/autosync/emitter_test.go internal/hooks/handler.go internal/hooks/pre_tool_use.go internal/hooks/post_tool_use.go internal/hooks/stop.go internal/cli/root.go cmd/mp-devx/main.go
git commit -m "feat(autosync): wire multi-entry raw event emission across runtime paths"
```

---

### Task 4: Build Detector Registry and Universal Signal Mapping

**Why:** The system must learn from generic dimensions, not language- or project-type hardcoding.

**What:** Create pluggable detector registry and baseline detectors for universal signal dimensions.

**How:** Define detector interface, registration, execution ordering, and normalized signal outputs.

**Key Design:** Core engine depends only on detector outputs (`dimension`, `value`, `confidence`, `evidence_refs`).

**Files:**
- Create: `internal/autosync/detector/types.go`
- Create: `internal/autosync/detector/registry.go`
- Create: `internal/autosync/detector/builtin.go`
- Create: `internal/autosync/detector/registry_test.go`

**Step 1: Write the failing tests**

```go
func TestRegistry_ExecutesBuiltinsDeterministically(t *testing.T) {}
func TestBuiltinDetectors_EmitUniversalDimensions(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/detector -run TestRegistry_ -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- Register builtins: `branching`, `workspace`, `command_contract`, `risk_profile`
- deterministic sorted execution
- merge signal list with stable ordering

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/detector -run TestRegistry_ -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/detector/types.go internal/autosync/detector/registry.go internal/autosync/detector/builtin.go internal/autosync/detector/registry_test.go
git commit -m "feat(autosync): add detector registry and universal signal builtins"
```

---

### Task 5: Implement Scoring Engine and Proposal Lifecycle

**Why:** Automatic policy activation must be confidence-governed and reversible.

**What:** Add scoring engine and proposal state transitions, including shadow/advisory/auto-candidate.

**How:** Compute support/conflict/sessions/stability/time-decay and evaluate activation thresholds.

**Key Design:** Safety-critical dimensions never go to auto-applied; they are forced into `manual-required`.

**Files:**
- Create: `internal/autosync/proposal/engine.go`
- Create: `internal/autosync/proposal/state.go`
- Create: `internal/autosync/proposal/store.go`
- Create: `internal/autosync/proposal/engine_test.go`

**Step 1: Write the failing tests**

```go
func TestProposal_AutoCandidateGate(t *testing.T) {}
func TestProposal_SafetyCriticalGoesManualRequired(t *testing.T) {}
func TestProposal_ConflictDemotesToShadow(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/proposal -run TestProposal_ -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- gate defaults:
  - `support>=8`
  - `sessions>=3`
  - `conflict_rate<15%`
  - `confidence>=0.95`
- state machine transitions and persistence to `proposals.jsonl`

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/proposal -run TestProposal_ -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/proposal/engine.go internal/autosync/proposal/state.go internal/autosync/proposal/store.go internal/autosync/proposal/engine_test.go
git commit -m "feat(autosync): add proposal scoring and lifecycle state machine"
```

---

### Task 6: Add Overlay Activation, Deny Confirmation, Revoke, Cooldown, and Data Reset

**Why:** User trust requires immediate enforceability of corrections and hard reset on policy rejection.

**What:** Implement deny-confirmation workflow (`delete` vs `skip-this-session`), overlay writer, applied ledger, revoke workflow, cooldown enforcement, and learning-data purge.

**How:** Add policy-denial decision handler that asks user choice, then branches into delete flow or session-only suppression flow.

**Key Design:** `delete` performs full revoke/reset path; `skip-this-session` only suppresses current-session injection and keeps persisted learning data untouched.

**Files:**
- Create: `internal/autosync/overlay/manager.go`
- Create: `internal/autosync/overlay/cooldown.go`
- Create: `internal/autosync/overlay/session_suppress.go`
- Create: `internal/autosync/overlay/manager_test.go`
- Modify: `internal/policy/resolve.go`

**Step 1: Write the failing tests**

```go
func TestOverlay_AutoApplyWritesAtomicFile(t *testing.T) {}
func TestOverlay_DenyAsksDeleteOrSkipSession(t *testing.T) {}
func TestOverlay_RevokeDeletesLearningDataAndSetsCooldown(t *testing.T) {}
func TestOverlay_SkipSessionSuppressesInjectionOnly(t *testing.T) {}
func TestResolver_ExcludesRevokedOrCoolingRules(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/overlay ./internal/policy -run TestOverlay_|TestResolver_ExcludesRevokedOrCoolingRules -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- `ApplyProposal` -> writes `overlays.auto.json` atomically
- deny path asks question and requires explicit choice:
  - `delete` -> remove overlay entry + purge related data + set cooldown + append audit
  - `skip-this-session` -> write session suppression only (no purge, no global delete)
- resolver merge precedence includes auto-learned overlay above session override

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/overlay ./internal/policy -run TestOverlay_|TestResolver_ExcludesRevokedOrCoolingRules -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/overlay/manager.go internal/autosync/overlay/cooldown.go internal/autosync/overlay/session_suppress.go internal/autosync/overlay/manager_test.go internal/policy/resolve.go
git commit -m "feat(autosync): add deny confirmation with delete or session-skip flows"
```

---

### Task 7: Implement `init-fingerprint` with Markdown Evidence Probing

**Why:** Cold-start quality depends on deterministic static evidence before runtime events exist.

**What:** Build doc-aware fingerprint scanner with evidence map and confidence outputs.

**How:** Parse required docs and common docs patterns, plus config/layout signals.

**Key Design:** Evidence precedence is fixed (`runtime > config > docs`) and docs are treated as prior hints.

**Files:**
- Create: `internal/autosync/fingerprint/scanner.go`
- Create: `internal/autosync/fingerprint/docs_probe.go`
- Create: `internal/autosync/fingerprint/scanner_test.go`
- Modify: `cmd/mp-devx/main.go`
- Modify: `cmd/mp-devx/main_test.go`

**Step 1: Write the failing tests**

```go
func TestFingerprint_ProbesRequiredDocs(t *testing.T) {}
func TestFingerprint_OutputIncludesEvidenceMapAndConfidence(t *testing.T) {}
func TestMPDevx_InitFingerprintAction(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/fingerprint ./cmd/mp-devx -run TestFingerprint_|TestMPDevx_InitFingerprintAction -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- probe files:
  - `README.md`, `CLAUDE.md`, `AGENTS.md`, `PRODUCT.md`
  - `ARCHITECTURE.md`, `CONTRIBUTING.md`, `docs/**/product*.md`, `docs/**/tech-stack*.md`, `docs/**/getting-started*.md`
- output vector + confidence + evidence map
- wire `mp-devx --action init-fingerprint`

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/fingerprint ./cmd/mp-devx -run TestFingerprint_|TestMPDevx_InitFingerprintAction -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/fingerprint/scanner.go internal/autosync/fingerprint/docs_probe.go internal/autosync/fingerprint/scanner_test.go cmd/mp-devx/main.go cmd/mp-devx/main_test.go
git commit -m "feat(autosync): add deterministic init-fingerprint with markdown evidence probing"
```

---

### Task 8: Add Cross-Project Semantic Preference Store and Migration Guards

**Why:** Reuse should improve cold start without leaking project-sensitive data or causing negative transfer.

**What:** Implement local-only desensitized semantic store with similarity threshold, cooldown, and rollback support.

**How:** Store only normalized preference patterns and aggregate scores with TTL.

**Key Design:** Never persist project private raw values in global semantic store.

**Files:**
- Create: `internal/autosync/semantic/store.go`
- Create: `internal/autosync/semantic/migrate.go`
- Create: `internal/autosync/semantic/store_test.go`

**Step 1: Write the failing tests**

```go
func TestSemanticStore_StoresPreferencePatternsOnly(t *testing.T) {}
func TestSemanticMigration_RequiresSimilarityThreshold(t *testing.T) {}
func TestSemanticMigration_CanRollbackOnConflictSpike(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/semantic -run TestSemantic -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- global store path: `~/.multipowers/policy/autosync/global.semantic.json`
- desensitize fields before write
- migration returns `shadow_only` when confidence/risk gate not met

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/semantic -run TestSemantic -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/semantic/store.go internal/autosync/semantic/migrate.go internal/autosync/semantic/store_test.go
git commit -m "feat(autosync): add cross-project semantic preference store with migration guardrails"
```

---

### Task 9: Build PolicyContext Snapshot and Prompt Injection (Including External Tools)

**Why:** Rules must be active in every run, otherwise learned policy has no runtime effect.

**What:** Generate deterministic `PolicyContext` and inject it into all prompt paths.

**How:** Integrate injection in mp workflow path and external CLI dispatch adapters.

**Key Design:** Injection references active rule IDs and excludes revoked/cooldown/session-suppressed rules; failure is non-blocking but auditable.

**Files:**
- Create: `internal/autosync/context/snapshot.go`
- Create: `internal/autosync/context/injector.go`
- Create: `internal/autosync/context/injector_test.go`
- Modify: `internal/hooks/handler.go`
- Modify: `internal/policy/dispatch.go`

**Step 1: Write the failing tests**

```go
func TestPolicyContext_IncludesActiveRulesWithReferences(t *testing.T) {}
func TestPolicyContext_ExcludesRevokedRules(t *testing.T) {}
func TestPolicyContext_ExcludesSessionSuppressedRules(t *testing.T) {}
func TestDispatchExternal_IncludesPolicyContextForPrompt(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/context ./internal/policy -run TestPolicyContext_|TestDispatchExternal_IncludesPolicyContextForPrompt -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- build deterministic context snapshot file
- inject into `/mp:*` prompt metadata
- for external command template support `{policy_context}` placeholder and fallback env/context-file adapter
- increment per-rule reference count on injection

**Step 4: Run tests to verify pass**

Run: `go test ./internal/autosync/context ./internal/policy -run TestPolicyContext_|TestDispatchExternal_IncludesPolicyContextForPrompt -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/context/snapshot.go internal/autosync/context/injector.go internal/autosync/context/injector_test.go internal/hooks/handler.go internal/policy/dispatch.go
git commit -m "feat(autosync): inject deterministic policy context into mp and external tool prompt paths"
```

---

### Task 10: Implement `mp policy` Commands (`sync/stats/gc/tune`)

**Why:** Operators need transparent control and observability without exposing low-level internals to users.

**What:** Add CLI surface for policy sync, apply/ignore/rollback/revoke, stats, GC, and tuning modes.

**How:** Extend `internal/cli/root.go` with `policy` subcommands and JSON/human outputs.

**Key Design:** Default `sync` is dry-run; auto-apply happens only for eligible non-safety proposals.

**Files:**
- Create: `internal/autosync/ops/service.go`
- Create: `internal/autosync/ops/service_test.go`
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`
- Modify: `docs/CLI-REFERENCE.md`

**Step 1: Write the failing tests**

```go
func TestPolicySync_DefaultDryRun(t *testing.T) {}
func TestPolicySync_ApplyIgnoreRollbackRevoke(t *testing.T) {}
func TestPolicyStatsAndTune(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/cli ./internal/autosync/ops -run TestPolicy -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- command contract:
  - `mp policy sync`
  - `mp policy sync --apply`
  - `mp policy sync --ignore <id>`
  - `mp policy sync --rollback <id>`
  - `mp policy sync --revoke <id>`
- `mp policy stats`
- `mp policy gc`
- `mp policy tune --mode balanced|accuracy|storage`

**Step 4: Run tests to verify pass**

Run: `go test ./internal/cli ./internal/autosync/ops -run TestPolicy -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/ops/service.go internal/autosync/ops/service_test.go internal/cli/root.go internal/cli/root_test.go docs/CLI-REFERENCE.md
git commit -m "feat(cli): add mp policy sync/stats/gc/tune command surface"
```

---

### Task 11: Integrate Hook and Doctor Governance for Drift/Unresolved High-Confidence Proposals

**Why:** Autonomous learning needs guardrails and drift alarms to remain trustworthy.

**What:** Add governance checks and hook-level enforcement/warnings for autosync state.

**How:** Extend doctor checks and hook metadata to include autosync risk signals.

**Key Design:** Hard guardrails block, while advisory issues warn and preserve user flow.

**Files:**
- Modify: `internal/doctor/registry.go`
- Modify: `internal/doctor/checks_local.go`
- Modify: `internal/doctor/checks_local_test.go`
- Modify: `internal/hooks/pre_tool_use.go`
- Modify: `internal/hooks/pre_tool_use_test.go`

**Step 1: Write the failing tests**

```go
func TestDoctor_AutoSyncDriftWarns(t *testing.T) {}
func TestDoctor_UnhandledHighConfidenceProposalWarns(t *testing.T) {}
func TestHook_HardSafetyRuleStillBlocks(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/doctor ./internal/hooks -run TestDoctor_AutoSync|TestHook_HardSafetyRuleStillBlocks -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- new doctor checks:
  - autosync-drift
  - autosync-unresolved-high-confidence
- hook enforcement:
  - block on safety-critical violations
  - attach autosync metadata/warnings in allow path

**Step 4: Run tests to verify pass**

Run: `go test ./internal/doctor ./internal/hooks -run TestDoctor_AutoSync|TestHook_HardSafetyRuleStillBlocks -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/doctor/registry.go internal/doctor/checks_local.go internal/doctor/checks_local_test.go internal/hooks/pre_tool_use.go internal/hooks/pre_tool_use_test.go
git commit -m "feat(governance): add autosync drift and unresolved proposal guardrails"
```

---

### Task 12: End-to-End Verification, Quota GC (`LRU + RefCount`), and Documentation Finalization

**Why:** Full-stack verification and operator docs are required before claiming readiness.

**What:** Add e2e coverage for learning/apply/revoke/inject, implement quota cleanup strategy, and finalize architecture docs.

**How:** Add integration tests and update docs with operational playbook and failure handling.

**Key Design:** GC chooses eviction by `least-recently-used + lowest cumulative references` score; high-reference data survives longer.

**Files:**
- Create: `internal/autosync/gc/collector.go`
- Create: `internal/autosync/gc/collector_test.go`
- Create: `internal/autosync/e2e_test.go`
- Modify: `docs/ARCHITECTURE.md`
- Modify: `docs/PLUGIN-ARCHITECTURE.md`
- Modify: `docs/COMMAND-REFERENCE.md`
- Modify: `docs/multipowers/README.md`

**Step 1: Write the failing tests**

```go
func TestGC_UsesLRUPlusReferenceCount(t *testing.T) {}
func TestE2E_AutoApplyInjectRevokeResetFlow(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/autosync/gc ./internal/autosync -run TestGC_|TestE2E_AutoApplyInjectRevokeResetFlow -v`  
Expected: FAIL

**Step 3: Write minimal implementation**

- GC scoring:
  - recency rank weight
  - cumulative reference count weight
  - tie-break by oldest timestamp
- e2e flow:
  - raw events -> proposal -> auto-apply -> prompt injection
  - user deny -> askquestion -> choose `skip-this-session` -> current-session injection suppression only
  - user deny -> askquestion -> choose `delete` -> revoke + cooldown -> zero-baseline relearn
- docs update with operational commands and troubleshooting

**Step 4: Run verification suite**

Run:

```bash
go test ./internal/autosync/... -v
go test ./internal/policy ./internal/hooks ./internal/doctor ./internal/cli -v
go test ./cmd/mp ./cmd/mp-devx -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/autosync/gc/collector.go internal/autosync/gc/collector_test.go internal/autosync/e2e_test.go docs/ARCHITECTURE.md docs/PLUGIN-ARCHITECTURE.md docs/COMMAND-REFERENCE.md docs/multipowers/README.md
git commit -m "feat(autosync): complete e2e learning loop, lru-refcount gc, and docs rollout"
```

---

## Final Verification Checklist (Must Be Executed Before Merge)

1. `go test ./... -v`
2. `go test -race ./internal/autosync/...`
3. `mp-devx --action doctor --dir . --verbose`
4. Manual smoke:
   - run `/mp:*` command and verify policy context injection metadata,
   - run an external vibe coding tool path and verify injection adapter output,
   - deny one active policy and verify askquestion branch:
     - `skip-this-session` suppresses only this session
     - `delete` removes rule and writes cooldown metadata.

## Rollback Plan

1. Disable autosync activation via feature flag in runtime settings (temporary kill switch).
2. Keep read-only ingestion enabled for observability while deactivating overlay writes.
3. Restore previous resolver merge behavior (ignore overlays) if safety incidents occur.
4. Retain audit logs (`applied.jsonl`) for incident analysis.
