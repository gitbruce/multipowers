# Multipowers Customization Hub

Start here for operator-focused guidance on this fork's custom overlay.

## Branch Principle
- `main` is an upstream mirror branch for `https://github.com/nyldn/claude-octopus/tree/main`.
- Do not put custom commits on `main`.
- Put all customization work on `multipowers`.
- Periodically merge `main` into `multipowers` with minimal-touch changes in upstream high-churn files.

## Directory Layout
- `custom/config/`: model, proxy, and persona lane configuration
- `custom/commands/`: command source overlays
- `custom/lib/`: internal helper libraries (not executed directly)
- `custom/scripts/`: executable scripts (`apply-custom-overlay.sh`, `sync-upstream.sh`)

## Contents
- getting-started.md
- customizations/*
- customizations/conductor-context.md
- sync/*
- troubleshooting.md
- reference/*
