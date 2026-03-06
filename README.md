# Multipowers

A Go-native multi-agent orchestration plugin for Claude Code.

> This repository is **based on** [nyldn/claude-octopus](https://github.com/nyldn/claude-octopus) and extends it with a Go-first runtime, stronger governance, and policy autosync learning.

## Table of Contents

- [Why This Fork](#why-this-fork)
- [Fork Differences (Commit-Based)](#fork-differences-commit-based)
- [Architecture Snapshot](#architecture-snapshot)
- [Requirements](#requirements)
- [Install](#install)
- [Quick Start](#quick-start)
- [Development](#development)
- [Testing](#testing)
- [Repository Layout](#repository-layout)
- [License](#license)

## Why This Fork

This fork focuses on runtime determinism, maintainability, and enterprise-style governance:

- Go-native orchestration and policy resolution on critical paths.
- Clear command boundary: `mp` (runtime) vs `mp-devx` (ops/devx).
- Structured governance checks and decision logs under `.multipowers`.
- Policy Auto Sync with prompt injection and deny confirmation flow.

## Fork Differences (Commit-Based)

Compared with upstream `main`, this fork branch (`go`) adds the following major differences based on commit history.

| Area | What changed in this fork | Example commits |
|---|---|---|
| Go-native orchestration runtime | Migrated orchestration planner/executor/synthesis to Go and aligned flow naming with runtime contracts. | `89a8440`, `bb982e0`, `e55a456` |
| Config-driven policy + dispatch | Added compiled runtime policy artifacts, config loaders/validators, and one-hop fallback dispatch behavior. | `e2fc72a`, `ddd46aa`, `41aaf93`, `23e6447`, `ef7e8b1` |
| Benchmark + smart routing | Added benchmark mode, async queue pipeline, JSONL persistence, judge scoring, and history-based routing override. | `d87f993`, `ff9d21a`, `4ba7861`, `88922be`, `a309b5b`, `34c16cc`, `b7c71e8` |
| Shared isolation + mailbox control plane | Added isolation policy/runtime, mailbox IPC, conflict monitor, deterministic gate decisions, and resource guardrails. | `36a5ee0`, `57f6a56`, `72f496e`, `8fef1bc`, `9e6b081`, `c5c7a3f`, `1b2037b`, `f2bd721`, `9233f58`, `43ebe08` |
| Runtime/ops command boundary | Formalized `mp` vs `mp-devx` ownership, migrated coverage/validate-runtime/cost-report to `mp-devx`. | `253ae8b`, `127f76a`, `52197ee`, `a704249`, `c1685ad` |
| Governance and doctor | Added/expanded governance checks, hook integration, and decision logging in `.multipowers`. | `53c1080`, `cb16789` |
| Policy Auto Sync learning loop | Added universal autosync domain (`internal/autosync`), deny confirmation (`delete` vs `skip-this-session`), and policy prompt injection (including external tool calls). | `2b59c12`, `830fe03`, `4981475`, `1c79154` |
| Naming/layout migration | Renamed and aligned plugin architecture toward `multipowers` + `.claude-plugin`/`.multipowers` conventions. | `ac17dfd`, `d21d73f` |

## Architecture Snapshot

High-level execution path in this fork:

1. User command enters `mp` runtime.
2. Policy resolver selects model/executor from compiled config.
3. Orchestration runs with optional isolation/mailbox control plane.
4. Hooks/doctor enforce governance and produce audit artifacts.
5. Policy Auto Sync ingests events, scores proposals, and injects active policy context into prompts.

Primary internal modules:

- `internal/orchestration`: planner, executor, synthesis loop.
- `internal/policy`: config compile/load/resolve/dispatch.
- `internal/hooks`: runtime hook enforcement.
- `internal/doctor`: governance diagnostics.
- `internal/autosync`: policy learning/injection/overlay controls.

## Requirements

- Claude Code `v2.1.34+`
- Go `1.22+` (for local build/test)
- Optional external providers:
  - Codex CLI (OpenAI)
  - Gemini CLI (Google)

## Install

Inside Claude Code:

```bash
/plugin marketplace add ${PWD}/.claude-plugin/marketplace.json
/plugin install mp@multipowers-plugins --scope user
```

Then run setup:

```bash
/mp:setup
```

## Quick Start

### Runtime commands (`mp`)

```bash
mp status --dir . --json
mp doctor --dir . --list
mp route --intent develop --dir . --json
mp policy sync --dir . --json
mp policy stats --dir . --json
```

### Workflow commands (plugin)

```text
/mp:discover <topic>
/mp:define <goal>
/mp:develop <task>
/mp:deliver <review-target>
/mp:embrace <end-to-end goal>
```

### Devx/Ops commands (`mp-devx`)

```bash
mp-devx --action doctor --dir . --verbose
mp-devx --action build-policy --config-dir config --output-dir .claude-plugin/runtime
mp-devx --action build-runtime
mp-devx --action init-fingerprint --dir . --json
```

## Development

```bash
go build ./cmd/mp
go build ./cmd/mp-devx
```

Build runtime artifacts:

```bash
mp-devx --action build-runtime
```

## Testing

Core suites used in this fork:

```bash
go test ./internal/autosync/... -v
go test ./internal/policy ./internal/hooks ./internal/doctor ./internal/cli -v
go test ./cmd/mp ./cmd/mp-devx -v
```

## Repository Layout

```text
cmd/
  mp/               # runtime CLI
  mp-devx/          # ops/devx CLI
internal/
  orchestration/    # planner/executor/synthesis
  policy/           # compile/load/resolve/dispatch
  hooks/            # runtime hook governance
  doctor/           # diagnostics checks
  autosync/         # policy learning + prompt injection
config/             # source policy/workflow/provider config
.claude-plugin/     # plugin manifest/runtime assets
.multipowers/       # runtime state and governance artifacts
```

## License

MIT. See [LICENSE](./LICENSE).
