# Workflow Skills: Quick Access to Multipowers Patterns

Multipowers includes **workflow skills** - lightweight wrappers that auto-invoke common multi-AI workflows. These activate automatically when you use certain natural language phrases.

## 🔬 Deep Research (`/mp:discover`)

**Auto-activates when you say:**
- "research this topic"
- "investigate how X works"
- "explore different approaches"

**What it does:** Runs the **Discover** workflow using the Go native engine.
- Spawns parallel agents (Codex, Gemini) to research from different perspectives.
- Claude synthesizes findings into a single Markdown report.

---

## 🔍 Quick Code Review (`/mp:deliver`)

**Auto-activates when you say:**
- "review this code"
- "check this PR"
- "quality check"

**What it does:** Runs the **Deliver** workflow.
- High-capability models perform architecture, security, and quality audits in parallel.
- Provides a weighted quality score and prioritized remediation steps.

---

## 🛡️ Adversarial Security (`/mp:debate`)

**Auto-activates when you say:**
- "security audit"
- "pentest this code"

**What it does:** Triggers a structured 3-way debate focused on security.
- **Red Team**: Attempts to find vulnerabilities and generate exploit PoCs.
- **Blue Team**: Reviews defenses and suggests patches.
- **Moderator**: Final synthesis and risk assessment.

---

## 🤖 Iterative Execution (`/mp:loop`)

**Auto-activates when you say:**
- "fix all lint errors"
- "refactor this until it passes tests"

**What it does:** Triggers the **Ralph Wiggum loop** engine.
- Execution continues iteratively until the task goal is reached or max iterations occur.
- Automatically handles linter or compiler feedback to self-correct.

---

## Architecture: How Skills work with the Go Engine

Understanding the distinction:

1. **Workflow Skills (Entry Points)**: Natural language triggers that make the tool easy to use.
2. **Go Engine (internal/orchestration)**: The robust, concurrent runtime that does the heavy lifting.
3. **Double Diamond (Methodology)**: The structured approach (Discover → Define → Develop → Deliver) enforced by the system.

---

[← Back to README](../README.md)
