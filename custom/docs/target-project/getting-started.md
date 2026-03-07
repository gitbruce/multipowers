# Getting Started (Target Project Users)

## Install Plugin (User Scope)

```text
/plugin marketplace add /mnt/f/src/ai/multipowers
/plugin install mp@multipowers-plugins --scope user
```

## Uninstall Plugin / Marketplace (User Scope)

```text
rm -rf ~/.claude/plugins/cache/multipowers-plugins
/plugin uninstall mp@multipowers-plugins --scope user
/plugin marketplace remove multipowers-plugins
```

## Initialize in Your Project

In your target project directory:

```text
/mp:init
```

Expected:
- creates `.multipowers/` in your target project
- initializes project context files and tracks registry
- writes the canonical registry to `.multipowers/tracks/tracks.md`
- creates spec-driven artifacts under `.multipowers/tracks/<track_id>/`
- creates `.multipowers/CLAUDE.md` and `.multipowers/FAQ.md`
- creates `.multipowers/context/runtime.json` for runtime/pre-run settings (`pre_run.enabled=false` by default)

Legacy note:
- `.multipowers/tracks.md` is not compatible with the current runtime and is not read.

## Optional: Configure Pre-Run Hooks

During `/mp:init`, you can configure pre-run hooks (for any runtime/toolchain).

Example use cases:
- activate an environment
- source a shell profile
- export required env vars

Hooks are stored in:
- `.multipowers/context/runtime.json`

All external provider executions (`codex`, `gemini`, `claude`) read this file before running.

## Run Spec-Driven Commands

- `/mp:plan`
- `/mp:discover`, `/mp:define`, `/mp:develop`, `/mp:deliver`
- `/mp:embrace`, `/mp:review`, `/mp:debate`, `/mp:research`

If context is missing or incomplete, the command blocks with `run_init` guidance; `/mp:init` is not auto-triggered silently.

## Manage Track Groups Explicitly

Once a spec track exists, advance implementation groups explicitly:

```text
mp track group-start --track-id <track_id> --group g1 --execution-mode workspace --json
mp track group-complete --track-id <track_id> --group g1 --commit-sha <sha> --json
```

Notes:
- all spec artifacts stay under `.multipowers/tracks/<track_id>/`
- `.multipowers/tracks/tracks.md` is the only runtime registry path
- an active group blocks the next spec pipeline step until commit and verification evidence are recorded

## FAQ Learning Loop

- `.multipowers/FAQ.md` is auto-generated and auto-refined from observed failures.
- Entries are categorized by error type and deduplicated.
- No manual maintenance, backup, or archive files are required.
