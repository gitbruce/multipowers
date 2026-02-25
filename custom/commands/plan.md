---
command: plan
description: Build a plan and save spec tracking under conductor/tracks (not .claude)
aliases:
  - build-plan
  - intent
---

# Plan - Conductor Track First

When user runs `/octo:plan`, always execute:

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh plan "$ARGUMENTS"
```

Rules:
- Do not write `.claude/session-plan.md` or `.claude/session-intent.md`
- Save spec-driven planning artifacts in `conductor/tracks/`
- Use checkbox tracking (`- [ ]`, `- [x]`) for tasks/subtasks
- Read `conductor/*.md` before planning or execution starts
