---
command: init
description: Initialize .multipowers context in target project via interactive wizard
---

# /octo:init

This command MUST run an interactive wizard. Do not skip questions.

## Mandatory Contract

1. Target path is `$PWD/.multipowers`.
2. If all required files exist, return success and stop:
- `.multipowers/product.md`
- `.multipowers/product-guidelines.md`
- `.multipowers/tech-stack.md`
- `.multipowers/workflow.md`
- `.multipowers/tracks.md`
3. If any required file is missing, you MUST start wizard immediately with `AskUserQuestion`.
4. Do not continue to any planning/task command until required files are created.

## Wizard Flow (Required)

1. Ask mode:
- Header: `Init Mode`
- Question: `Conductor context is missing. How do you want to initialize it?`
- Options:
  - `Interactive (Required)` - Ask guided questions then generate files.

2. If `Interactive`, ask one batched `AskUserQuestion` set for:
- product summary
- target users
- primary goal
- runtime/framework
- constraints

3. Create required files under `.multipowers/` and `.multipowers/code_styleguides/`:
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
- Do not generate context files from non-interactive defaults.
- Do not jump to `/octo:plan` questions before required files exist.
