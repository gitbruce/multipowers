# Product Vision: Build the Multipowers Tool

## Vision

Multipowers is an orchestration tool for AI coding agents. Its goal is to make target-project iteration reliable through multi-model collaboration with:

- stable context anchoring
- workflow enforcement
- role-based execution
- traceable verification and governance

## Product Positioning

- **Project type**: tooling project (not a business application)
- **Subject**: how to build the Multipowers tool
- **Audience**: tool maintainers
- **Tech stack**: Bash, Python, and Node.js tooling

## Strategic Inputs

1. From conductor: adopt `setup` and `new track` to stabilize project background and delivery tracks.
2. From superpowers: apply explicit methodology workflows for major changes.
3. From role-driven systems: route workflow nodes to the most suitable specialist role.
4. From multi-CLI bridging patterns: execute role calls through non-interactive `claude`, `codex`, and `gemini` CLIs.

## Design Principles

- **Workflow first, roles second**: in standard routing, choose a workflow before choosing node-level roles.
- **Fast for small changes, disciplined for major changes**: optimize both speed and quality.
- **Context and execution separation**: keep `conductor/context/*` stable and carry change details in tracks.
- **Evidence before claims**: every completion claim must include checks and documentation updates.
