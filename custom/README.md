# Multipowers Overlay

This directory contains fork-specific customizations designed to minimize conflicts with upstream core files.

- `config/`: declarative policy (models, proxy, persona lane mapping)
- `commands/`: custom command source files
- `templates/`: source templates for generated `.multipowers/` artifacts
- `references/`: upstream source mapping and attribution notes
- `docs/`: operator documentation for this customization layer

Runtime policy:
- No shell runtime dependencies.
- Execution entrypoints are Go binaries:
  - `scripts/mp`
  - `go run ./cmd/octo-devx`
