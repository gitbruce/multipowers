# Skills Update for Role-Based Workflow Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Update all skills in `/skills` folder to be concise and specify main role per skill with workflow steps that may involve different roles from `config/roles.default.json`

**Architecture:** Read each skill, identify main role (router/architect/coder/librarian), add role specification to workflow steps that require different roles, remove redundant content while preserving essential information

**Tech Stack:** Bash (git), Markdown (skill files), JSON (role config)

---

## Context Summary

**Available Roles** (from `config/roles.default.json`):
- **router**: Orchestrator & Router - routes work by task type and workflow
- **architect**: Planning, architecture, and review/verification (Gemini 1.5 Pro)
- **coder**: TDD implementation executor (Deepseek Coder)
- **librarian**: Research and documentation (Gemini Flash)

**Key Requirements:**
1. Make skill content concise (remove redundancy)
2. Each skill has a main role
3. Workflow steps may specify different roles when needed
4. Roles must come from `config/roles.default.json`

---

## Task 1: Update brainstorming/SKILL.md

**Files:**
- Modify: `skills/brainstorming/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/brainstorming/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **architect** (architecture design)

Current Conductor Integration Notes already references architect role.

**Step 3: Simplify and add role specifications to workflow**

Update skill structure:
- Keep frontmatter (name, description)
- Keep "Overview" (1-2 sentences max)
- Keep "The Process" but add role labels to each step
- Remove redundant examples
- Keep Conductor Integration Notes with explicit role mapping

**Step 4: Verify changes**

```bash
git diff skills/brainstorming/SKILL.md
```

Expected: Content is 30-40% shorter, roles specified for workflow steps

**Step 5: Commit**

```bash
git add skills/brainstorming/SKILL.md
git commit -m "refactor: brainstorming skill - concise with role specifications"
```

---

## Task 2: Update writing-plans/SKILL.md

**Files:**
- Modify: `skills/writing-plans/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/writing-plans/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **architect** (implementation planning)

Current Conductor Integration Notes references architect role.

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Compress "Overview" to core principle only
- Keep "Bite-Sized Task Granularity" (essential)
- Keep "Plan Document Header" template (essential)
- Simplify "Task Structure" example
- Compress "Remember" to bullet points
- Keep "Execution Handoff" but mark roles for each option
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/writing-plans/SKILL.md
```

Expected: ~40% shorter, roles specified in execution handoff

**Step 5: Commit**

```bash
git add skills/writing-plans/SKILL.md
git commit -m "refactor: writing-plans skill - concise with role specifications"
```

---

## Task 3: Update systematic-debugging/SKILL.md

**Files:**
- Modify: `skills/systematic-debugging/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/systematic-debugging/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **coder** (primary debugger), but research steps use **librarian**

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Compress "Overview" to core principle
- Keep "The Iron Law" (essential)
- Keep "When to Use" as bullet list
- Keep "Four Phases" but add role labels:
  - Phase 1: Root Cause Investigation (main: coder)
  - Phase 2: Pattern Analysis (coder)
  - Phase 3: Hypothesis and Testing (coder)
  - Phase 4: Implementation (coder)
- Keep "Red Flags" list
- Compress "Quick Reference" table
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/systematic-debugging/SKILL.md
```

Expected: ~30% shorter, role labels on phases

**Step 5: Commit**

```bash
git add skills/systematic-debugging/SKILL.md
git commit -m "refactor: systematic-debugging skill - concise with role specifications"
```

---

## Task 4: Update test-driven-development/SKILL.md

**Files:**
- Modify: `skills/test-driven-development/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/test-driven-development/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **coder** (TDD implementation)

**Step 3: Simplify while preserving TDD discipline**

Update skill structure:
- Keep frontmatter
- Compress "Overview" to core principle
- Keep "The Iron Law" (essential)
- Keep "Red-Green-Refactor" flowchart
- Compress each phase description
- Keep "Good Tests" comparison table
- Compress "Common Rationalizations" table
- Keep "Verification Checklist"
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/test-driven-development/SKILL.md
```

Expected: ~35% shorter, TDD discipline preserved

**Step 5: Commit**

```bash
git add skills/test-driven-development/SKILL.md
git commit -m "refactor: test-driven-development skill - concise with role specifications"
```

---

## Task 5: Update subagent-driven-development/SKILL.md

**Files:**
- Modify: `skills/subagent-driven-development/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/subagent-driven-development/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **router** (orchestrates workflow, delegates to coder/architect)

Workflow involves:
- Router: coordination
- Coder: implementation nodes
- Architect: spec compliance and code quality review nodes

**Step 3: Simplify and add explicit role mapping**

Update skill structure:
- Keep frontmatter
- Compress "Overview"
- Keep "When to Use" flowchart
- Compress "The Process" flowchart description
- Add explicit role mapping table:
  - Implementer subagent → coder
  - Spec reviewer → architect
  - Code quality reviewer → architect
- Compress "Example Workflow"
- Keep Conductor Integration Notes with expanded role mapping

**Step 4: Verify changes**

```bash
git diff skills/subagent-driven-development/SKILL.md
```

Expected: ~30% shorter, explicit role mapping for all subagents

**Step 5: Commit**

```bash
git add skills/subagent-driven-development/SKILL.md
git commit -m "refactor: subagent-driven-development skill - concise with role specifications"
```

---

## Task 6: Update executing-plans/SKILL.md

**Files:**
- Modify: `skills/executing-plans/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/executing-plans/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **coder** (plan execution)

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep "The Process" steps (already minimal)
- Add role label to Step 5: "Complete Development" → calls architect (via finishing skill)
- Keep "When to Stop" section
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/executing-plans/SKILL.md
```

Expected: Minimal changes (already concise), add role labels

**Step 5: Commit**

```bash
git add skills/executing-plans/SKILL.md
git commit -m "refactor: executing-plans skill - add role specifications"
```

---

## Task 7: Update requesting-code-review/SKILL.md

**Files:**
- Modify: `skills/requesting-code-review/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/requesting-code-review/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **architect** (code review)

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep "When to Request Review" bullet list
- Keep "How to Request" steps
- Compress "Example"
- Keep "Integration with Workflows"
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/requesting-code-review/SKILL.md
```

Expected: Minor compression, role labels added

**Step 5: Commit**

```bash
git add skills/requesting-code-review/SKILL.md
git commit -m "refactor: requesting-code-review skill - concise with role specifications"
```

---

## Task 8: Update receiving-code-review/SKILL.md

**Files:**
- Modify: `skills/receiving-code-review/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/receiving-code-review/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **coder** (receives feedback, implements fixes)

Reviewer role: **architect** (provides review)

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep key process steps
- Add role labels: reviewer = architect, implementer = coder
- Compress examples
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/receiving-code-review/SKILL.md
```

Expected: ~25% shorter, role labels added

**Step 5: Commit**

```bash
git add skills/receiving-code-review/SKILL.md
git commit -m "refactor: receiving-code-review skill - concise with role specifications"
```

---

## Task 9: Update dispatching-parallel-agents/SKILL.md

**Files:**
- Modify: `skills/dispatching-parallel-agents/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/dispatching-parallel-agents/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **router** (coordinates parallel dispatch)

Agents dispatched use **coder** role.

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep "When to Use" flowchart
- Compress "The Pattern" steps
- Add role label: dispatched agents → coder
- Compress "Agent Prompt Structure"
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/dispatching-parallel-agents/SKILL.md
```

Expected: ~30% shorter, role labels added

**Step 5: Commit**

```bash
git add skills/dispatching-parallel-agents/SKILL.md
git commit -m "refactor: dispatching-parallel-agents skill - concise with role specifications"
```

---

## Task 10: Update writing-skills/SKILL.md

**Files:**
- Modify: `skills/writing-skills/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/writing-skills/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **architect** (skill design and documentation)

Testing may use **router** (for subagent coordination).

**Step 3: Simplify while preserving essential content**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (core TDD mapping)
- Keep "When to Create a Skill"
- Keep "SKILL.md Structure" but compress examples
- Keep "Claude Search Optimization" (CSO) - essential for discovery
- Compress "Flowchart Usage" section
- Keep "The Iron Law" (essential)
- Compress "Common Rationalizations" table
- Keep "Skill Creation Checklist" (essential)
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/writing-skills/SKILL.md
```

Expected: ~25% shorter, all essential content preserved

**Step 5: Commit**

```bash
git add skills/writing-skills/SKILL.md
git commit -m "refactor: writing-skills skill - concise with role specifications"
```

---

## Task 11: Update verification-before-completion/SKILL.md

**Files:**
- Modify: `skills/verification-before-completion/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/verification-before-completion/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **coder** (self-verification before claiming done)

**Step 3: Simplify while preserving gate function**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep "The Iron Law" (essential)
- Keep "The Gate Function" (essential)
- Keep "Common Failures" table
- Keep "Red Flags" list
- Compress "Rationalization Prevention" table
- Keep "Key Patterns"
- Keep "Why This Matters" (brief)
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/verification-before-completion/SKILL.md
```

Expected: ~20% shorter (already concise), minor compression

**Step 5: Commit**

```bash
git add skills/verification-before-completion/SKILL.md
git commit -m "refactor: verification-before-completion skill - concise with role specifications"
```

---

## Task 12: Update finishing-a-development-branch/SKILL.md

**Files:**
- Modify: `skills/finishing-a-development-branch/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/finishing-a-development-branch/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **router** (presents options, executes user choice)

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep "The Process" steps (already minimal)
- Keep "Quick Reference" table
- Compress "Common Mistakes"
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/finishing-a-development-branch/SKILL.md
```

Expected: Minor compression (already concise), add main role label

**Step 5: Commit**

```bash
git add skills/finishing-a-development-branch/SKILL.md
git commit -m "refactor: finishing-a-development-branch skill - add role specifications"
```

---

## Task 13: Update using-git-worktrees/SKILL.md

**Files:**
- Modify: `skills/using-git-worktrees/SKILL.md`

**Step 1: Read current skill**

```bash
cat skills/using-git-worktrees/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **router** (setup coordination)

**Step 3: Simplify and add role specifications**

Update skill structure:
- Keep frontmatter
- Keep "Overview" (already concise)
- Keep key setup steps
- Add role label: router orchestrates setup
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/using-git-worktrees/SKILL.md
```

Expected: Minimal changes (already concise)

**Step 5: Commit**

```bash
git add skills/using-git-worktrees/SKILL.md
git commit -m "refactor: using-git-worktrees skill - add role specifications"
```

---

## Task 14: Update using-superpowers/SKILL.md

**Files:**
- Modify: `skills/using-superpowers/SKILL.md

**Step 1: Read current skill**

```bash
cat skills/using-superpowers/SKILL.md
```

**Step 2: Identify main role and workflow steps**

Main role: **router** (entry point, skill selection)

**Step 3: Simplify while preserving essential guidance**

Update skill structure:
- Keep frontmatter
- Keep "How to Access Skills" flowchart
- Keep "Skill Priority" list
- Keep "The Rule" and red flags (essential)
- Compress examples
- Keep Conductor Integration Notes

**Step 4: Verify changes**

```bash
git diff skills/using-superpowers/SKILL.md
```

Expected: ~20% shorter, essential rules preserved

**Step 5: Commit**

```bash
git add skills/using-superpowers/SKILL.md
git commit -m "refactor: using-superpowers skill - concise with role specifications"
```

---

## Task 15: Create Role Mapping Summary Document

**Files:**
- Create: `docs/role-mapping-summary.md`

**Step 1: Create role mapping summary**

```bash
cat > docs/role-mapping-summary.md << 'EOF'
# Skills to Role Mapping Summary

This document maps each skill to its main role and specifies which roles are involved in workflow steps.

## Role Definitions

| Role | Description | Model | Primary Function |
|------|-------------|--------|------------------|
| router | Orchestrator & Router | Claude 3.5 Sonnet | Routes work, delegates tasks |
| architect | Architect & Planner | Gemini 1.5 Pro | Architecture, planning, review |
| coder | Deep Worker | Deepseek Coder | TDD implementation |
| librarian | Researcher | Gemini Flash | Research, documentation |

## Skills Role Mapping

| Skill | Main Role | Other Roles in Workflow |
|-------|-----------|------------------------|
| brainstorming | architect | - |
| writing-plans | architect | - |
| systematic-debugging | coder | librarian (research) |
| test-driven-development | coder | - |
| subagent-driven-development | router | coder (implement), architect (review) |
| executing-plans | coder | architect (via finishing skill) |
| requesting-code-review | architect | - |
| receiving-code-review | coder | architect (reviewer) |
| dispatching-parallel-agents | router | coder (agents) |
| writing-skills | architect | router (testing) |
| verification-before-completion | coder | - |
| finishing-a-development-branch | router | - |
| using-git-worktrees | router | - |
| using-superpowers | router | - |

## Workflow Patterns

### Pattern 1: Router-Led Workflows
Skills: subagent-driven-development, dispatching-parallel-agents, finishing-a-development-branch

Router coordinates and delegates:
- Implementation tasks → coder
- Review tasks → architect

### Pattern 2: Architect-Led Design
Skills: brainstorming, writing-plans, requesting-code-review

Architect produces outputs:
- Architecture specifications
- Implementation plans
- Code review feedback

### Pattern 3: Coder-Led Execution
Skills: test-driven-development, systematic-debugging, executing-plans

Coder executes with discipline:
- TDD workflow
- Root cause investigation
- Plan implementation
