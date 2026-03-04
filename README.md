<p align="center">
  <img src="assets/social-preview.jpg" alt="Claude Octopus - Multi-tentacled orchestrator for Claude Code" width="640">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Claude_Code-Plugin-blueviolet" alt="Claude Code Plugin">
  <img src="https://img.shields.io/badge/Double_Diamond-Design_Thinking-orange" alt="Double Diamond">
  <img src="https://img.shields.io/badge/License-MIT-green" alt="MIT License">
  <img src="https://img.shields.io/badge/Version-8.12.17-blue" alt="Version 8.12.17">
  <img src="https://img.shields.io/badge/Claude_Code-v2.1.34+-blueviolet" alt="Requires Claude Code v2.1.34+">
</p>

# Claude Octopus

**Multi-AI orchestration plugin for Claude Code** - Run Codex, Gemini, and Claude simultaneously with 29 expert personas, Double Diamond workflows, and 44 specialized skills.

> *Three AI perspectives in the time it takes for one. Structured workflows that actually get followed.*

---

## Install

Inside Claude Code, run:

```
/plugin marketplace add /mnt/f/src/ai/claude-octopus/.claude-plugin/marketplace.json
/plugin install mp@multipowers-plugins --scope user
```

Then configure your AI providers:

```
/mp:setup
```

The setup wizard detects what you have installed, shows what's missing, and walks you through configuration. You only need **one** external provider (Codex or Gemini) to get multi-AI features - Claude is built-in.

**Requirements:** Claude Code v2.1.34+

---

## Quick Start

### For Developers

```
/mp:research OAuth 2.1 implementation patterns     # Multi-AI research
/mp:review                                          # Security-aware code review
/mp:tdd                                             # Red-green-refactor with discipline
/mp:debug                                           # Systematic 4-phase debugging
/mp:security                                        # OWASP vulnerability scan
/mp:persona list                                    # List available personas
/mp:persona security-auditor review auth flow       # Force specific persona
/mp:embrace build user authentication               # Full lifecycle: research to delivery
```

### For Knowledge Workers

```
/mp:research competitor landscape for B2B SaaS      # Multi-source synthesis
/mp:prd                                             # AI-optimized PRD with 100-point scoring
/mp:brainstorm                                      # Creative thought partner session
/mp:debate build vs buy for analytics platform      # Structured three-way AI debate
/mp:docs                                            # Export to PPTX, DOCX, PDF
/mp:embrace write market entry strategy             # Full lifecycle: research to deliverable
```

### The Smart Router

Don't remember the exact command? Just describe what you need:

```
/mp research microservices patterns      -> routes to discover phase
/mp build user authentication            -> routes to develop phase
/mp review this PR for security issues   -> routes to deliver phase
/mp compare Redis vs DynamoDB            -> routes to debate
```

The router parses your intent and selects the right workflow. Above 80% confidence it auto-routes; between 70-80% it confirms; below 70% it asks for clarification.

---

## How It Works

### Multi-AI Orchestration

Claude Octopus coordinates three AI providers - Codex, Gemini, and Claude - running them in parallel across every workflow. This isn't just a debate feature. Multi-AI orchestration powers the entire plugin:

- **`/mp:embrace`** runs a full 4-phase lifecycle where Codex and Gemini research independently in the Discover phase, build consensus in Define, propose competing implementations in Develop, then cross-review in Deliver
- **`/mp:extract`** uses Codex to analyze code structure while Gemini maps the design system, with Claude synthesizing tokens, components, and architecture into exportable formats
- **`/mp:research`** sends the same question to all providers simultaneously, then synthesizes three independent analyses into one report
- **`/mp:review`** has Codex check code quality and patterns while Gemini scans for security and edge cases, with Claude producing the final assessment
- **`/mp:debate`** structures a formal multi-round argument where each provider takes and defends a position

| Provider | Powered By | Role Across Workflows |
|----------|-----------|----------|
| Codex (OpenAI) | GPT-5.3-Codex | Implementation depth - code patterns, technical analysis, architecture proposals |
| Gemini (Google) | Gemini 3 Pro | Ecosystem breadth - alternative approaches, security review, research synthesis |
| Claude (Anthropic) | Sonnet 4.5 / Opus 4.6 | Orchestration and synthesis - quality gates, final recommendations, consensus building |

Each workflow uses providers differently. Research runs them in parallel. Define runs them sequentially for coherent problem scoping. Develop runs them in parallel for competing proposals, then merges through a 75% consensus quality gate. Deliver cross-validates with adversarial review.

**Graceful degradation:** Works with 1, 2, or 3 providers. With one external provider, you get dual-perspective analysis. With none, you still get all 29 personas, structured workflows, and every skill - multi-AI orchestration simply runs through Claude alone.

### Double Diamond Workflows

Four structured phases adapted from the UK Design Council's proven methodology:

```
  DISCOVER          DEFINE           DEVELOP          DELIVER
  (diverge)        (converge)       (diverge)        (converge)
     /\               /\               /\               /\
    /  \             /  \             /  \             /  \
   / Re \           / Sc \           / Bu \           / Va \
  / search\        / ope  \         / ild  \         / lidate\
 /________\       /________\       /________\       /________\
```

| Phase | Command | Alias | What Happens |
|-------|---------|-------|--------------|
| Discover | `/mp:discover` | `/mp:probe` | Multi-AI research and broad exploration |
| Define | `/mp:define` | `/mp:grasp` | Requirements clarification with consensus building |
| Develop | `/mp:develop` | `/mp:tangle` | Implementation with quality gates (75% threshold) |
| Deliver | `/mp:deliver` | `/mp:ink` | Adversarial review, security checks, go/no-go scoring |
| **All 4** | `/mp:embrace` | - | Complete lifecycle in one command |

Each phase has quality gates that must pass before proceeding. If a gate fails, the workflow pauses for revision rather than shipping questionable work.

Run phases individually or chain them. `/mp:embrace` runs all four in sequence, with configurable autonomy:

- **Supervised** (default) - Review and approve after each phase
- **Semi-autonomous** - Auto-proceed unless a quality gate fails
- **Autonomous** - Run all 4 phases without intervention

### 29 Expert Personas

Specialized AI agents that activate automatically based on your request. Each persona has domain expertise, a preferred AI provider, and memory that persists across sessions.

**Software Engineering** (11 personas)
backend-architect, frontend-developer, cloud-architect, devops-troubleshooter, deployment-engineer, database-architect, security-auditor, performance-engineer, code-reviewer, debugger, incident-responder

**Specialized Development** (6 personas)
ai-engineer, typescript-pro, python-pro, graphql-architect, test-automator, tdd-orchestrator

**Documentation & Communication** (5 personas)
docs-architect, product-writer, academic-writer, exec-communicator, content-analyst

**Research & Strategy** (4 personas)
research-synthesizer, ux-researcher, strategy-analyst, business-analyst

**Creative & Design** (3 personas)
thought-partner, mermaid-expert, context-manager

**How activation works:** Personas trigger proactively based on intent detection. When you say "audit my API for vulnerabilities," the security-auditor activates automatically. When you say "write a research paper," academic-writer takes over. No explicit invocation needed.

```
"I need a security audit of my auth code"       -> security-auditor persona
"Review my API design for scalability"           -> backend-architect persona
"Help me write a PRD for the new feature"        -> product-writer persona
"Research market sizing for AI developer tools"  -> strategy-analyst persona
"Create a sequence diagram for the auth flow"    -> mermaid-expert persona
```

### Context-Aware Intelligence

Claude Octopus auto-detects whether you're doing development work or knowledge work and adapts everything: research sources, output formats, review criteria, and persona selection.

**Dev mode** (activates in code repositories): Research targets libraries and patterns. Output is code and tests. Reviews check security and performance.

**Knowledge mode** (`/mp:km on`): Research targets market data and strategy. Output is PRDs and reports. Reviews check clarity and evidence quality.

Auto-detection uses file signatures - `package.json` triggers dev mode, business keywords trigger knowledge mode. Override anytime with `/mp:km on|off|auto` or `/mp:dev`.

---

## Developer Workflows

### Code Review

```
/mp:review
```

Multi-perspective code review combining Codex (code quality, patterns), Gemini (security, edge cases), and Claude (synthesis, recommendations). Checks architecture, security vulnerabilities, performance bottlenecks, and maintainability.

### Test-Driven Development

```
/mp:tdd
```

Enforces red-green-refactor discipline. Write failing tests first, implement minimally to pass, then refactor with confidence. The TDD orchestrator prevents skipping steps.

### Debugging

```
/mp:debug
```

Systematic 4-phase debugging: Investigate (gather evidence), Analyze (form hypotheses), Hypothesize (rank causes), Implement (fix and verify). No more random `console.log` scattering.

### Security Audit

```
/mp:security
```

OWASP Top 10 compliance checking, vulnerability detection, dependency scanning, and adversarial security testing. The security-auditor persona brings specialized knowledge of attack vectors.

### Design System Extraction

```
/mp:extract ./my-app                                    # Interactive mode
/mp:extract ./my-app --mode design --storybook true     # Design system with Storybook
/mp:extract ./my-app --depth deep --multi-ai force      # Deep analysis, all providers
/mp:extract https://example.com --mode design           # From live website
```

Reverse-engineers design tokens (W3C format), components (React/Vue/Svelte), architecture (service boundaries, API contracts), and features. Outputs JSON, CSS, Markdown, CSV.

---

## Knowledge Worker Workflows

### Deep Research

```
/mp:research competitor landscape for B2B SaaS tools
```

Multi-source synthesis combining Codex (technical analysis), Gemini (ecosystem research), and Claude (strategic synthesis). The research skill asks 3 clarifying questions (depth, focus, format) before execution, so you get exactly what you need.

### PRD Writing

```
/mp:prd
```

Write AI-optimized PRDs scored against a 100-point framework. Structures requirements in sequential phases with P0/P1/P2 priority levels and explicit boundary definitions. Score existing PRDs with `/mp:prd-score`.

### AI Debate

```
/mp:debate build vs buy for analytics platform
```

Structured three-way debate between Codex, Gemini, and Claude. Each takes a position, provides evidence, and responds to counterarguments. Multiple styles: quick (1 round), thorough (2-3 rounds), adversarial (active critique), or collaborative (build on ideas).

### Brainstorming

```
/mp:brainstorm
```

Creative thought partner session using Pattern Spotting, Paradox Hunting, Naming the Unnamed, and Contrast Creation techniques. Helps surface hidden insights and unexpected connections.

### Document Delivery

```
/mp:docs
```

Export your work to PPTX, DOCX, or PDF. Converts markdown deliverables into professional document formats ready for stakeholder review.

### Content Analysis

```
/mp:pipeline https://example.com/article
```

Multi-stage content analysis pipeline. Reverse-engineers article anatomy, extracts recreatable patterns and frameworks, identifies psychological techniques and structural elements.

---

## All Commands

### Core Workflows
| Command | Description |
|---------|-------------|
| `/mp:embrace` | Full Double Diamond workflow (all 4 phases) |
| `/mp:discover` | Discovery phase - multi-AI research |
| `/mp:define` | Definition phase - requirements and scope |
| `/mp:develop` | Development phase - implementation with quality gates |
| `/mp:deliver` | Delivery phase - review and validation |
| `/mp:research` | Deep research with multi-source synthesis |

### Development
| Command | Description |
|---------|-------------|
| `/mp:tdd` | Test-driven development (red-green-refactor) |
| `/mp:debug` | Systematic debugging with methodical investigation |
| `/mp:review` | Expert code review with security analysis |
| `/mp:security` | OWASP compliance and vulnerability detection |
| `/mp:quick` | Fast execution without full workflow overhead |

### AI & Decisions
| Command | Description |
|---------|-------------|
| `/mp:debate` | Structured three-way AI debate |
| `/mp:loop` | Iterate until exit criteria pass |
| `/mp:brainstorm` | Creative thought partner session |
| `/mp:meta-prompt` | Generate optimized prompts |
| `/mp:multi` | Force multi-provider execution (manual override) |

### Planning & Docs
| Command | Description |
|---------|-------------|
| `/mp:prd` | AI-optimized PRD writing |
| `/mp:prd-score` | Score PRDs against 100-point framework |
| `/mp:plan` | Strategic plan builder (doesn't execute) |
| `/mp:docs` | Export to PPTX, DOCX, PDF |
| `/mp:pipeline` | Content analysis and pattern extraction |
| `/mp:extract` | Design system and product reverse-engineering |

### Project Lifecycle
| Command | Description |
|---------|-------------|
| `/mp:status` | Project progress dashboard |
| `/mp:resume` | Restore context from previous session |
| `/mp:ship` | Finalize with multi-AI validation |
| `/mp:issues` | Cross-session issue tracking |
| `/mp:rollback` | Checkpoint recovery (git tags) |

### Mode & Configuration
| Command | Description |
|---------|-------------|
| `/mp:km` | Toggle Knowledge Work mode |
| `/mp:dev` | Switch to Dev Work mode |
| `/mp:model-config` | Configure AI provider models at runtime |
| `/mp:persona` | Run a specific persona or list configured personas |
| `/mp:setup` | Provider setup wizard |
| `/mp:sys-setup` | System configuration status |

### Phase Aliases
| Command | Same as |
|---------|---------|
| `/mp:probe` | `/mp:discover` |
| `/mp:grasp` | `/mp:define` |
| `/mp:tangle` | `/mp:develop` |
| `/mp:ink` | `/mp:deliver` |

---

## 43 Skills

Skills are the engine behind commands and personas. They activate automatically when needed - you don't invoke them directly.

**Workflow Phases** - flow-discover, flow-define, flow-develop, flow-deliver

**Research & Knowledge** - skill-deep-research, skill-debate, skill-debate-integration, skill-thought-partner, skill-meta-prompt, skill-knowledge-work

**Code Quality** - skill-code-review, skill-quick-review, skill-security-audit, skill-adversarial-security, skill-security-framing, skill-audit

**Development** - skill-tdd, skill-debug, skill-verify, skill-validate, skill-iterative-loop, skill-finish-branch, skill-parallel-agents

**Architecture & Planning** - skill-architecture, skill-prd, skill-writing-plans, skill-decision-support, skill-intent-contract

**Content & Docs** - skill-doc-delivery, skill-content-pipeline, skill-visual-feedback

**Project Lifecycle** - skill-status, skill-issues, skill-rollback, skill-resume, skill-resume-enhanced, skill-ship

**Task & Session** - skill-task-management, skill-task-management-v2, skill-quick

**Mode & Config** - skill-context-detection, sys-configure, extract-skill

**How skills relate to commands:** Commands are what you type. Skills are what runs. When you run `/mp:review`, it activates the skill-code-review skill, which invokes the code-reviewer persona, which routes to the appropriate AI providers. You interact with commands; skills handle execution.

---

## Project Lifecycle

Track state across sessions with the `.claude-octopus/` directory:

```
.claude-octopus/
├── PROJECT.md      # Vision and requirements
├── ROADMAP.md      # Phase breakdown
├── STATE.md        # Current position and history
├── config.json     # Workflow preferences
├── ISSUES.md       # Cross-session issue tracking
└── LESSONS.md      # Lessons learned (preserved across rollbacks)
```

Created automatically on first `/mp:embrace`. Use `/mp:status` for a progress dashboard, `/mp:resume` to continue where you left off, `/mp:issues` to track cross-session problems, and `/mp:rollback` to restore from git tag checkpoints.

LESSONS.md is intentionally preserved across rollbacks - mistakes are worth remembering.

---

## Model Configuration

Configure which AI models power each provider:

```
/mp:model-config
```

Supports runtime model selection with 4-tier precedence:
1. Environment variables (`OCTOPUS_CODEX_MODEL`, `OCTOPUS_GEMINI_MODEL`)
2. Runtime overrides
3. Config file settings
4. Built-in defaults (GPT-5.3-Codex, Gemini 3 Pro, Claude Sonnet 4.5)

For premium tasks, complexity-based routing automatically upgrades to Opus 4.6.

Optional proxy settings for Codex/Gemini can be added in `~/.claude-octopus/config/providers.json`:

```json
{
  "proxy": {
    "enabled": true,
    "port": 7890
  }
}
```

---

## Benchmark + Smart Routing (Optional)

You can enable async benchmark collection for `/mp:*` code-related requests, then optionally let history override routing.

Configure in `config/orchestration.yaml`:

```yaml
benchmark_mode:
  enabled: true
  async_enabled: true
  force_all_models_on_code: true
  judge_model: "claude-opus"

smart_routing:
  enabled: false
  min_samples_per_model: 10
```

Behavior:
- `benchmark_mode.enabled=true`: code-related `/mp:*` requests can fan out to all available models and write daily JSONL records under `~/.claude-octopus/metrics`.
- `smart_routing.enabled=false`: no history override is applied.
- `smart_routing.enabled=true`: override is applied only when a similar-scenario model has at least `min_samples_per_model` judged samples.
- Benchmark queue/store/judge failures are best-effort only and do not fail the main workflow result.

---

## Cost Transparency

You see cost estimates **before** execution. Interactive research asks 3 questions (depth, focus, format) then shows exactly what will run and how much it costs.

**Most users pay nothing extra.** If you authenticate Codex via `codex login` (ChatGPT account) and Gemini via Google account OAuth, usage is covered by your existing subscriptions - no per-token API charges. Claude is included with your Claude Code subscription.

| Auth Method | Codex Cost | Gemini Cost | Claude Cost |
|-------------|-----------|-------------|-------------|
| OAuth login (recommended) | Included in ChatGPT Plus/Pro/Team | Included in Google AI subscription | Included in Claude Code |
| API key (`OPENAI_API_KEY` / `GEMINI_API_KEY`) | ~$0.02-0.10/query | ~$0.01-0.03/query | Included in Claude Code |

If you do use API keys, here are the per-workflow estimates:

| Scenario | Time | Est. API Cost |
|----------|------|---------------|
| Quick research | 1-2 min | $0.01-0.02 |
| Standard research | 2-3 min | $0.02-0.05 |
| Deep dive | 4-5 min | $0.05-0.10 |
| AI debate | 5-10 min | $0.08-0.15 |
| Code review | 3-5 min | $0.04-0.08 |
| Full workflow | 15-25 min | $0.20-0.40 |

Works without external providers too - you still get 29 personas, all workflows, context-aware intelligence, and every skill. Multi-AI features activate only when providers are available.

---

## FAQ

**Do I need all three AI providers?**
No. One external provider (Codex or Gemini) plus the built-in Claude gives you multi-AI features. Both external providers gives maximum diversity. No external providers still gives you personas, workflows, and skills.

**Will this break my existing Claude Code setup?**
No. Claude Octopus only activates with the `mp` prefix or `/mp:*` commands. Results are stored separately in `~/.claude-octopus/`. Uninstalls cleanly with no residual configuration changes.

**Can I use it without external AIs?**
Yes. You get all 29 personas, structured workflows, context intelligence, task management, and every skill. Multi-AI features (parallel analysis, debate, consensus) won't activate without external providers.

**What happens if an external provider times out?**
The workflow continues with available providers. If Codex fails, Gemini and Claude complete the work. If both fail, Claude handles it solo. You'll see the provider status in the visual indicators.

**What's the difference between `/mp:quick` and full workflows?**
`/mp:quick` skips the structured phases and quality gates - it's a fast path for ad-hoc tasks that don't need the full Double Diamond treatment. Use it for simple tasks; use `/mp:embrace` for complex features.

**Can I share `.claude-octopus/` state across a team?**
Yes. The `.claude-octopus/` directory is designed to be committed to your repository. Team members can use `/mp:resume` to pick up where others left off. ISSUES.md and LESSONS.md provide cross-session continuity.

**Does this work offline?**
Partially. Claude (via Claude Code) works with your subscription. External providers (Codex, Gemini) require internet access. All personas and workflow logic run locally.

**How do I see what's happening under the hood?**
Visual indicators show active providers in real-time. For deeper debugging, check logs in `~/.claude-octopus/logs/` or enable debug mode.

**How do I update?**
Run `/plugin` > Installed > update, or reinstall:
```
/plugin uninstall claude-octopus@multipowers-plugins
/plugin install mp@multipowers-plugins --scope user
```

---

## Documentation

- [Visual Indicators](docs/VISUAL-INDICATORS.md) - Understanding provider status
- [Command Reference](docs/COMMAND-REFERENCE.md) - All commands in detail
- [Architecture](docs/ARCHITECTURE.md) - How it works internally
- [Plugin Architecture](docs/PLUGIN-ARCHITECTURE.md) - Plugin structure
- [Native Integration](docs/NATIVE-INTEGRATION.md) - Claude Code TaskCreate/TaskUpdate
- [Debug Mode](docs/DEBUG_MODE.md) - Troubleshooting workflows
- [Full Changelog](CHANGELOG.md) - Complete version history

---

## Attribution

- **[wolverin0/claude-skills](https://github.com/wolverin0/claude-skills)** - AI Debate Hub for structured three-way debates. MIT License.
- **[UK Design Council](https://www.designcouncil.org.uk/our-resources/the-double-diamond/)** - Double Diamond methodology.

---

## Contributing

1. [Report Issues](https://github.com/nyldn/claude-octopus/issues)
2. Submit PRs following existing code style
3. Development: `git clone --recursive https://github.com/nyldn/claude-octopus.git && make test`

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

MIT - see [LICENSE](LICENSE)

<p align="center">
  <a href="https://github.com/nyldn">nyldn</a> | MIT License | <a href="https://github.com/nyldn/claude-octopus/issues">Report Issues</a>
</p>
