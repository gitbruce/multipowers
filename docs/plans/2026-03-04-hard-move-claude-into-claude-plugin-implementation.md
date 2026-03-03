# Hard Move .claude Into .claude-plugin Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Move the entire workspace `.claude` tree under `.claude-plugin/.claude` so plugin installation from local path or GitHub URL always resolves commands/skills from the plugin package root, and the repository root no longer exposes local `/command` entries. Root `.claude` must be deleted and stay deleted on `go`.

**Architecture:** Keep `.claude-plugin` as the only publishable plugin source root (`marketplace.json -> source: "./.claude-plugin"`). Move runtime markdown assets (commands, skills, references, state) into `.claude-plugin/.claude/*` and keep manifest references as `./.claude/...` relative to plugin source root. Update validation/tests/scripts to enforce the new physical layout and prevent regressions. For future `main -> go` merges (where `main` layout remains unchanged), use a deterministic import workflow that refreshes `.claude-plugin/.claude` from `main:.claude` and then removes root `.claude` in the same commit.

**Tech Stack:** Go (`internal/validation`, `internal/devx`), JSON manifests (`.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json`, `config/sync/claude-structure-rules.json`), Bash scripts (`scripts/*.sh`), Claude Code plugin manager.

---

### Task 1: Add Failing Layout Guard Test (TDD Red)

**Files:**
- Create: `internal/validation/claude_layout_test.go`
- Test: `internal/validation/claude_layout_test.go`

**Step 1: Write the failing test**

```go
package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestClaudeAssetsArePackagedUnderPluginRoot(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve caller")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))

	if _, err := os.Stat(filepath.Join(root, ".claude")); err == nil {
		t.Fatalf("root .claude must not exist after hard-move")
	}

	mustExist := []string{
		filepath.Join(root, ".claude-plugin", ".claude", "commands", "mp.md"),
		filepath.Join(root, ".claude-plugin", ".claude", "skills", "skill-prd.md"),
		filepath.Join(root, ".claude-plugin", ".claude", "references", "validation-gates.md"),
		filepath.Join(root, ".claude-plugin", ".claude", "state", "state-manager.md"),
	}
	for _, p := range mustExist {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing moved asset %s: %v", p, err)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run TestClaudeAssetsArePackagedUnderPluginRoot -count=1`
Expected: FAIL because `.claude` still exists at repository root before move.

**Step 3: Commit failing test scaffold**

```bash
git add internal/validation/claude_layout_test.go
git commit -m "test(validation): add failing guard for .claude hard-move layout"
```

### Task 2: Hard-Move `.claude` Directory Into `.claude-plugin/.claude`

**Files:**
- Move: `.claude` -> `.claude-plugin/.claude`
- Move: `.claude-plugin/commands/persona.md` -> `.claude-plugin/.claude/commands/persona.md`
- Delete (if empty): `.claude-plugin/commands`

**Step 1: Execute filesystem move**

```bash
mv .claude .claude-plugin/.claude
mkdir -p .claude-plugin/.claude/commands
if [ -f .claude-plugin/commands/persona.md ]; then
  mv .claude-plugin/commands/persona.md .claude-plugin/.claude/commands/persona.md
fi
rmdir .claude-plugin/commands 2>/dev/null || true
```

**Step 2: Verify moved structure**

Run: `find .claude-plugin/.claude -maxdepth 2 -type d | sort`
Expected: contains `.claude-plugin/.claude/commands`, `skills`, `references`, `state`.

**Step 3: Re-run the red test**

Run: `go test ./internal/validation -run TestClaudeAssetsArePackagedUnderPluginRoot -count=1`
Expected: PASS.

**Step 4: Commit move**

```bash
git add -A .claude .claude-plugin
git commit -m "refactor(layout): hard-move .claude under .claude-plugin/.claude"
```

### Task 3: Fix Plugin Manifest Registration for Persona and Packaged Paths

**Files:**
- Modify: `.claude-plugin/plugin.json`
- Modify: `internal/validation/persona_namespace_test.go`
- Test: `internal/validation/persona_namespace_test.go`

**Step 1: Update failing test expectations first**

```go
pluginPersona := filepath.Join(root, ".claude-plugin", ".claude", "commands", "persona.md")
if _, err := os.Stat(pluginPersona); err != nil {
	t.Fatalf("plugin persona command must exist: %v", err)
}

if !strings.Contains(content, "./.claude/commands/persona.md") {
	t.Fatalf("plugin.json must register persona from .claude/commands/persona.md")
}
if strings.Contains(content, "./.claude-plugin/commands/persona.md") {
	t.Fatalf("plugin.json must not register persona from legacy .claude-plugin/commands path")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run TestPersonaCommandIsPluginNamespacedOnly -count=1`
Expected: FAIL until manifest and persona path are aligned.

**Step 3: Write minimal manifest fix**

Update `.claude-plugin/plugin.json` persona entry to:

```json
"./.claude/commands/persona.md"
```

(no `./.claude-plugin/commands/persona.md` remains)

**Step 4: Run test to verify it passes**

Run: `go test ./internal/validation -run TestPersonaCommandIsPluginNamespacedOnly -count=1`
Expected: PASS.

**Step 5: Commit**

```bash
git add .claude-plugin/plugin.json internal/validation/persona_namespace_test.go
git commit -m "fix(plugin): register persona from packaged .claude command path"
```

### Task 4: Update Runtime Path Validation Tests for New Location

**Files:**
- Modify: `internal/validation/command_mp_path_test.go`
- Modify: `internal/validation/no_shell_runtime.go`
- Modify: `internal/validation/no_shell_runtime_test.go`
- Test: `internal/validation/*.go`

**Step 1: Write failing test updates**

Change command path root in `command_mp_path_test.go`:

```go
path := filepath.Join(root, ".claude-plugin", ".claude", "commands", name)
```

Change no-shell test fixture references:

```go
refs := []string{".claude-plugin/.claude/commands/persona.md:bash ./scripts/build.sh"}
refs := []string{".claude-plugin/.claude/commands/persona.md:${CLAUDE_PLUGIN_ROOT}/bin/mp persona --json"}
```

**Step 2: Run tests to verify red state**

Run: `go test ./internal/validation -count=1`
Expected: FAIL before scanner candidates are updated.

**Step 3: Write minimal implementation**

Update scanner candidates in `no_shell_runtime.go`:

```go
candidates := []string{
	".claude-plugin/.claude/commands",
	".claude-plugin/.claude/skills",
	".claude-plugin",
	"custom/docs/tool-project",
	".github/workflows",
	"Makefile",
	"docs/COMMAND-REFERENCE.md",
	"docs/CLI-REFERENCE.md",
}
```

**Step 4: Run tests to verify green state**

Run: `go test ./internal/validation -count=1`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/validation/command_mp_path_test.go internal/validation/no_shell_runtime.go internal/validation/no_shell_runtime_test.go
git commit -m "test(validation): align runtime path checks with .claude-plugin/.claude layout"
```

### Task 5: Update Structure Parity Rules and DevX Assertions

**Files:**
- Modify: `config/sync/claude-structure-rules.json`
- Modify: `internal/devx/structure_rules_test.go`
- Test: `internal/devx/*.go`

**Step 1: Update failing assertion first**

Replace root-target assertion with packaged-target assertion in `structure_rules_test.go`:

```go
if !strings.Contains(rule.TargetRoot, ".claude-plugin/.claude/") {
	t.Fatalf("target root must use packaged .claude path: %s", rule.TargetRoot)
}
```

**Step 2: Run test to confirm it fails before config change**

Run: `go test ./internal/devx -run TestLoadStructureRules_RootTargetsUseClaudeRoot -count=1`
Expected: FAIL.

**Step 3: Update config targets**

In `config/sync/claude-structure-rules.json`, keep `source_root` as upstream `.claude/*`, but move all `target_root` to:
- `.claude-plugin/.claude/commands`
- `.claude-plugin/.claude/skills`
- `.claude-plugin/.claude/references`
- `.claude-plugin/.claude/state`
- forked command targets under `.claude-plugin/.claude/commands/*.md`

**Step 4: Run devx tests**

Run: `go test ./internal/devx -count=1`
Expected: PASS.

**Step 5: Commit**

```bash
git add config/sync/claude-structure-rules.json internal/devx/structure_rules_test.go
git commit -m "chore(sync): point parity targets to packaged .claude paths"
```

### Task 6: Fix Utility Scripts That Still Assume Root `.claude`

**Files:**
- Modify: `scripts/build-openclaw.sh`
- Modify: `scripts/fix-command-frontmatter.sh`

**Step 1: Write failing check command**

Run: `rg -n "PLUGIN_ROOT/.claude|PROJECT_ROOT/.claude" scripts/build-openclaw.sh scripts/fix-command-frontmatter.sh`
Expected: at least one hit (legacy path assumption).

**Step 2: Write minimal implementation**

Update script roots:

```bash
SKILLS_DIR="$PLUGIN_ROOT/.claude-plugin/.claude/skills"
COMMANDS_DIR="$PLUGIN_ROOT/.claude-plugin/.claude/commands"
```

```bash
COMMANDS_DIR="$PROJECT_ROOT/.claude-plugin/.claude/commands"
```

**Step 3: Validate scripts syntax**

Run: `bash -n scripts/build-openclaw.sh scripts/fix-command-frontmatter.sh`
Expected: no output, exit code 0.

**Step 4: Commit**

```bash
git add scripts/build-openclaw.sh scripts/fix-command-frontmatter.sh
git commit -m "fix(scripts): read commands/skills from packaged .claude location"
```

### Task 7: End-to-End Regression Verification

**Files:**
- Test/verify only (no file creation required)

**Step 1: Run repository checks**

Run:

```bash
go test ./internal/validation -count=1
go test ./internal/devx -count=1
./scripts/validate-claude-structure.sh -source-ref main -target-ref HEAD -dry-run
```

Expected: all pass.

**Step 2: Ensure no root `.claude` remains**

Run: `test ! -d .claude && echo "OK: root .claude removed"`
Expected: prints `OK: root .claude removed`.

**Step 3: Ensure packaged tree is complete**

Run:

```bash
ls -1 .claude-plugin/.claude/commands | wc -l
ls -1 .claude-plugin/.claude/skills | wc -l
```

Expected: command/skill counts match prior root layout (excluding intentionally removed files like `skill-persona.md`).

**Step 4: Claude Code manual verification (outside repo tests)**

Run in Claude Code:

```text
/plugin uninstall mp@multipowers-plugins
/plugin marketplace add https://github.com/gitbruce/multipowers
/plugin install mp@multipowers-plugins --scope user
# restart Claude Code
/mp:persona list
```

Expected:
- `/mp:persona` works
- no fallback to `/persona`
- installed plugin cache contains packaged `.claude` commands/skills under `mp/<version>/...`

**Step 5: Final commit**

```bash
git add -A
git commit -m "refactor(plugin-layout): hard-move .claude under .claude-plugin with path parity"
```

### Task 8: Add Repeatable Main->Go Import Script (Merge-Friendly with Root `.claude` Deleted)

**Files:**
- Create: `scripts/sync-main-claude-into-plugin.sh`
- Modify: `README.md` (or `docs/CLI-REFERENCE.md`) with merge playbook snippet

**Step 1: Write the script**

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Pull latest upstream/main asset tree into workspace path.
git checkout main -- .claude

# Replace packaged runtime assets with upstream snapshot.
rm -rf .claude-plugin/.claude
mv .claude .claude-plugin/.claude

# Safety check: root .claude must not remain.
test ! -d .claude

echo "Synced main:.claude -> .claude-plugin/.claude"
```

**Step 2: Validate script behavior**

Run:

```bash
bash -n scripts/sync-main-claude-into-plugin.sh
./scripts/sync-main-claude-into-plugin.sh
test ! -d .claude && echo "OK: root .claude removed after sync"
```

Expected: script succeeds and root `.claude` is absent.

**Step 3: Document merge playbook**

Add a short operation sequence:

```bash
git fetch upstream --prune
git merge upstream/main -X theirs
./scripts/sync-main-claude-into-plugin.sh
go test ./internal/validation -count=1
go test ./internal/devx -count=1
git add -A
git commit -m "chore(sync): refresh packaged .claude assets from main"
```

**Step 4: Commit**

```bash
git add scripts/sync-main-claude-into-plugin.sh README.md docs/CLI-REFERENCE.md
git commit -m "chore(sync): add deterministic main-to-packaged-claude import script"
```

### Task 9: Post-Merge Safety Net (Optional but Recommended)

**Files:**
- Modify: `internal/validation/persona_namespace_test.go` (or create dedicated layout guard test)

**Step 1: Add regression assertion for legacy path ban**

```go
if strings.Contains(content, "./.claude-plugin/commands/") {
	t.Fatalf("plugin.json must not use legacy .claude-plugin/commands path")
}
```

**Step 2: Verify test**

Run: `go test ./internal/validation -count=1`
Expected: PASS.

**Step 3: Commit**

```bash
git add internal/validation/persona_namespace_test.go
git commit -m "test(validation): block legacy command path regressions"
```
