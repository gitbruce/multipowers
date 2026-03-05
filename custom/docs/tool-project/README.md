# Tool Project Maintainer Guide

Audience: maintainers of this repository (`multipowers` fork), not end users.

## Naming Baseline

- Slash command namespace: `/mp:*`
- Plugin id: `mp`
- Marketplace id: `multipowers-plugins`

## Branch and Sync Discipline

1. `main` mirrors `upstream/main` only.
2. All custom work stays on `multipowers`.
3. Never merge `multipowers` back into `main`.
4. Regular sync path: `upstream/main -> main -> multipowers`.

Primary runbook:
- [upstream-sync-playbook.md](/mnt/f/src/ai/multipowers/custom/docs/sync/upstream-sync-playbook.md)
- [conflict-resolution.md](/mnt/f/src/ai/multipowers/custom/docs/sync/conflict-resolution.md)

## Customization Architecture

Keep customization isolated in `custom/*`:
- `custom/config/`: model, persona-lane, proxy, and Conductor setup protocol config

- `scripts/`: developer wrappers for local build/run (`mp`, `mp-devx`)
- `custom/templates/`: target-project templates (`conductor/*`, `CLAUDE.md`, `FAQ.md`)

## Operational Commands

```bash
./scripts/mp-devx sync
./scripts/mp-devx overlay
go test ./...
```

## Marketplace and Plugin (User Scope)

Install / refresh:

```text
/plugin marketplace add /mnt/f/src/ai/multipowers/.claude-plugin/marketplace.json
/plugin install mp@multipowers-plugins --scope user
```

Uninstall / remove:

```text
/plugin uninstall mp@multipowers-plugins --scope user
/plugin marketplace remove multipowers-plugins
```

## What To Verify After Sync

1. Overlay applies cleanly.
2. `/mp:persona list` works.
3. `/mp:init` uses `custom/config/setup.toml` and writes into target project `.multipowers/`.
4. Spec-driven outputs go under `.multipowers/tracks/<track_id>/`.
5. Auto-learning FAQ updates `.multipowers/FAQ.md` without writing to `$HOME` or tool project paths.

## Reference Docs

- [config-schema.md](/mnt/f/src/ai/multipowers/custom/docs/reference/config-schema.md)
- [compatibility.md](/mnt/f/src/ai/multipowers/custom/docs/reference/compatibility.md)
- [conductor-context.md](/mnt/f/src/ai/multipowers/custom/docs/customizations/conductor-context.md)
- [troubleshooting.md](/mnt/f/src/ai/multipowers/custom/docs/tool-project/troubleshooting.md)

## Tool Project Conductor Notes

Legacy tool-project Conductor context docs are stored at:
- [product-vision.md](/mnt/f/src/ai/multipowers/custom/docs/tool-project/product-vision.md)
- [product.md](/mnt/f/src/ai/multipowers/custom/docs/tool-project/product.md)
- [tech-stack.md](/mnt/f/src/ai/multipowers/custom/docs/tool-project/tech-stack.md)
