# Multipowers Overlay

This directory contains fork-specific customizations designed to minimize conflicts with upstream core files.

- `config/`: declarative policy (models, proxy, persona lane mapping)
- `commands/`: custom command source files (for example `persona.md`)
- `lib/`: reusable helper libraries sourced by runtime and scripts
- `scripts/`: executable entry scripts (`apply-custom-overlay.sh`, `sync-upstream.sh`)
- `templates/`: source templates for generated `.multipowers/` context artifacts
- `references/`: upstream source mapping and attribution notes
- `docs/`: operator documentation for this customization layer
