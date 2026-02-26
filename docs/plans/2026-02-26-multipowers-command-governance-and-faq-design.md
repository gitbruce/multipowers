# Multipowers Command Governance & Auto-FAQ Design

## 1. Goal
Build a low-conflict, upstream-friendly governance layer for `/mp:*` commands that improves execution reliability and prevents repeated failures in target projects.

This design extends existing multipowers customizations with:
- standardized preflight checks
- provider degradation policy
- strict artifact/output boundaries
- runtime precondition contract enforcement
- automatic FAQ generation and deduped refinement
- a stable target-project `/.multipowers/CLAUDE.md` referencing auto-generated `/.multipowers/FAQ.md`

## 2. Design Scope
In scope:
- command contract and preflight governance for spec-driven commands
- target-project context and artifact rules
- error learning loop (auto FAQ)
- doc structure for tool-project and target-project users

Out of scope:
- full provider SDK rewrites
- changing upstream main branch architecture
- manual FAQ curation workflows

## 3. Architecture Principles
- Main branch remains upstream-aligned with minimum divergence.
- Multipowers branch carries custom behavior mostly under `custom/*` and thin hooks in high-churn files.
- All target-project outputs must stay under `/<target>/.multipowers/`.
- No writes to `$HOME` or tool project when operating on target projects.
- Spec-driven command behavior is consistent and centrally enforced.

## 4. Borrowed Concepts and Mapping (from referenced guide)
1. Instruction hierarchy and boundaries -> explicit command contract and precedence model.
2. Standardized command preflight -> one reusable preflight pipeline for spec-driven commands.
3. Provider fallback handling -> continue with remaining providers if quorum remains.
4. Runtime contract -> mandatory pre-run command application from runtime config.
5. Artifact isolation -> hard output boundary under target `/.multipowers/`.
6. Wizarded context bootstrap -> `/mp:init` setup-first behavior.
7. Observability -> structured failure events for diagnosis and FAQ extraction.
8. Dual-view docs -> separate tool maintenance docs and target usage docs.
9. Anti-regression knowledge loop -> `CLAUDE.md` + auto-maintained `FAQ.md`.

## 5. Command-Level Governance Design
### 5.1 Spec-driven command set
Spec-driven commands include at least:
- `/mp:plan`
- `/mp:discover`
- `/mp:define`
- `/mp:develop`
- `/mp:deliver`
- `/mp:embrace`
- `/mp:research`

### 5.2 Mandatory Step 0 preflight
Before command-specific logic:
1. Resolve project root (target project cwd).
2. Verify required context files exist in `/.multipowers/`:
   - `product.md`
   - `product-guidelines.md`
   - `tech-stack.md`
   - `workflow.md`
   - `tracks.md`
   - `CLAUDE.md`
3. If missing, force `/mp:init` wizard.
4. Re-check context; fail-fast if still incomplete.
5. If `/.multipowers/context/runtime.json` exists, load runtime preconditions.
6. Execute pre-run commands per fail-fast policy.
7. Proceed to command body only after Step 0 success.

## 6. Provider Orchestration and Fallback Rules
- Debate-style multi-LLM workflows should attempt up to 3 providers (Claude + Codex + Gemini).
- Minimum viable quorum is 2 providers.
- If one provider fails, continue with remaining two.
- If available providers drop below 2, fail with actionable diagnostics.
- Proxy routing for Codex/Gemini must be applied consistently at every invocation path.
- Host for proxy is dynamically detected (not hardcoded to `127.0.0.1`).

## 7. Runtime Preconditions Contract
- Runtime config source: `/.multipowers/context/runtime.json`.
- Used by all commands/skills that invoke tools/providers.
- Supports user-defined pre-run command lists (language-agnostic, not Python-specific).
- Policy: `fail-fast` (confirmed).
- Guard requirement: runtime file is optional for context readiness; when present, its pre-run policy is enforced.
- If any pre-run command fails, stop immediately and report which step failed.

## 8. CLAUDE.md and FAQ.md Design
### 8.1 Template ownership and generation
Tool project templates:
- `custom/templates/CLAUDE.md`
- `custom/templates/FAQ.md`

Generated in target project by `/mp:init`:
- `/.multipowers/CLAUDE.md`
- `/.multipowers/FAQ.md`

### 8.2 CLAUDE.md role
- Stable, human-readable project contract.
- Cross-language baseline conventions.
- Explicit link/reference to `/.multipowers/FAQ.md` for learned failure-avoidance patterns.

### 8.3 FAQ.md role
- Fully auto-generated and auto-maintained.
- Error-type-based sections (not command-based), e.g.:
  - proxy/network
  - timeout
  - auth/permission
  - missing context
  - model unavailable
  - path boundary violations
  - provider capacity

### 8.4 FAQ update behavior (no manual maintenance)
Trigger points:
- command non-zero exit
- timeout
- retry threshold reached
- quality gate failure
- task completion summarization

Update algorithm:
1. Ingest new failure events.
2. Normalize and classify by error type.
3. Deduplicate by key:
   - `error_type + normalized_root_cause + normalized_fix`
4. Refine wording into concise entries.
5. Rewrite `/.multipowers/FAQ.md` (regenerate style, not append-only).
6. Enforce max entry cap (e.g., top 120 by recency and frequency).

Explicit requirements:
- No “other/low-frequency merge bucket”.
- No monthly archive.
- No backup file generation.

## 9. Data Flow
1. User runs `/mp:*`.
2. Step 0 preflight validates context and runtime preconditions.
3. Command executes provider/tool logic.
4. Errors and retries emit structured failure events.
5. FAQ synthesizer deduplicates/refines and rewrites `/.multipowers/FAQ.md`.
6. Future runs read `CLAUDE.md` + `FAQ.md` as anti-regression guidance.

## 10. Error Handling
- Context missing after init: hard fail, no downstream execution.
- Runtime pre-run failure: hard fail (`fail-fast`).
- Provider partial failure (3->2): continue.
- Provider collapse (<2) for debate/multi-LLM: fail with remediation tips.
- Write-boundary violation attempt: block and report.

## 11. Testing Strategy
Contract tests:
- Step 0 runs before spec-driven command body.
- Missing context always forces init path.
- Runtime preconditions are always loaded and executed.
- Proxy env is applied in all Codex/Gemini paths.
- Debate quorum behavior works (3->2 continue, <2 fail).
- No writes to `$HOME` or tool project when target project execution is active.
- FAQ generation dedups and rewrites deterministically.

Regression tests:
- Ensure `.claude/commands/*` remain close to upstream except minimal required guard text.
- Ensure `main` branch remains upstream-syncable with minimal conflict footprint.

## 12. Documentation Structure
Tool project docs (maintenance):
- customization architecture
- sync/merge discipline
- command contracts
- FAQ generation internals

Target project docs (usage):
- setup/init
- runtime preconditions
- command behavior and expectations
- troubleshooting linked to FAQ categories

## 13. Proposed File-Level Changes (Design-level)
- Add templates:
  - `custom/templates/CLAUDE.md`
  - `custom/templates/FAQ.md`
- Update init/config docs to include template generation and FAQ loop.
- Update target-project docs to explain automatic FAQ lifecycle.
- Add/extend tests for preflight, runtime contract, provider quorum, FAQ dedup/regeneration.

## 14. Acceptance Criteria
- `/mp:init` creates `/.multipowers/CLAUDE.md` and `/.multipowers/FAQ.md` from templates.
- Spec-driven commands enforce Step 0 preflight before command-specific workflow.
- Step 0 required context files are `product.md`, `product-guidelines.md`, `tech-stack.md`, `workflow.md`, `tracks.md`, `CLAUDE.md`.
- Runtime preconditions from `runtime.json` are optional to presence, mandatory to enforce when configured, and fail-fast.
- Multi-LLM debate uses up to 3 providers, continues with 2, fails below 2.
- All artifacts/temp outputs remain under target `/.multipowers/`.
- FAQ auto-updates, dedups, refines, and stays bounded without manual edits.
- Documentation updated for both tool maintainers and target users.

## 15. Rollout Plan (High-level)
1. Implement template + init generation.
2. Consolidate Step 0 guard in command contracts and orchestration path.
3. Enforce runtime precondition loading across all execution paths.
4. Implement provider quorum fallback.
5. Implement FAQ event schema and synthesizer.
6. Add contract/regression tests.
7. Update docs and verify behavior in a clean target project.
