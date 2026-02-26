---
name: model-config
description: Configure AI provider models for Claude Octopus workflows
version: 2.0.0
category: configuration
tags: [config, models, providers, codex, gemini, routing]
created: 2025-01-21
updated: 2026-02-14
---

# Model Configuration

Configure which AI models are used by Claude Octopus workflows. This allows you to:
- Configure Codex and Gemini defaults/overrides
- Configure Codex phase routing (`phase_routing`) for Codex-backed steps
- Align model choices with the project routing policy
- Control cost/performance tradeoffs per project

## Environment Constraints

In this environment:
- Codex available model: `gpt-5.3-codex` only
- Gemini available model: `gemini-3-pro-preview`
- Claude routing uses Sonnet/Opus lanes and is mapped by your Claude Code defaults

## Principles-Aligned Policy

Current project policy (aligned with `agents/config.yaml`, `workflows/embrace.yaml`, and `bin/octo`):

- Planning, architecture, and important decisions -> Codex (`gpt-5.3-codex`)
- Heavy coding/implementation -> Claude Opus (`claude-opus`; mapped by your Claude Code env, e.g. GLM-5)
- Documentation and test-case authoring -> Claude Sonnet (`claude`; mapped by your Claude Code env, e.g. GLM-4.7)
- External-world research -> Gemini (`gemini-3-pro-preview`)
- Quality checks:
  - Heavy/high-token -> Claude Opus
  - Light/lower-token -> Codex

Important scope note:
- `/octo:model-config` configures Codex/Gemini models in `~/.claude-octopus/config/providers.json`.
- Claude model family selection is primarily controlled by routing (`claude` vs `claude-opus`) and your Claude Code defaults (`ANTHROPIC_DEFAULT_*`), not by this command.

## Usage

```bash
# View current configuration (models + phase routing)
/octo:model-config

# Set codex model (persistent)
/octo:model-config codex gpt-5.3-codex

# Set gemini model (persistent)
/octo:model-config gemini gemini-3-pro-preview

# Set session-only override (doesn't modify config file)
/octo:model-config codex gpt-5.3-codex --session

# Configure phase routing (which codex model to use in which phase)
/octo:model-config phase deliver gpt-5.3-codex
/octo:model-config phase develop gpt-5.3-codex

# Reset to defaults
/octo:model-config reset codex
/octo:model-config reset phases
/octo:model-config reset all
```

## Model Precedence

Models are selected using a 5-tier precedence system:

1. **Environment variables** (highest priority)
   - `OCTOPUS_CODEX_MODEL` - Override all codex model selection
   - `OCTOPUS_GEMINI_MODEL` - Override all gemini model selection

2. **Task hints** (contextual override from calling code)
   - Task hints still exist in runtime, but with your setup Codex stays on `gpt-5.3-codex`

3. **Phase routing config** (per-phase model selection)
   - Stored in `~/.claude-octopus/config/providers.json` → `phase_routing`

4. **Config file defaults / session overrides**
   - Stored in `~/.claude-octopus/config/providers.json` → `providers` / `overrides`

5. **Hard-coded defaults** (lowest priority)
   - Codex: `gpt-5.3-codex`
   - Gemini: `gemini-3-pro-preview`

## Supported Models

### Codex (Available in this Environment)

| Model | Context | Speed | Best For | Cost |
|-------|---------|-------|----------|------|
| `gpt-5.3-codex` | 400K | ~65 tok/s | Planning, architecture, lighter quality checks | $1.75/$14.00 per MTok |

### OpenRouter Models (v8.11.0)

| Agent Type | Model | Context | Best For | Cost |
|------------|-------|---------|----------|------|
| `openrouter-glm5` | `z-ai/glm-5` | 203K | Code review (77.8% SWE-bench, lowest hallucination) | $0.80/$2.56 per MTok |
| `openrouter-kimi` | `moonshotai/kimi-k2.5` | **262K** | Research, large context, multimodal | $0.45/$2.25 per MTok |
| `openrouter-deepseek` | `deepseek/deepseek-r1` | 164K | Deep reasoning (visible `<think>` traces) | $0.70/$2.50 per MTok |

Requires `OPENROUTER_API_KEY` to be set. These are automatically selected when OpenRouter is the chosen provider via `get_tiered_agent_v2()` task routing.

### Gemini (Google)

| Model | Best For | Cost |
|-------|----------|------|
| `gemini-3-pro-preview` | Premium quality research | $2.50/$10.00 per MTok |
| `gemini-3-flash-preview` | Fast, low-cost tasks | $0.25/$1.00 per MTok |

## Phase Routing (Codex Only)

`phase_routing` controls Codex model selection for Codex-backed steps only.

| Phase | Recommended Codex Model | Rationale |
|-------|--------------------------|-----------|
| `discover` | `gpt-5.3-codex` | Technical analysis depth |
| `define` | `gpt-5.3-codex` | Planning and architecture decisions |
| `develop` | `gpt-5.3-codex` | Codex-side integration/review within develop |
| `deliver` | `gpt-5.3-codex` | Lighter quality checks on Codex lane |
| `review` | `gpt-5.3-codex` | Explicit light review lane |
| `security` | `gpt-5.3-codex` | Codex security pass (non-heavy path) |
| `research` | `gpt-5.3-codex` | Codex technical synthesis path |

### Customizing Phase Routing

```bash
# Use full codex for review/decision phases
/octo:model-config phase deliver gpt-5.3-codex
/octo:model-config phase review gpt-5.3-codex

# Reset phase routing to defaults
/octo:model-config reset phases
```

## Examples

### Planning/Architecture First (Codex)
```bash
/octo:model-config codex gpt-5.3-codex
/octo:model-config phase define gpt-5.3-codex
/octo:model-config phase review gpt-5.3-codex
```

### External Research (Gemini)
```bash
/octo:model-config gemini gemini-3-pro-preview
```

### Claude Lane Reminder
```bash
# Claude Sonnet/Opus selection is controlled by routing and Claude Code defaults:
# ANTHROPIC_DEFAULT_SONNET_MODEL / ANTHROPIC_DEFAULT_OPUS_MODEL
# Use persona and flow routing for heavy-vs-light Claude behavior.
```

## Configuration File

Location: `~/.claude-octopus/config/providers.json`

```json
{
  "version": "2.0",
  "providers": {
    "codex": {
      "model": "gpt-5.3-codex",
      "fallback": "gpt-5.3-codex"
    },
    "gemini": {
      "model": "gemini-3-pro-preview",
      "fallback": "gemini-3-pro-preview"
    }
  },
  "phase_routing": {
    "discover": "gpt-5.3-codex",
    "define":   "gpt-5.3-codex",
    "develop":  "gpt-5.3-codex",
    "deliver":  "gpt-5.3-codex",
    "review":   "gpt-5.3-codex",
    "security": "gpt-5.3-codex",
    "research": "gpt-5.3-codex"
  },
  "overrides": {}
}
```

## Codex Model Guidance (This Environment)

| Model | Best For |
|-------|----------|
| `gpt-5.3-codex` | Planning, architecture, important decisions, lighter quality checks |

For heavy implementation and heavy-token quality checks, rely on Claude Opus routing instead of forcing Codex.

## Requirements

- `python3` - Used for JSON read/write operations

## Notes

- Model names are validated against known models but unknown models are still accepted
- Invalid models will fail when workflows execute
- Environment variables override all other settings including phase routing
- Phase routing only affects Codex model selection (Gemini has its own model defaults; Claude family is selected by routing plus Claude Code defaults)
- Cost implications vary significantly between models - see pricing table above
- **Gemini sandbox modes** (`OCTOPUS_GEMINI_SANDBOX`):
  - `headless` (default, v8.10.0) - Stdin-based prompt delivery with `-p ""`, `-o text`, `--approval-mode yolo`
  - `interactive` - Launch Gemini in interactive mode (for manual use)
  - `auto-accept` - Legacy alias for `headless`
  - `prompt-mode` - Legacy alias for `interactive`

---

## EXECUTION CONTRACT (Mandatory)

When the user invokes `/octo:model-config`, you MUST:

1. **Parse arguments** to determine action:
   - No args → View current configuration including phase routing
   - `<provider> <model>` → Set model (persistent)
   - `<provider> <model> --session` → Set model (session only)
   - `phase <phase> <model>` → Set phase-specific model routing
   - `reset <provider|phases|all>` → Reset to defaults

2. **View Configuration** (no args):
   ```bash
   # Check environment variables
   env | grep OCTOPUS_

   # Show config file contents
   if [[ -f ~/.claude-octopus/config/providers.json ]]; then
     cat ~/.claude-octopus/config/providers.json
   else
     echo "No configuration file found (using defaults)"
   fi
   ```

3. **Set Model** (`<provider> <model>` or with `--session`):
   ```bash
   # Call set_provider_model from orchestrate.sh
   source "${CLAUDE_PLUGIN_ROOT}/bin/octo"
   set_provider_model <provider> <model> [--session]

   # Show updated configuration
   cat ~/.claude-octopus/config/providers.json
   ```

4. **Set Phase Routing** (`phase <phase> <model>`):
   ```bash
   # Update phase_routing in config file
   local config_file="${HOME}/.claude-octopus/config/providers.json"
   python3 - "$config_file" "$phase" "$model" <<'PY'
import json, sys
cfg, phase, model = sys.argv[1], sys.argv[2], sys.argv[3]
with open(cfg, "r", encoding="utf-8") as f:
    data = json.load(f)
data.setdefault("phase_routing", {})[phase] = model
with open(cfg + ".tmp", "w", encoding="utf-8") as f:
    json.dump(data, f, indent=2, ensure_ascii=False)
    f.write("\n")
PY
   mv "${config_file}.tmp" "$config_file"
   echo "✓ Set phase routing: $phase → $model"
   python3 - <<'PY'
import json, os
cfg=os.path.expanduser("~/.claude-octopus/config/providers.json")
with open(cfg, "r", encoding="utf-8") as f:
    print(json.dumps(json.load(f).get("phase_routing", {}), indent=2, ensure_ascii=False))
PY
   ```

5. **Reset Model** (`reset <provider|phases|all>`):
   ```bash
   # Call reset_provider_model from orchestrate.sh
   source "${CLAUDE_PLUGIN_ROOT}/bin/octo"
   reset_provider_model <provider>

   # For phases: reset phase_routing to defaults
   # For all: reset both providers and phase routing

   # Show updated configuration
   cat ~/.claude-octopus/config/providers.json
   ```

6. **Provide guidance** on:
   - Which models are appropriate for which tasks/phases
   - Cost implications of available models and routing choices
   - How to use environment variables for temporary changes
   - Claude Sonnet/Opus routing behavior in this environment

### Validation Gates

- Parsed arguments correctly
- Action determined (view/set/set-phase/reset)
- Functions called with Bash tool (not simulated)
- Configuration displayed to user
- Clear confirmation messages shown

### Prohibited Actions

- Assuming configuration without reading the file
- Suggesting edits without using the provided functions
- Skipping validation of provider names
- Ignoring errors from JSON updates or function calls
