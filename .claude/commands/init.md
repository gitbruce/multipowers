---
command: init
description: Initialize conductor context in target project via interactive wizard
---

# /octo:init

This command MUST run an interactive wizard. Do not skip questions.

## Mandatory Contract

1. Target path is `$PWD/conductor`.
2. If all required files exist, return success and stop:
- `conductor/product.md`
- `conductor/product-guidelines.md`
- `conductor/tech-stack.md`
- `conductor/workflow.md`
- `conductor/tracks.md`
3. If any required file is missing, you MUST start wizard immediately with `AskUserQuestion`.
4. Do not continue to any planning/task command until required files are created.

## Wizard Flow (Required)

1. Ask mode:
- Header: `Init Mode`
- Question: `Conductor context is missing. How do you want to initialize it?`
- Options:
  - `Interactive (Recommended)` - Ask guided questions then generate files.
  - `Quick defaults` - Generate files with safe defaults.

2. If `Interactive`, ask one batched `AskUserQuestion` set for:
- product summary
- target users
- primary goal
- runtime/framework
- constraints

3. Create required files under `conductor/` and `conductor/code_styleguides/`:
- `product.md`
- `product-guidelines.md`
- `tech-stack.md`
- `workflow.md`
- `tracks.md`

4. Confirm completion by re-checking the five required files.
- If still missing, report failure and stop.
- If complete, report success and stop.

## Prohibited

- Do not bypass wizard with silent analysis-only behavior.
- Do not jump to `/octo:plan` questions before required files exist.
