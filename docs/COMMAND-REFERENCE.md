# Command Reference

## Public commands

| Command | Purpose |
|---|---|
| `/mp:init` | Create the required `.multipowers` initialization artifacts. |
| `/mp:model-config` | Change user-specific model configuration. |
| `/mp:brainstorm` | Start exploration using the upstream brainstorming workflow. |
| `/mp:design` | Turn exploration into a concrete design direction. |
| `/mp:plan` | Write the implementation plan. |
| `/mp:execute` | Execute the plan and finish the branch inside the same flow. |
| `/mp:debug` | Enter the direct debugging path. |
| `/mp:debate` | Run all configured models in a structured debate. |
| `/mp:status` | Inspect current mainline state and track progress. |
| `/mp:doctor` | Diagnose runtime and environment readiness. |
| `/mp:resume` | Resume the last interrupted mainline command. |
| `/mp:setup` | Prepare or verify setup required by the mainline flow. |

## Mainline rule

Every mainline and special-entry command checks for `/mp:init` artifacts first. If they are missing, the runtime blocks, sends the user to `/mp:init`, and preserves resume metadata.

## Command intent

- `brainstorm`: open the problem space
- `design`: choose a direction
- `plan`: turn the direction into executable work
- `execute`: implement and finish
- `debug`: diagnose directly
- `debate`: compare all configured models
