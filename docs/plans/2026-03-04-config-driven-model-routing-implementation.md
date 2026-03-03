# Config-Driven Model Routing Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Move workflow/persona model selection and executor dispatch to config-driven policy (`config/*.yaml`), compile it into `.claude-plugin/runtime/policy.json`, and make runtime resolve models/executors without hardcoded model strings in command/skill flows.

**Architecture:** Introduce a new `internal/policy` pipeline: load YAML source config (`workflows.yaml`, `agents.yaml`, `executors.yaml`) -> validate schema/semantics -> compile normalized runtime policy JSON. Runtime paths (hooks + CLI workflow/persona entrypoints) resolve execution contracts from compiled policy. External executors use hard enforcement with one automatic fallback; Claude Code path uses hint enforcement.

**Tech Stack:** Go (new `internal/policy` package, `cmd/mp-devx` actions, `internal/hooks`, `internal/cli`, `internal/workflows`), YAML/JSON configs in `config/`, runtime artifacts in `.claude-plugin/runtime/`, Go tests (`go test ./...`).

---

## Task Status Tracker (MUST UPDATE AFTER EVERY TASK/SUBTASK)

| ID | Name | Status |
|---|---|---|
| T01 | Define policy schemas and source config files | COMPLETED |
| T01-S01 | Workflow config with optional task level | COMPLETED |
| T01-S02 | Agent + executor config base | COMPLETED |
| T02 | Build parser and semantic validator (`internal/policy`) | COMPLETED |
| T02-S01 | YAML loader + schema validation tests | COMPLETED |
| T02-S02 | Semantic checks (fallback, executor mapping) | COMPLETED |
| T03 | Compile runtime policy JSON artifact | COMPLETED |
| T03-S01 | Compiler output model + golden tests | COMPLETED |
| T03-S02 | Build action integration in `mp-devx` | COMPLETED |
| T04 | Integrate runtime resolver into hooks and workflow entrypoints | COMPLETED |
| T04-S01 | Resolver contract + workflow task resolution | COMPLETED |
| T04-S02 | Hook metadata + persona/workflow consumption | COMPLETED |
| T05 | Implement executor dispatch and one-hop fallback | COMPLETED |
| T05-S01 | External hard enforcement with model arg injection | COMPLETED |
| T05-S02 | Automatic single-hop fallback (cross-provider allowed) | COMPLETED |
| T06 | Add `/mp:config` visibility toggle (default show) | COMPLETED |
| T06-S01 | Persist `show_model_routing` setting | COMPLETED |
| T06-S02 | Respect toggle in output rendering | COMPLETED |
| T07 | Replace hardcoded model paths and add guardrails | NOT_STARTED |
| T07-S01 | Remove runtime hardcoding dependency points | NOT_STARTED |
| T07-S02 | Add hardcoded-model scan tests | NOT_STARTED |
| T08 | Build/release integration for runtime-only `.claude-plugin` | NOT_STARTED |
| T08-S01 | Build binaries + policy in one pipeline | NOT_STARTED |
| T08-S02 | Runtime read-only contract checks | NOT_STARTED |
| T09 | End-to-end verification and migration docs | NOT_STARTED |
| T09-S01 | E2E tests for resolver/enforcement/fallback | NOT_STARTED |
| T09-S02 | Operator docs + evidence capture | NOT_STARTED |

### Status Update Rule (Mandatory)

After every subtask and every task:

1. Update this table row from `NOT_STARTED` -> `IN_PROGRESS` -> `COMPLETED`.
2. Persist machine-readable status in state metrics:

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.<task_or_subtask_id> --value <IN_PROGRESS|COMPLETED> --json
```

Expected: JSON response status is `ok`.

---

### Task T01: Define Policy Schemas And Source Config Files

**Why:** Without a stable source schema, resolver behavior remains ambiguous and junior engineers will create incompatible config.

**What:** Introduce three source config files in `config/` and define field contracts for workflow/task policy, agent policy, and executor mapping.

**How:** Start with failing parser tests that expect required fields and optional task-level workflow overrides; then create minimum valid YAML files.

**Key Design Considerations:**
- `workflows.yaml` must support 2 levels: workflow default + optional `tasks`.
- Task names can be semantic names or `task_1`, `task_2`, ...
- Executor mapping is centralized in `executors.yaml`; workflow/agent files must not carry command templates.

**Files:**
- Create: `config/workflows.yaml`
- Create: `config/agents.yaml`
- Create: `config/executors.yaml`
- Create: `internal/policy/schema_test.go`

#### Subtask T01-S01: Workflow config with optional task level

**Why:** User explicitly requires per-workflow multi-model task routing.

**What:** Define `workflows.yaml` with `default` + optional `tasks` sections.

**How:** Write tests first, then create minimal config that passes.

**Key Design Considerations:** Task-level policy overrides workflow default, and omitted task falls back safely.

**Step 1: Write failing tests**

Add test cases in `internal/policy/schema_test.go`:
- valid workflow with `default` only
- valid workflow with `tasks.task_1`
- invalid workflow missing `model` and `executor_profile`

**Step 2: Run tests to verify failure**

Run: `go test ./internal/policy -run TestWorkflowSchema -count=1`
Expected: FAIL because parser/types do not exist yet.

**Step 3: Create `config/workflows.yaml` minimal contract**

Example seed:

```yaml
version: "1"
workflows:
  define:
    default:
      model: gpt-5.3-codex
      executor_profile: codex_cli
      fallback_policy: cross_provider_once
    tasks:
      task_1:
        model: gpt-5.3-codex
        executor_profile: codex_cli
      task_2:
        model: gemini-3-pro-preview
        executor_profile: gemini_cli
```

**Step 4: Update status**

Run:

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T01-S01 --value COMPLETED --json
```

#### Subtask T01-S02: Agent + executor config base

**Why:** Agent model policy and executor routing must be independently configurable.

**What:** Add `agents.yaml` and `executors.yaml` initial definitions.

**How:** Define explicit mapping for current key models and fallback behavior.

**Key Design Considerations:**
- `claude_code` uses `enforcement: hint`.
- external executors use `enforcement: hard`.
- fallback supports cross-provider, one hop only.

**Step 1: Create `config/agents.yaml`**

```yaml
version: "1"
agents:
  backend-architect:
    model: gpt-5.3-codex
    executor_profile: codex_cli
    fallback_policy: cross_provider_once
  security-auditor:
    model: claude-opus-4.6
    executor_profile: claude_code
    fallback_policy: none
```

**Step 2: Create `config/executors.yaml`**

```yaml
version: "1"
executors:
  codex_cli:
    kind: external_cli
    command_template: ["codex", "exec", "-m", "{model}", "-C", "{project_dir}", "{prompt}"]
    enforcement: hard
  gemini_cli:
    kind: external_cli
    command_template: ["gemini", "-m", "{model}", "-p", "{prompt}"]
    enforcement: hard
  claude_code:
    kind: claude_code
    enforcement: hint
fallback_policies:
  cross_provider_once:
    max_hops: 1
    chain:
      - from: gpt-5.3-codex
        to: gemini-3-pro-preview
      - from: gemini-3-pro-preview
        to: claude-sonnet-4.5
```

**Step 3: Update task status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T01-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T01 --value COMPLETED --json
```

**Step 4: Commit**

```bash
git add config/workflows.yaml config/agents.yaml config/executors.yaml internal/policy/schema_test.go
git commit -m "feat(policy): add source config schemas for workflows agents executors"
```

---

### Task T02: Build Parser And Semantic Validator (`internal/policy`)

**Why:** Config-driven architecture fails without strict semantic validation; runtime should never discover broken config late.

**What:** Implement YAML loaders, typed structs, and semantic validators.

**How:** TDD: parser tests -> semantic tests -> minimal implementation.

**Key Design Considerations:**
- parser error messages must include file path + key path.
- semantic checks must be deterministic and order-independent.

**Files:**
- Create: `internal/policy/types.go`
- Create: `internal/policy/load.go`
- Create: `internal/policy/validate.go`
- Create: `internal/policy/load_test.go`
- Modify: `go.mod` (if yaml dependency needed)

#### Subtask T02-S01: YAML loader + schema validation tests

**Why:** Junior engineers need immediate feedback on malformed config.

**What:** Parse three yaml files into typed structures.

**How:** Add table tests for success/failure cases.

**Key Design Considerations:** return structured errors (`file`, `field`, `reason`).

**Step 1: Write failing tests in `load_test.go`**

Cover:
- missing `version`
- unknown executor profile in workflow entry
- valid parse path

**Step 2: Run to verify red**

Run: `go test ./internal/policy -run TestLoadSourceConfig -count=1`
Expected: FAIL.

**Step 3: Implement `types.go` + `load.go`**

- Define `SourceConfig` aggregate.
- Add `LoadSourceConfig(root string)` that reads `config/workflows.yaml`, `config/agents.yaml`, `config/executors.yaml`.

**Step 4: Run test to green**

Run: `go test ./internal/policy -run TestLoadSourceConfig -count=1`
Expected: PASS.

**Step 5: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T02-S01 --value COMPLETED --json
```

#### Subtask T02-S02: Semantic checks (fallback, executor mapping)

**Why:** Syntax-valid config can still be runtime-invalid.

**What:** Enforce semantic rules.

**How:** Add validator with explicit rule list + tests.

**Key Design Considerations:**
- fallback `max_hops` must be exactly `1` for this phase.
- every referenced `executor_profile` must exist.
- every model in fallback chain must resolve to an executor.

**Step 1: Write failing semantic tests**

Create tests:
- fallback chain with 2 hops => fail
- agent references missing executor profile => fail
- valid cross-provider single-hop => pass

**Step 2: Implement `validate.go`**

Expose `ValidateSourceConfig(cfg SourceConfig) error`.

**Step 3: Run tests**

Run: `go test ./internal/policy -run TestValidateSourceConfig -count=1`
Expected: PASS.

**Step 4: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T02-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T02 --value COMPLETED --json
git add internal/policy/types.go internal/policy/load.go internal/policy/validate.go internal/policy/load_test.go
# include dependency file if changed
if git diff --name-only --cached | rg -q '^go\.mod$|^go\.sum$'; then git add go.mod go.sum; fi
git commit -m "feat(policy): add source config loader and semantic validator"
```

---

### Task T03: Compile Runtime Policy JSON Artifact

**Why:** Runtime must consume immutable compiled policy from `.claude-plugin/runtime/policy.json` to enforce dev/runtime separation.

**What:** Create compiler that normalizes source config into runtime policy.

**How:** Add compiler + golden tests + devx action.

**Key Design Considerations:**
- compiler output must be stable sorted order.
- include checksum metadata for drift detection.

**Files:**
- Create: `internal/policy/compile.go`
- Create: `internal/policy/runtime_policy.go`
- Create: `internal/policy/compile_test.go`
- Modify: `cmd/mp-devx/main.go`
- Modify: `cmd/mp-devx/main_test.go`

#### Subtask T03-S01: Compiler output model + golden tests

**Why:** Prevent accidental policy shape regressions.

**What:** Define runtime JSON schema and compile function.

**How:** Golden-file style tests.

**Key Design Considerations:** preserve source references for explainability.

**Pseudo-code (compiler):**

```text
load source cfg
validate cfg
for each workflow:
  add workflow.default to runtime index
  for each workflow.task:
    add task override to runtime index
for each agent:
  add agent policy to runtime index
for each model pattern / executor profile:
  build model->executor lookup
for each fallback rule:
  build one-hop fallback map
emit RuntimePolicy JSON with deterministic key order
```

**Step 1: Write failing compile tests**

Run: `go test ./internal/policy -run TestCompileRuntimePolicy -count=1`
Expected: FAIL.

**Step 2: Implement compiler and runtime structs**

**Step 3: Generate expected test fixture**

Create fixture: `internal/policy/testdata/runtime_policy.golden.json`

**Step 4: Run tests**

Run: `go test ./internal/policy -run TestCompileRuntimePolicy -count=1`
Expected: PASS.

**Step 5: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T03-S01 --value COMPLETED --json
```

#### Subtask T03-S02: Build action integration in `mp-devx`

**Why:** Policy compilation must be part of standard build pipeline.

**What:** Add `-action build-policy` and `-action build-runtime`.

**How:** Reuse compiler package and write artifact to `.claude-plugin/runtime/policy.json`.

**Key Design Considerations:** build must fail fast if policy invalid.

**Step 1: Add failing tests in `cmd/mp-devx/main_test.go`**

- `build-policy` action success path
- invalid config returns non-zero

**Step 2: Implement actions in `cmd/mp-devx/main.go`**

- `build-policy` => compile and write JSON
- `build-runtime` => run `build-policy` then binaries build

**Step 3: Run tests**

Run: `go test ./cmd/mp-devx -run TestRun_ActionBuildPolicy -count=1`
Expected: PASS.

**Step 4: Run policy build command**

Run: `./scripts/mp-devx -action build-policy`
Expected: prints success and writes `.claude-plugin/runtime/policy.json`.

**Step 5: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T03-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T03 --value COMPLETED --json
git add internal/policy/compile.go internal/policy/runtime_policy.go internal/policy/compile_test.go internal/policy/testdata/runtime_policy.golden.json cmd/mp-devx/main.go cmd/mp-devx/main_test.go
git commit -m "feat(policy): compile runtime policy and expose mp-devx build actions"
```

---

### Task T04: Integrate Runtime Resolver Into Hooks And Workflow Entry Points

**Why:** Config is useless unless runtime actually resolves model/executor from compiled policy.

**What:** Add resolver API and wire it into hooks + workflow/persona paths.

**How:** implement resolver first, then consume in hooks and CLI/workflows.

**Key Design Considerations:** resolver should accept normalized input (`scope,name,task`) and return explainable contract.

**Files:**
- Create: `internal/policy/resolve.go`
- Create: `internal/policy/resolve_test.go`
- Modify: `internal/hooks/handler.go`
- Modify: `internal/hooks/session_start.go`
- Modify: `internal/workflows/persona.go`
- Modify: `internal/cli/root.go`

#### Subtask T04-S01: Resolver contract + workflow task resolution

**Why:** Need predictable precedence logic for workflow task overrides.

**What:** Implement `ResolveExecutionContract(...)`.

**How:** TDD with precedence tests.

**Key Design Considerations:** precedence must be `workflow.task` -> `workflow.default` -> global default.

**Pseudo-code (resolver):**

```text
input: scope, name, task(optional), projectDir
policy := loadCompiledPolicy(projectDir)
if scope == workflow:
  if task exists and policy has workflow.task:
    selected := workflow.task
  else if workflow.default exists:
    selected := workflow.default
  else error
if scope == agent:
  selected := agent policy
executor := mapModelToExecutor(selected.model)
fallback := lookupFallback(selected.model)
return contract(selected, executor, fallback)
```

**Step 1: Write failing resolver tests**

Run: `go test ./internal/policy -run TestResolveExecutionContract -count=1`
Expected: FAIL.

**Step 2: Implement resolver**

**Step 3: Run tests**

Run: `go test ./internal/policy -run TestResolveExecutionContract -count=1`
Expected: PASS.

**Step 4: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T04-S01 --value COMPLETED --json
```

#### Subtask T04-S02: Hook metadata + persona/workflow consumption

**Why:** The actual command path must stop relying on hardcoded model values.

**What:** Replace `modelroute` usage with compiled policy resolver metadata path.

**How:** Wire resolver in `UserPromptSubmit` and session start, then route persona/workflow entry calls through contract resolution.

**Key Design Considerations:** preserve backward-compatible metadata keys during migration window.

**Step 1: Write failing integration tests**

- update `internal/hooks/handler_test.go` expected metadata fields.
- add test for workflow task resolution fallback.

**Step 2: Implement hook and entrypoint wiring**

**Step 3: Run tests**

Run: `go test ./internal/hooks ./internal/workflows ./internal/cli -count=1`
Expected: PASS.

**Step 4: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T04-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T04 --value COMPLETED --json
git add internal/policy/resolve.go internal/policy/resolve_test.go internal/hooks/handler.go internal/hooks/session_start.go internal/workflows/persona.go internal/cli/root.go internal/hooks/handler_test.go
# include any newly created tests for workflows/cli if added
git add internal/workflows/*.go internal/cli/*.go 2>/dev/null || true
git commit -m "feat(runtime): resolve workflow and agent execution contracts from compiled policy"
```

---

### Task T05: Implement Executor Dispatch And One-Hop Fallback

**Why:** Enforcement mode and fallback behavior are core runtime guarantees requested by user.

**What:** Implement dispatch that differentiates `external_cli` hard mode and `claude_code` hint mode, with auto single-hop fallback.

**How:** Add dispatcher module with unit tests for success/failure/fallback.

**Key Design Considerations:** fallback triggers only once and only on external hard-failure.

**Files:**
- Create: `internal/policy/dispatch.go`
- Create: `internal/policy/dispatch_test.go`
- Modify: `internal/workflows/persona.go` (consume dispatcher)
- Modify: `internal/providers/*.go` (if model arg API changes required)

#### Subtask T05-S01: External hard enforcement with model arg injection

**Why:** Hard mode must guarantee selected model is passed to tool invocation.

**What:** Inject `{model}` into command template and validate generated argv.

**How:** test command rendering + execution path.

**Key Design Considerations:** never silently drop model arg in hard mode.

**Pseudo-code (dispatch hard mode):**

```text
if executor.kind == external_cli:
  argv := renderTemplate(executor.command_template, model, prompt, project_dir)
  run external process(argv)
  if exit==0: return success
  else if fallback exists: attempt fallback once
  else return hard-fail
```

**Step 1: Write failing tests for argv rendering and hard-fail**

Run: `go test ./internal/policy -run TestDispatchExternalHardMode -count=1`
Expected: FAIL.

**Step 2: Implement dispatch logic for hard mode**

**Step 3: Run tests**

Run: `go test ./internal/policy -run TestDispatchExternalHardMode -count=1`
Expected: PASS.

**Step 4: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T05-S01 --value COMPLETED --json
```

#### Subtask T05-S02: Automatic single-hop fallback (cross-provider allowed)

**Why:** Requested behavior requires transparent degradation path without user confirmation.

**What:** Add one-hop fallback with metadata (`degraded`, `fallback_from`, `fallback_to`).

**How:** test failure of primary model and success on fallback executor.

**Key Design Considerations:**
- max one fallback attempt
- fallback may be external or `claude_code`
- if fallback target is `claude_code`, emit hint contract and continue

**Pseudo-code (fallback):**

```text
res := dispatch(primary)
if res.success: return res
if !primary.is_external || !primary.has_fallback: return fail
fb := resolve(fallback_model)
res2 := dispatch(fb)
if res2.success:
  res2.degraded = true
  res2.fallback_from = primary.model
  res2.fallback_to = fb.model
  return res2
return fail
```

**Step 1: Write failing tests for fallback success and fallback fail**

Run: `go test ./internal/policy -run TestDispatchOneHopFallback -count=1`
Expected: FAIL.

**Step 2: Implement fallback execution path**

**Step 3: Run tests**

Run: `go test ./internal/policy -run TestDispatchOneHopFallback -count=1`
Expected: PASS.

**Step 4: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T05-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T05 --value COMPLETED --json
git add internal/policy/dispatch.go internal/policy/dispatch_test.go internal/workflows/persona.go internal/providers/*.go
git commit -m "feat(dispatch): enforce external model constraints and one-hop automatic fallback"
```

---

### Task T06: Add `/mp:config` Visibility Toggle (Default Show)

**Why:** User requires model/fallback visibility to be configurable and default visible.

**What:** Add CLI config command and persistent runtime setting.

**How:** Introduce settings storage + output-gating logic.

**Key Design Considerations:** default value must be true when setting missing.

**Files:**
- Create: `internal/settings/runtime_settings.go`
- Create: `internal/settings/runtime_settings_test.go`
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`

#### Subtask T06-S01: Persist `show_model_routing` setting

**Why:** Toggle must survive process restarts.

**What:** Store in project-local state file (or dedicated settings file).

**How:** add get/set helper with defaults.

**Key Design Considerations:** avoid breaking existing `tracks` state schema.

**Step 1: Write failing tests for default true and set false**

Run: `go test ./internal/settings -run TestShowModelRoutingSetting -count=1`
Expected: FAIL.

**Step 2: Implement settings helper**

Default behavior:
- missing key => `true`

**Step 3: Run tests**

Run: `go test ./internal/settings -run TestShowModelRoutingSetting -count=1`
Expected: PASS.

**Step 4: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T06-S01 --value COMPLETED --json
```

#### Subtask T06-S02: Respect toggle in output rendering

**Why:** Config is meaningless unless output respects it.

**What:** Add `/mp config` or `/mp:config` mapping and gate routing display.

**How:** extend CLI command handling and response formatting.

**Key Design Considerations:** keep JSON contract stable; hide only user-facing text when disabled.

**Step 1: Write failing CLI tests in `root_test.go`**

Cases:
- `mp config show-model-routing off`
- response hides `requested_model/effective_model` text fields

**Step 2: Implement CLI command branch**

**Step 3: Run tests**

Run: `go test ./internal/cli -run TestConfigShowModelRouting -count=1`
Expected: PASS.

**Step 4: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T06-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T06 --value COMPLETED --json
git add internal/settings/runtime_settings.go internal/settings/runtime_settings_test.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(config): add show_model_routing toggle with default visible"
```

---

### Task T07: Replace Hardcoded Model Paths And Add Guardrails

**Why:** Long-term drift prevention requires enforcement, not convention.

**What:** Remove direct runtime hardcoding usage and add scanning tests.

**How:** deprecate old route defaults from runtime decision path, add guard tests.

**Key Design Considerations:** keep old code readable during migration but ensure it is not used for active resolution.

**Files:**
- Modify: `internal/modelroute/route.go` (mark deprecated or remove from active path)
- Create: `internal/validation/model_hardcode_guard_test.go`
- Modify: `internal/workflows/persona_test.go` (avoid brittle fixed model assertions where inappropriate)

#### Subtask T07-S01: Remove runtime hardcoding dependency points

**Why:** Existing defaults in code conflict with config-only objective.

**What:** Ensure runtime resolution uses compiled policy path only.

**How:** update call sites; keep old helper only for compatibility warning if needed.

**Key Design Considerations:** don't break hooks metadata contract abruptly; map old fields to new values.

**Step 1: Add failing tests asserting resolver source path is compiled policy**

Run: `go test ./internal/hooks ./internal/policy -run TestResolverUsesCompiledPolicy -count=1`
Expected: FAIL.

**Step 2: Update call sites and tests**

**Step 3: Run tests**

Run: `go test ./internal/hooks ./internal/policy -count=1`
Expected: PASS.

**Step 4: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T07-S01 --value COMPLETED --json
```

#### Subtask T07-S02: Add hardcoded-model scan tests

**Why:** Prevent new hardcoded model regressions outside config fixtures.

**What:** Add test scanning repository for model strings outside allowed paths.

**How:** use allowlist-driven test.

**Key Design Considerations:** allow model strings in:
- `config/*.yaml`
- explicit test fixtures
- changelog/release notes docs

**Step 1: Implement guard test**

Create: `internal/validation/model_hardcode_guard_test.go`

**Step 2: Run tests**

Run: `go test ./internal/validation -run TestNoHardcodedModelsOutsideConfig -count=1`
Expected: PASS.

**Step 3: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T07-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T07 --value COMPLETED --json
git add internal/modelroute/route.go internal/validation/model_hardcode_guard_test.go internal/workflows/persona_test.go internal/hooks/*.go internal/policy/*.go
git commit -m "test(policy): enforce no hardcoded runtime models outside config"
```

---

### Task T08: Build/Release Integration For Runtime-Only `.claude-plugin`

**Why:** Build outputs must consistently produce runtime artifacts and keep `.claude-plugin` read-only in steady state.

**What:** Update scripts and devx actions to generate both binaries and policy.

**How:** wire `build-runtime` into `scripts/build.sh` and add checks.

**Key Design Considerations:**
- build order: policy first, binaries second (or fail fast at policy stage).
- artifact paths fixed under `.claude-plugin`.

**Files:**
- Modify: `scripts/build.sh`
- Modify: `scripts/mp-devx`
- Modify: `cmd/mp-devx/main.go`
- Create: `internal/validation/runtime_artifact_test.go`

#### Subtask T08-S01: Build binaries + policy in one pipeline

**Why:** One command should prepare full runtime package.

**What:** ensure `scripts/build.sh` creates:
- `.claude-plugin/runtime/policy.json`
- `.claude-plugin/bin/mp`
- `.claude-plugin/bin/mp-devx`

**How:** call `./scripts/mp-devx -action build-runtime`.

**Key Design Considerations:** command should be idempotent.

**Step 1: Write failing artifact test**

Run: `go test ./internal/validation -run TestRuntimeBuildArtifactsExist -count=1`
Expected: FAIL.

**Step 2: Update build scripts and actions**

**Step 3: Run build + test**

Run:

```bash
./scripts/build.sh
go test ./internal/validation -run TestRuntimeBuildArtifactsExist -count=1
```

Expected: PASS.

**Step 4: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T08-S01 --value COMPLETED --json
```

#### Subtask T08-S02: Runtime read-only contract checks

**Why:** Prevent direct runtime edits bypassing source config.

**What:** Add validation that runtime policy is generated artifact, not hand-edited source of truth.

**How:** add checksum or generation marker validation.

**Key Design Considerations:** clear remediation: edit `config/*.yaml` then rebuild.

**Step 1: Add validation test for generated marker/checksum**

**Step 2: Implement marker/checksum in compiler output**

**Step 3: Run tests**

Run: `go test ./internal/validation ./internal/policy -count=1`
Expected: PASS.

**Step 4: Update status + commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T08-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T08 --value COMPLETED --json
git add scripts/build.sh scripts/mp-devx cmd/mp-devx/main.go internal/validation/runtime_artifact_test.go internal/policy/*.go
git commit -m "build(runtime): generate policy and binaries for runtime-only plugin bundle"
```

---

### Task T09: End-To-End Verification And Migration Docs

**Why:** Final safety net for behavior, operability, and handoff quality.

**What:** Add E2E tests and update docs for new config-driven flow.

**How:** verify workflow task routing, agent routing, hard/hint behavior, and fallback visibility toggle.

**Key Design Considerations:** tests must prove user-confirmed semantics, not just compile success.

**Files:**
- Create: `internal/policy/e2e_resolve_dispatch_test.go`
- Modify: `docs/ARCHITECTURE.md`
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `custom/docs/tool-project/tech-stack.md`
- Create: `docs/plans/evidence/model-routing/2026-03-04-config-driven-routing-verification.md`

#### Subtask T09-S01: E2E tests for resolver/enforcement/fallback

**Why:** Core runtime promises must be executable assertions.

**What:** add tests for:
- workflow task-specific model selection
- external hard enforcement with model arg
- one-hop automatic fallback to cross-provider or claude_code
- `/mp:config` show/hide routing details

**How:** table-driven integration tests with mocked executor runner.

**Key Design Considerations:** include both success and terminal failure cases.

**Step 1: Write failing E2E tests**

Run: `go test ./internal/policy -run TestE2E_ConfigDrivenRouting -count=1`
Expected: FAIL.

**Step 2: Implement missing wiring and mocks**

**Step 3: Run E2E tests**

Run: `go test ./internal/policy -run TestE2E_ConfigDrivenRouting -count=1`
Expected: PASS.

**Step 4: Update status**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T09-S01 --value COMPLETED --json
```

#### Subtask T09-S02: Operator docs + evidence capture

**Why:** Deployment and maintenance teams need explicit migration instructions.

**What:** document new source-of-truth and runtime artifact flow.

**How:** update key docs and capture real command transcript.

**Key Design Considerations:** doc must clearly distinguish dev-time vs run-time responsibilities.

**Step 1: Update docs**

- `docs/ARCHITECTURE.md`
- `docs/CLI-REFERENCE.md`
- `custom/docs/tool-project/tech-stack.md`

**Step 2: Run full verification suite**

Run:

```bash
go test ./... -count=1
./scripts/mp-devx -action build-policy
./scripts/build.sh
```

Expected: all pass, runtime artifacts present.

**Step 3: Record evidence transcript**

Create: `docs/plans/evidence/model-routing/2026-03-04-config-driven-routing-verification.md`

Include exact commands + key output snippets.

**Step 4: Update status + final commit**

```bash
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T09-S02 --value COMPLETED --json
go run ./cmd/mp state set --dir . --key plan.config_model_routing.T09 --value COMPLETED --json
git add internal/policy/e2e_resolve_dispatch_test.go docs/ARCHITECTURE.md docs/CLI-REFERENCE.md custom/docs/tool-project/tech-stack.md docs/plans/evidence/model-routing/2026-03-04-config-driven-routing-verification.md
git commit -m "docs(test): verify and document config-driven model routing runtime"
```

---

## Final Verification Gate (Must Pass Before Merge)

Run in order:

```bash
go test ./internal/policy ./internal/hooks ./internal/cli ./internal/workflows ./internal/validation ./cmd/mp-devx -count=1
go test ./... -count=1
./scripts/mp-devx -action build-policy
./scripts/build.sh
```

Expected:
- All tests pass
- `.claude-plugin/runtime/policy.json` exists and is fresh
- `.claude-plugin/bin/mp` and `.claude-plugin/bin/mp-devx` exist
- hooks/CLI metadata show model/fallback when `show_model_routing=true`
- fallback behavior is one-hop only

---

## Open Questions To Resolve During Execution

1. Whether `internal/modelroute` is fully deleted or retained as deprecated wrapper for one release.
2. Whether `/mp:config` should be implemented as `mp config` subcommand only or also explicit markdown command alias.
3. Whether generated policy checksum is embedded in JSON body or sidecar file.

