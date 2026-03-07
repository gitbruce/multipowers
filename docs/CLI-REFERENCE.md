# CLI Reference

## Mainline runtime commands

```bash
mp init --dir <project> --prompt '{...}' --json
mp brainstorm --dir <project> --prompt "..." --json
mp design --dir <project> --prompt "..." --json
mp plan --dir <project> --prompt "..." --json
mp execute --dir <project> --prompt "..." --json
mp debug --dir <project> --prompt "..." --json
mp debate --dir <project> --prompt "..." --json
```

## Operational commands

```bash
mp status --dir <project> --json
mp doctor --dir <project> --json
mp resume --dir <project> --json
mp setup --dir <project> --json
mp config show-model-routing --dir <project> --json
```

## Init gating

All mainline and exception-path commands require initialized context. When context is missing, the CLI returns a blocked response with:

- `recommended_command=/mp:init`
- `resume_command=<original command>`
- `resume_prompt=<original prompt>`

## Devx commands

```bash
mp-devx --action sync-superpowers
mp-devx --action build-policy --config-dir config --output-dir .claude-plugin/runtime
mp-devx --action build-runtime
mp-devx --action doctor --dir . --json
```
