# Multipowers: Multi-Model Orchestration Engine

Transform Claude Code into a role-driven, multi-model orchestration system with Conductor context anchoring.

## Quick Start

```bash
# Install dependencies
npm install

# First-time health check (strict mode expects missing context)
./bin/multipowers doctor

# Non-destructive bootstrap/repair
./bin/multipowers init --repair

# Fill required context files, then verify again
./bin/multipowers doctor

# Create a track
./bin/multipowers track new my-feature
```

## Architecture

### Conductor System
- **Context-First**: Immutable background in `conductor/context/*.md` is injected per task.
- **Track-Based**: Work is organized via `conductor/tracks/*.md`.
- **Role-Driven**: Specialized personas execute different lifecycle phases.

### Roles
- **Sisyphus**: Orchestrator & router (Claude Code main session)
- **Prometheus**: Architect & planner (Gemini)
- **Hephaestus**: TDD implementer (Codex)
- **Oracle**: Reviewer & verifier (Gemini)
- **Librarian**: Research role (Gemini Flash)

### Bridge Script
- `bin/ask-role` routes tasks to external CLIs
- loads `conductor/context/*.md` with priority-aware budget trimming
- uses `conductor/config/roles.json` first, then `config/roles.default.json`
- writes structured records for context preparation + connector execution
- propagates `request_id` across both log record types

## Directory Structure

```text
multipowers/
├── bin/
│   ├── multipowers              # CLI entry (init/doctor/update/track)
│   └── ask-role                 # Role dispatch bridge
├── config/
│   ├── roles.default.json
│   ├── roles.schema.json
│   └── mcp.default.json
├── connectors/
│   ├── codex.py
│   ├── gemini.py
│   └── utils.py
├── conductor/
│   ├── config/
│   ├── context/
│   └── tracks/
├── docs/
│   ├── design/
│   └── plans/
├── outputs/
│   └── runs/                    # structured JSONL runtime logs
├── scripts/
│   ├── validate_roles.py
│   ├── check_context_quality.py
│   └── check_plan_evidence.py
├── templates/
└── tests/
```

## Context Enforcement Mode

Use `MULTIPOWERS_CONTEXT_MODE` to align `doctor` and `ask-role` behavior:

- `strict` (default): missing required context files fails fast.
- `lenient`: missing context produces warnings and continues.

```bash
# strict (default)
./bin/multipowers doctor

# lenient mode
MULTIPOWERS_CONTEXT_MODE=lenient ./bin/multipowers doctor
MULTIPOWERS_CONTEXT_MODE=lenient ./bin/ask-role prometheus "analyze this"
```

## Init & Repair Modes

`bin/multipowers init` supports three modes:

- `./bin/multipowers init`: create `conductor/` once, no overwrite.
- `./bin/multipowers init --repair`: non-destructive; only fills missing files.
- `./bin/multipowers init --force --yes`: destructive re-create; explicit confirmation required.

## Development Workflow

### Setup
1. `./bin/multipowers init --repair`
2. Fill required context files:
   - `conductor/context/product.md`
   - `conductor/context/product-guidelines.md`
   - `conductor/context/workflow.md`
   - `conductor/context/tech-stack.md`
3. Optionally override roles in `conductor/config/roles.json`
4. Run `./bin/multipowers doctor`

### Track Lifecycle
1. `./bin/multipowers track new <feature-name>`
2. `./bin/multipowers track start <track-name>`
3. Execute design/planning/implementation
4. `./bin/multipowers track complete <track-name>`

### Role Dispatch
```bash
./bin/ask-role prometheus "Brainstorm architecture for auth"
./bin/ask-role hephaestus "Implement task with TDD"
./bin/ask-role oracle "Review diff for blocking issues"
```

## Path Conventions

- Design docs: `docs/design/YYYY-MM-DD-<feature>-design.md`
- Implementation plans: `docs/plans/YYYY-MM-DD-<feature>.md`
- Tracks: `conductor/tracks/track-YYYYMMDD-<feature>.md`
- Runtime logs: `outputs/runs/YYYY-MM-DD.jsonl`

## Context Budget

`ask-role` uses `MULTIPOWERS_CONTEXT_BUDGET` (default `128000` tokens) and trims context by file priority:
1. `conductor/context/product.md`
2. `conductor/context/product-guidelines.md`
3. `conductor/context/workflow.md`
4. `conductor/context/tech-stack.md`

When trimming occurs:
- stderr includes warning + truncated file list
- JSONL includes `event=context_prepared` with `truncated_files`

## Governance Modes

### Private Mode (default)
- `.gitignore` excludes `conductor/context/*.md` and `conductor/tracks/*.md`

### Traceable Mode
- uncomment ignore lines to version control context/tracks for team auditability

## Validation & Testing

```bash
# CLI health checks (context + schema + quality)
./bin/multipowers doctor

# Core tests
npm test --silent

# Fresh onboarding smoke test
bash tests/opencode/test-onboarding-smoke.sh

# Optional integration tests
bash tests/opencode/run-tests.sh --integration
```

## CI

### Core + Governance (`.github/workflows/core-tests.yml`)
- `core-tests` job:
  - `bash -n bin/multipowers bin/ask-role`
  - `python3 -m py_compile connectors/*.py scripts/*.py`
  - `npm test --silent`
- `governance-checks` job:
  - `python3 scripts/check_plan_evidence.py`
  - `python3 scripts/check_context_quality.py --context-dir templates/conductor/context --quiet`

### Integration (`.github/workflows/integration-tests.yml`)
- Triggered by `workflow_dispatch` and nightly `schedule`
- Runs `bash tests/opencode/run-tests.sh --integration`

## License

[Add license information here]
