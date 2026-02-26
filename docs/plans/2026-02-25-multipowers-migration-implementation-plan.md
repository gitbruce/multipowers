# Multipowers Overlay Migration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Migrate from the current `multipowers` customization model to a low-conflict overlay architecture on a new `multipowers` design branch, with explicit handling of the current in-progress rebase (`resolve` or `abort`) and operator-first documentation.

**Architecture:** Keep upstream core files mostly unchanged and move customization policy into `custom/` (config + overlay libraries + command source), with thin hooks in `bin/mp` and a deterministic overlay-apply script. Build a merge-based sync workflow around `main -> multipowers-design` with idempotent overlay reapplication and contract tests.

**Tech Stack:** Bash (`scripts/*.sh`), Markdown docs (`docs/multipowers/*`), shell test harness (`tests/helpers/test-framework.sh`), git branching/merge workflow.

---

## Execution Tracker (Update Checkboxes During Execution)

Execution note (2026-02-25): direct full merge of `multipowers` into new `multipowers-design` produced an excessive conflict set. The migration completed using the same target architecture but via selective overlay reimplementation plus validation, then merge-based sync workflow verification.

### Task 1: Stabilize Repository State Before Migration
- [x] Task 1 complete
- [x] T1.1 Capture current state and save patch snapshots
- [x] T1.2 Create safety backup branch from detached HEAD
- [x] T1.3 Attempt rebase resolve path (`git rebase --continue`)
- [x] T1.4 Fallback to rebase abort path if needed (`git rebase --abort`)
- [x] T1.5 Verify clean/non-conflicted status
- [x] T1.6 Create checkpoint commit

### Task 2: Create New Migration Branch and Baseline
- [x] Task 2 complete
- [x] T2.1 Fast-forward `main` to `upstream/main`
- [x] T2.2 Create `multipowers-design` branch
- [x] T2.3 Merge existing `multipowers` baseline
- [x] T2.4 Verify ancestry/log
- [x] T2.5 Create baseline checkpoint commit

### Task 3: Scaffold Overlay Directory and Config Contracts
- [x] Task 3 complete
- [x] T3.1 Add failing config contract test
- [x] T3.2 Run test and confirm fail
- [x] T3.3 Create `custom/config/*.json` and `custom/README.md`
- [x] T3.4 Re-run test and confirm pass
- [x] T3.5 Commit

### Task 4: Add Overlay Loader and Routing Libraries
- [x] Task 4 complete
- [x] T4.1 Add failing overlay loader/routing test
- [x] T4.2 Run test and confirm fail
- [x] T4.3 Create `custom/lib/*.sh`
- [x] T4.4 Re-run test and confirm pass
- [x] T4.5 Commit

### Task 5: Integrate Thin Hooks into Core Orchestrator
- [x] Task 5 complete
- [x] T5.1 Add failing orchestrator hook test
- [x] T5.2 Run test and confirm fail
- [x] T5.3 Add optional hook source/invocation in `bin/mp`
- [x] T5.4 Re-run test and confirm pass
- [x] T5.5 Commit

### Task 6: Migrate `/mp:persona` to Overlay Source-of-Truth
- [x] Task 6 complete
- [x] T6.1 Add failing command sync test
- [x] T6.2 Run test and confirm fail
- [x] T6.3 Create `custom/commands/persona.md` and `scripts/mp-devx overlay`
- [x] T6.4 Run command/frontmatter registration tests and confirm pass
- [x] T6.5 Commit

### Task 7: Restore/Add Merge-Based Sync Script
- [x] Task 7 complete
- [x] T7.1 Add failing sync integration test
- [x] T7.2 Run test and confirm fail
- [x] T7.3 Create `scripts/mp-devx sync`
- [x] T7.4 Re-run test and confirm pass
- [x] T7.5 Commit

### Task 8: Add Operator-First Multipowers Docs Hub
- [x] Task 8 complete
- [x] T8.1 Add failing docs navigation test
- [x] T8.2 Run test and confirm fail
- [x] T8.3 Create docs hub + feature pages + root/docs index pointers
- [x] T8.4 Run docs tests and confirm pass
- [x] T8.5 Commit

### Task 9: Add Compatibility and Rebase Recovery Playbooks
- [x] Task 9 complete
- [x] T9.1 Add compatibility matrix content
- [x] T9.2 Add resolve/abort runbooks in conflict guide
- [x] T9.3 Add manual verification checklist
- [x] T9.4 Execute checklist and record outcomes
- [x] T9.5 Commit

### Task 10: End-to-End Validation and Cutover
- [x] Task 10 complete
- [x] T10.1 Run targeted unit/integration tests
- [x] T10.2 Run smoke sync on `multipowers-design`
- [x] T10.3 Validate persona + model/proxy behavior
- [x] T10.4 Verify branch is PR-ready
- [x] T10.5 Finalize commit if needed

### Task 11: Post-Cutover Operating Procedure
- [x] Task 11 complete
- [x] T11.1 Document routine sync command sequence
- [x] T11.2 Document conflict SLA and fallback path
- [x] T11.3 Execute command sequence once
- [x] T11.4 Add example sync transcript
- [x] T11.5 Commit docs update

### Final Verification Gate
- [x] F1 Verify clean git status
- [x] F2 Run `tests/integration/test-sync-overlay.sh`
- [x] F3 Run `tests/test-command-registration.sh`
- [x] F4 Run `tests/test-model-config-simple.sh`

### Task 1: Stabilize Repository State Before Migration (Resolve or Abort Rebase)

**Files:**
- Modify: none
- Create: none
- Test: repo state via `git status`

**Step 1: Capture current state and safety snapshot**

Run:
```bash
git status --short --branch
git rev-parse --short HEAD
git diff > /tmp/multipowers-rebase-working.patch
git diff --cached > /tmp/multipowers-rebase-staged.patch
```
Expected: patch files created and current detached `HEAD` shown.

**Step 2: Create safety branch from current detached HEAD**

Run:
```bash
git switch -c backup/rebase-2026-02-25
```
Expected: new local backup branch created.

**Step 3: Decision gate - try resolving current rebase first (preferred when conflicts are nearly done)**

Run:
```bash
git rebase --continue
```
Expected: either next conflict appears or rebase completes.

If conflicts remain, resolve only these files, then continue:
```bash
git add .claude-plugin/plugin.json README.md docs/COMMAND-REFERENCE.md bin/mp
git rebase --continue
```
Expected: rebase completes and branch name is `multipowers` (not detached).

**Step 4: Decision gate fallback - abort rebase if resolution is too costly (>30 min)**

Run:
```bash
git rebase --abort
```
Expected: rebase state cleared; branch returns to pre-rebase pointer.

**Step 5: Verify clean branch state for migration**

Run:
```bash
git status --short --branch
```
Expected: on named branch (`multipowers` or backup branch), no `UU` entries.

**Step 6: Commit checkpoint**

Run:
```bash
git add -A
git commit -m "chore: checkpoint before multipowers overlay migration" || true
```
Expected: commit created if staged changes exist; otherwise no-op.

### Task 2: Create New Migration Branch and Baseline

**Files:**
- Modify: none
- Create: none
- Test: branch ancestry verification

**Step 1: Ensure `main` tracks upstream baseline**

Run:
```bash
git fetch upstream origin
git switch main
git merge --ff-only upstream/main
```
Expected: `main` fast-forwards to upstream tip.

**Step 2: Create new design branch**

Run:
```bash
git switch -c multipowers-design
```
Expected: branch created from updated `main`.

**Step 3: Bring over current multipowers behavior for migration reference**

Run:
```bash
git merge --no-ff multipowers -m "chore: import existing multipowers baseline for overlay migration"
```
Expected: existing custom behavior present on `multipowers-design`.

**Step 4: Verify baseline commit graph**

Run:
```bash
git log --oneline --decorate -n 12
```
Expected: merge commit linking `multipowers` into `multipowers-design`.

**Step 5: Commit checkpoint**

Run:
```bash
git commit --allow-empty -m "chore: establish multipowers-design migration baseline"
```
Expected: explicit checkpoint commit.

### Task 3: Scaffold Overlay Directory and Config Contracts (TDD)

**Files:**
- Create: `custom/README.md`
- Create: `custom/config/models.json`
- Create: `custom/config/proxy.json`
- Create: `custom/config/persona-lanes.json`
- Create: `tests/unit/test-custom-config-contracts.sh`
- Test: `tests/unit/test-custom-config-contracts.sh`

**Step 1: Write failing config contract test**

```bash
cat > tests/unit/test-custom-config-contracts.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$PROJECT_ROOT/tests/helpers/test-framework.sh"

test_case "models.json exists with required keys"
if python3 -e '.providers and .role_routing and .fallback_lane' "$PROJECT_ROOT/custom/config/models.json" >/dev/null 2>&1; then
  test_pass "models.json contract valid"
else
  test_fail "models.json missing required keys"
fi

test_case "proxy.json exists with required keys"
if python3 -e '.enabled != null and .providers and .host and .port' "$PROJECT_ROOT/custom/config/proxy.json" >/dev/null 2>&1; then
  test_pass "proxy.json contract valid"
else
  test_fail "proxy.json missing required keys"
fi

test_case "persona-lanes.json exists with required keys"
if python3 -e '.personas and .fallback_lane' "$PROJECT_ROOT/custom/config/persona-lanes.json" >/dev/null 2>&1; then
  test_pass "persona-lanes contract valid"
else
  test_fail "persona-lanes missing required keys"
fi
SH
chmod +x tests/unit/test-custom-config-contracts.sh
```

**Step 2: Run test to verify failure**

Run:
```bash
bash tests/unit/test-custom-config-contracts.sh
```
Expected: FAIL because `custom/config/*.json` does not yet exist.

**Step 3: Add minimal config files and overlay README**

```bash
mkdir -p custom/config
cat > custom/config/models.json <<'JSON'
{
  "providers": {
    "codex": "gpt-5.3-codex",
    "gemini": "gemini-3-pro-preview",
    "claude_heavy": "claude-opus",
    "claude_light": "claude-sonnet"
  },
  "role_routing": {
    "heavy_coding": "claude_heavy",
    "docs_and_tests": "claude_light",
    "architecture_review_decision": "codex",
    "external_search_business": "gemini"
  },
  "fallback_lane": "claude_light"
}
JSON

cat > custom/config/proxy.json <<'JSON'
{
  "enabled": true,
  "providers": ["codex", "gemini"],
  "host": "127.0.0.1",
  "port": 7890,
  "no_proxy": ["localhost", "127.0.0.1"]
}
JSON

cat > custom/config/persona-lanes.json <<'JSON'
{
  "personas": {
    "backend-architect": "codex",
    "code-reviewer": "codex",
    "business-analyst": "gemini",
    "docs-architect": "claude_light"
  },
  "fallback_lane": "claude_light"
}
JSON

cat > custom/README.md <<'MD'
# Multipowers Overlay

This directory contains fork-specific customizations designed to minimize conflicts with upstream core files.

- `config/`: declarative policy (models, proxy, persona lane mapping)
- runtime hooks live in `custom/lib/`
- docs live in `docs/multipowers/`
MD
```

**Step 4: Run test to verify pass**

Run:
```bash
bash tests/unit/test-custom-config-contracts.sh
```
Expected: PASS for all three config contracts.

**Step 5: Commit**

```bash
git add custom/README.md custom/config/*.json tests/unit/test-custom-config-contracts.sh
git commit -m "feat(custom): add overlay config contracts and baseline files"
```

### Task 4: Add Overlay Loader and Routing Libraries (TDD)

**Files:**
- Create: `custom/lib/overlay-loader.sh`
- Create: `custom/lib/model-routing.sh`
- Create: `custom/lib/proxy-routing.sh`
- Create: `tests/unit/test-custom-overlay-loader.sh`
- Test: `tests/unit/test-custom-overlay-loader.sh`

**Step 1: Write failing loader/routing unit test**

```bash
cat > tests/unit/test-custom-overlay-loader.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$PROJECT_ROOT/tests/helpers/test-framework.sh"

source "$PROJECT_ROOT/custom/lib/overlay-loader.sh"
source "$PROJECT_ROOT/custom/lib/model-routing.sh"
source "$PROJECT_ROOT/custom/lib/proxy-routing.sh"

test_case "resolve_custom_model_for_role returns codex for architecture review"
model="$(resolve_custom_model_for_role architecture_review_decision)"
if [[ "$model" == "gpt-5.3-codex" ]]; then
  test_pass "role routing works"
else
  test_fail "expected gpt-5.3-codex, got: $model"
fi

test_case "proxy export string includes configured host/port"
proxy_url="$(custom_proxy_url_for_provider codex)"
if [[ "$proxy_url" == "http://127.0.0.1:7890" ]]; then
  test_pass "proxy routing works"
else
  test_fail "unexpected proxy URL: $proxy_url"
fi
SH
chmod +x tests/unit/test-custom-overlay-loader.sh
```

**Step 2: Run test to verify failure**

Run:
```bash
bash tests/unit/test-custom-overlay-loader.sh
```
Expected: FAIL because overlay libs do not exist.

**Step 3: Write minimal overlay libs**

```bash
mkdir -p custom/lib
cat > custom/lib/overlay-loader.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail
CUSTOM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CUSTOM_CONFIG_DIR="$CUSTOM_ROOT/config"

custom_require_config() {
  local file="$1"
  [[ -f "$CUSTOM_CONFIG_DIR/$file" ]]
}
SH

cat > custom/lib/model-routing.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail

resolve_custom_model_for_role() {
  local role="$1"
  local config_file="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/config/models.json"
  local lane
  lane="$(python3 -r --arg role "$role" '.role_routing[$role] // .fallback_lane' "$config_file")"
  python3 -r --arg lane "$lane" '.providers[$lane]' "$config_file"
}
SH

cat > custom/lib/proxy-routing.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail

custom_proxy_url_for_provider() {
  local provider="$1"
  local config_file="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/config/proxy.json"
  local enabled host port
  enabled="$(python3 -r '.enabled' "$config_file")"
  host="$(python3 -r '.host' "$config_file")"
  port="$(python3 -r '.port' "$config_file")"
  if [[ "$enabled" != "true" ]]; then
    echo ""
    return 0
  fi
  if python3 -e --arg p "$provider" '.providers | index($p)' "$config_file" >/dev/null 2>&1; then
    echo "http://${host}:${port}"
  else
    echo ""
  fi
}
SH
chmod +x custom/lib/*.sh
```

**Step 4: Run test to verify pass**

Run:
```bash
bash tests/unit/test-custom-overlay-loader.sh
```
Expected: PASS.

**Step 5: Commit**

```bash
git add custom/lib/*.sh tests/unit/test-custom-overlay-loader.sh
git commit -m "feat(custom): add overlay loader and routing libraries"
```

### Task 5: Integrate Thin Hooks into Core Orchestrator (TDD)

**Files:**
- Modify: `bin/mp`
- Create: `tests/unit/test-custom-orchestrate-hooks.sh`
- Test: `tests/unit/test-custom-orchestrate-hooks.sh`

**Step 1: Write failing hook-presence test**

```bash
cat > tests/unit/test-custom-orchestrate-hooks.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$PROJECT_ROOT/tests/helpers/test-framework.sh"

test_case "orchestrate sources custom overlay loader"
if grep -q 'custom/lib/overlay-loader.sh' "$PROJECT_ROOT/bin/mp"; then
  test_pass "overlay loader source found"
else
  test_fail "overlay loader source missing"
fi

test_case "orchestrate can call custom proxy resolver"
if grep -q 'custom_proxy_url_for_provider' "$PROJECT_ROOT/bin/mp"; then
  test_pass "custom proxy hook found"
else
  test_fail "custom proxy hook missing"
fi
SH
chmod +x tests/unit/test-custom-orchestrate-hooks.sh
```

**Step 2: Run test to verify failure**

Run:
```bash
bash tests/unit/test-custom-orchestrate-hooks.sh
```
Expected: FAIL until hook lines are added.

**Step 3: Add minimal non-breaking hooks in `bin/mp`**

Add near existing `source` block:
```bash
# Source custom overlay hooks (multipowers, optional)
if [[ -f "${PLUGIN_DIR}/custom/lib/overlay-loader.sh" ]]; then
    source "${PLUGIN_DIR}/custom/lib/overlay-loader.sh" 2>/dev/null || true
fi
if [[ -f "${PLUGIN_DIR}/custom/lib/model-routing.sh" ]]; then
    source "${PLUGIN_DIR}/custom/lib/model-routing.sh" 2>/dev/null || true
fi
if [[ -f "${PLUGIN_DIR}/custom/lib/proxy-routing.sh" ]]; then
    source "${PLUGIN_DIR}/custom/lib/proxy-routing.sh" 2>/dev/null || true
fi
```

In proxy setup path, add conditional override call:
```bash
if declare -F custom_proxy_url_for_provider >/dev/null 2>&1; then
    custom_proxy_url="$(custom_proxy_url_for_provider "$provider")"
    if [[ -n "$custom_proxy_url" ]]; then
        export http_proxy="$custom_proxy_url" https_proxy="$custom_proxy_url"
        export HTTP_PROXY="$custom_proxy_url" HTTPS_PROXY="$custom_proxy_url"
    fi
fi
```

**Step 4: Run test to verify pass**

Run:
```bash
bash tests/unit/test-custom-orchestrate-hooks.sh
```
Expected: PASS.

**Step 5: Commit**

```bash
git add bin/mp tests/unit/test-custom-orchestrate-hooks.sh
git commit -m "feat(core): add thin optional hooks for multipowers overlay"
```

### Task 6: Migrate `/mp:persona` to Overlay Source-of-Truth (TDD)

**Files:**
- Create: `custom/commands/persona.md`
- Create: `scripts/mp-devx overlay`
- Modify: `.claude-plugin/.claude/commands/persona.md`
- Test: `tests/unit/test-command-frontmatter.sh`
- Test: `tests/test-command-registration.sh`

**Step 1: Write failing test for command source sync script**

```bash
cat > tests/unit/test-custom-command-sync.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$PROJECT_ROOT/tests/helpers/test-framework.sh"

test_case "apply-custom-overlay script exists and executable"
if [[ -x "$PROJECT_ROOT/scripts/mp-devx overlay" ]]; then
  test_pass "overlay apply script exists"
else
  test_fail "missing scripts/mp-devx overlay"
fi
SH
chmod +x tests/unit/test-custom-command-sync.sh
```

**Step 2: Run test to verify failure**

Run:
```bash
bash tests/unit/test-custom-command-sync.sh
```
Expected: FAIL before script exists.

**Step 3: Create command source and apply script**

```bash
mkdir -p custom/commands
cp .claude-plugin/.claude/commands/persona.md custom/commands/persona.md

cat > scripts/mp-devx overlay <<'SH'
#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

mkdir -p "$ROOT_DIR/.claude-plugin/.claude/commands"
cp "$ROOT_DIR/custom/commands/persona.md" "$ROOT_DIR/.claude-plugin/.claude/commands/persona.md"

python3 empty "$ROOT_DIR/custom/config/models.json" >/dev/null
python3 empty "$ROOT_DIR/custom/config/proxy.json" >/dev/null
python3 empty "$ROOT_DIR/custom/config/persona-lanes.json" >/dev/null

echo "Overlay applied successfully"
SH
chmod +x scripts/mp-devx overlay
```

**Step 4: Run tests and verify pass**

Run:
```bash
bash tests/unit/test-custom-command-sync.sh
bash tests/unit/test-command-frontmatter.sh
bash tests/test-command-registration.sh
```
Expected: PASS with persona command still registered.

**Step 5: Commit**

```bash
git add custom/commands/persona.md scripts/mp-devx overlay .claude-plugin/.claude/commands/persona.md tests/unit/test-custom-command-sync.sh
git commit -m "feat(custom): make persona command overlay-managed"
```

### Task 7: Restore/Add Sync Script for Merge-Based Upstream Updates (TDD)

**Files:**
- Create: `scripts/mp-devx sync`
- Create: `tests/integration/test-sync-overlay.sh`
- Test: `tests/integration/test-sync-overlay.sh`

**Step 1: Write failing integration test for sync script**

```bash
cat > tests/integration/test-sync-overlay.sh <<'SH'
#!/usr/bin/env bash
set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$PROJECT_ROOT/tests/helpers/test-framework.sh"

test_case "sync-upstream script exists"
if [[ -x "$PROJECT_ROOT/scripts/mp-devx sync" ]]; then
  test_pass "sync-upstream exists"
else
  test_fail "scripts/mp-devx sync missing"
fi

test_case "apply-custom-overlay is callable"
if "$PROJECT_ROOT/scripts/mp-devx overlay" >/dev/null 2>&1; then
  test_pass "overlay apply callable"
else
  test_fail "overlay apply failed"
fi
SH
chmod +x tests/integration/test-sync-overlay.sh
```

**Step 2: Run test to verify failure**

Run:
```bash
bash tests/integration/test-sync-overlay.sh
```
Expected: FAIL because sync script is currently missing.

**Step 3: Create `scripts/mp-devx sync`**

```bash
cat > scripts/mp-devx sync <<'SH'
#!/usr/bin/env bash
set -euo pipefail

UPSTREAM_REMOTE="${UPSTREAM_REMOTE:-upstream}"
MAIN_BRANCH="${MAIN_BRANCH:-main}"
TARGET_BRANCH="${TARGET_BRANCH:-multipowers-design}"

current_branch="$(git branch --show-current || true)"

git fetch "$UPSTREAM_REMOTE" origin

git switch "$MAIN_BRANCH"
git merge --ff-only "$UPSTREAM_REMOTE/$MAIN_BRANCH"

git switch "$TARGET_BRANCH"
git merge "$MAIN_BRANCH" -m "chore(sync): merge main into $TARGET_BRANCH"

"$(dirname "$0")/mp-devx overlay"

echo "Sync complete for $TARGET_BRANCH"

if [[ -n "$current_branch" && "$current_branch" != "$TARGET_BRANCH" ]]; then
  git switch "$current_branch"
fi
SH
chmod +x scripts/mp-devx sync
```

**Step 4: Run test to verify pass**

Run:
```bash
bash tests/integration/test-sync-overlay.sh
```
Expected: PASS.

**Step 5: Commit**

```bash
git add scripts/mp-devx sync tests/integration/test-sync-overlay.sh
git commit -m "feat(sync): add merge-based upstream sync script with overlay apply"
```

### Task 8: Add Operator-First Multipowers Docs Hub and Feature Pages (TDD)

**Files:**
- Create: `docs/INDEX.md`
- Create: `docs/multipowers/README.md`
- Create: `docs/multipowers/getting-started.md`
- Create: `docs/multipowers/daily-usage.md`
- Create: `docs/multipowers/customizations/models-and-lanes.md`
- Create: `docs/multipowers/customizations/proxy-routing.md`
- Create: `docs/multipowers/customizations/persona-command.md`
- Create: `docs/multipowers/sync/upstream-sync-playbook.md`
- Create: `docs/multipowers/sync/conflict-resolution.md`
- Create: `docs/multipowers/troubleshooting.md`
- Create: `docs/multipowers/reference/config-schema.md`
- Create: `docs/multipowers/reference/compatibility.md`
- Create: `docs/multipowers/reference/faq.md`
- Modify: `README.md`
- Test: `tests/integration/test-readme-compliance.sh`

**Step 1: Write failing docs navigation test**

Create `tests/unit/test-multipowers-docs.sh`:
```bash
#!/usr/bin/env bash
set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "$PROJECT_ROOT/tests/helpers/test-framework.sh"

required=(
  "$PROJECT_ROOT/docs/multipowers/README.md"
  "$PROJECT_ROOT/docs/multipowers/customizations/models-and-lanes.md"
  "$PROJECT_ROOT/docs/multipowers/customizations/proxy-routing.md"
  "$PROJECT_ROOT/docs/multipowers/customizations/persona-command.md"
)
for f in "${required[@]}"; do
  test_case "doc exists: $(basename "$f")"
  [[ -f "$f" ]] && test_pass "exists" || test_fail "missing $f"
done
```

**Step 2: Run test to verify failure**

Run:
```bash
chmod +x tests/unit/test-multipowers-docs.sh
bash tests/unit/test-multipowers-docs.sh
```
Expected: FAIL because docs not yet created.

**Step 3: Create docs with required section template**

For each feature page (`models-and-lanes.md`, `proxy-routing.md`, `persona-command.md`), include these headings:
```markdown
## What Changed From Upstream
## Why This Exists
## How To Use
## Operational Impact
## Rollback Path
```

Add root README pointer section:
```markdown
## Multipowers Customization

This fork maintains an operator-first customization layer. Start at `docs/multipowers/README.md`.
```

**Step 4: Run docs tests**

Run:
```bash
bash tests/unit/test-multipowers-docs.sh
bash tests/integration/test-readme-compliance.sh
```
Expected: PASS.

**Step 5: Commit**

```bash
git add docs/INDEX.md docs/multipowers README.md tests/unit/test-multipowers-docs.sh
git commit -m "docs(multipowers): add operator-first customization docs hub and feature pages"
```

### Task 9: Add Migration Compatibility and Rollback Playbooks

**Files:**
- Modify: `docs/multipowers/reference/compatibility.md`
- Modify: `docs/multipowers/sync/conflict-resolution.md`
- Test: manual command verification checklist

**Step 1: Add compatibility matrix content**

Include matrix columns:
- upstream version/tag
- multipowers overlay version
- required migration action
- known caveats

**Step 2: Add conflict resolution runbooks**

Include both flows:
- `Resolve rebase` runbook (with `git add` + `git rebase --continue` loop)
- `Abort rebase` runbook (`git rebase --abort`, `git switch`, re-run scripted sync)

**Step 3: Add manual verification checklist**

Commands:
```bash
git status --short --branch
./scripts/mp-devx sync
./scripts/mp-devx overlay
./bin/mp persona list
```
Expected: clean state and functional persona command lane output.

**Step 4: Run verification checklist**

Run commands above and record outcomes in docs where relevant.

**Step 5: Commit**

```bash
git add docs/multipowers/reference/compatibility.md docs/multipowers/sync/conflict-resolution.md
git commit -m "docs(sync): add rebase resolve/abort runbooks and compatibility matrix"
```

### Task 10: End-to-End Validation and Cutover to New Branch

**Files:**
- Modify: none (validation only)
- Test: targeted unit/integration suite and smoke commands

**Step 1: Run targeted test suite**

Run:
```bash
bash tests/unit/test-custom-config-contracts.sh
bash tests/unit/test-custom-overlay-loader.sh
bash tests/unit/test-custom-orchestrate-hooks.sh
bash tests/unit/test-custom-command-sync.sh
bash tests/unit/test-multipowers-docs.sh
bash tests/integration/test-sync-overlay.sh
bash tests/test-command-registration.sh
bash tests/test-model-config-simple.sh
```
Expected: all PASS.

**Step 2: Run smoke sync flow on new branch**

Run:
```bash
git switch multipowers-design
./scripts/mp-devx sync
```
Expected: successful merge from `main`, overlay reapplied, no runtime errors.

**Step 3: Validate persona + model/proxy behavior manually**

Run:
```bash
./bin/mp persona list
./bin/mp model-config
```
Expected: persona command works and configured model lanes reflect overlay policy.

**Step 4: Create cutover branch/PR checkpoint**

Run:
```bash
git log --oneline --decorate -n 20
git status --short
```
Expected: clean branch ready for PR from `multipowers-design`.

**Step 5: Commit (if any remaining generated artifacts)**

```bash
git add -A
git commit -m "chore(release): finalize multipowers overlay migration" || true
```

### Task 11: Post-Cutover Operating Procedure (Merge-Based Sync)

**Files:**
- Modify: `docs/multipowers/sync/upstream-sync-playbook.md`
- Test: dry-run operator checklist

**Step 1: Document routine sync sequence**

```bash
git fetch upstream origin
git switch main
git merge --ff-only upstream/main
git switch multipowers-design
git merge main -m "chore(sync): merge main into multipowers-design"
./scripts/mp-devx overlay
bash tests/integration/test-sync-overlay.sh
```

**Step 2: Document conflict SLA and fallback**

- if unresolved >30 min: stop manual conflicting, run abort runbook.
- restore from `backup/rebase-<date>`.

**Step 3: Validate playbook commands**

Run each command sequence once in local environment.

**Step 4: Record example sync transcript in docs**

Include command + expected output summary.

**Step 5: Commit**

```bash
git add docs/multipowers/sync/upstream-sync-playbook.md
git commit -m "docs(ops): finalize multipowers merge-based sync operating procedure"
```

## Final Verification Gate (Must Pass Before Declaring Complete)

Run:
```bash
git status --short --branch
bash tests/integration/test-sync-overlay.sh
bash tests/test-command-registration.sh
bash tests/test-model-config-simple.sh
```
Expected:
- clean branch
- all required tests passing
- sync script + overlay apply deterministic and repeatable

## Notes for Executor

- Keep all policy changes in `custom/*`; avoid expanding custom logic directly inside core files.
- Keep `bin/mp` hook footprint minimal.
- Do not skip the rebase decision gate; migration starts only after a stable branch state is achieved.
- Prefer small commits per task exactly as listed.
