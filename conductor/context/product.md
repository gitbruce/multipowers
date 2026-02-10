# Product: Multipowers Tool

## What We Build

Multipowers is a tooling layer for vibe-coding agents. It enables reliable iteration on a target project by combining:

- conductor-style setup and track lifecycle
- superpowers-style development methodology
- scene-aware role dispatch
- multi-model CLI execution (`claude`, `codex`, `gemini`)

## Core Boundaries

- `conductor/` in this repository is **for building Multipowers itself**.
- `templates/conductor/` is **for bootstrapping user target projects**.

## Primary Objectives

1. Keep stable project background in context docs and avoid repeating it in every task.
2. Route tasks through either a fast lane or a standard workflow lane.
3. Ensure role dispatch is explicit and observable at each execution node.
4. Require verification and documentation sync for significant changes.

## Non-Goals

- Replacing external model CLIs with embedded SDK coupling.
- Turning every tiny fix into heavyweight process.
- Allowing undocumented “done” claims without evidence.

## Acceptance Signals

- Maintainers can run setup + track commands and keep context stable.
- Standard workflows run as workflow graphs (not ad-hoc role jumping).
- Review nodes can enforce specialist roles (e.g., Architect) when required.
- Major changes leave artifact evidence: file list, checks run, docs updated.
