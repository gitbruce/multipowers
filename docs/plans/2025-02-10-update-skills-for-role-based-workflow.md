# Skills Update for Role-Based Workflow Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make all skills in `skills/*/SKILL.md` concise while preserving helpful guidance, and add explicit role contracts: one main role per skill plus step-level role overrides (when needed) using only roles from `config/roles.default.json`.

**Architecture:** Use a shared edit contract + per-skill role matrix. Keep workflow logic intact, remove redundancy, and annotate workflow steps with `[Role: ...]` labels where non-main specialists are required.

**Tech Stack:** Markdown (`SKILL.md`), Bash (`rg`, `git diff`), JSON role source (`config/roles.default.json`)

---

## Verification Result (2026-02-10)

The previous version of this plan **did not fully satisfy** requirements:

1. **Not concise enough:** repetitive per-skill boilerplate and per-task commit blocks made the plan hard to scan.
2. **Role-source drift risk:** role usage was not consistently tied to `config/roles.default.json`.
3. **Step-role detail inconsistent:** some skills had a main role but no explicit step-level overrides.
4. **Document integrity issue:** malformed/incomplete Markdown near the end reduced reliability.

This revision resolves those gaps.

---

## Source of Truth

### Allowed Roles (only these names)

From `config/roles.default.json`:

- `router`
- `architect`
- `coder`
- `librarian`
- `reviewer-claude`

### Context Rules (from `conductor/context/*.md`)

- **Workflow first, roles second:** pick workflow lane first, then assign node roles.
- **Fast lane vs standard lane:** small changes can route directly; major/complex work uses explicit workflow nodes.
- **Evidence before claims:** verification is required before completion statements.
- **Role contracts:** router orchestrates; architect plans/reviews; coder implements with TDD; librarian researches; reviewer-claude performs strict audit review.

### Required Skill Format Addition

Add a concise role block to every updated skill:

```markdown
## Role Contract
- Main Role: <router|architect|coder|librarian|reviewer-claude>
- Workflow Step Roles:
  1. <step summary> [Role: <allowed role>]
  2. <step summary> [Role: <allowed role>]
```

Rules:
- Exactly one `Main Role` per skill.
- Only add non-main roles where the workflow truly requires specialization.
- Do not invent role names outside `config/roles.default.json`.

---

## Per-Skill Role Matrix (Required)

| Skill | Main Role | Workflow Step Roles (include overrides) | Keep Helpful Info (must preserve) |
|---|---|---|---|
| `skills/brainstorming/SKILL.md` | `architect` | Intake and routing `[Role: router]` → clarify and design `[Role: architect]` → implementation handoff `[Role: router]` | One-question-at-a-time discipline, alternatives with trade-offs, incremental validation, design-doc handoff |
| `skills/writing-plans/SKILL.md` | `architect` | Plan authoring `[Role: architect]` → execution-mode handoff `[Role: router]` → execution ownership note `[Role: coder]` | Required plan header, bite-sized tasks, explicit commands, execution handoff options |
| `skills/systematic-debugging/SKILL.md` | `coder` | Root-cause + hypothesis workflow `[Role: coder]` → external evidence lookup `[Role: librarian]` (if needed) → architecture challenge after repeated failures `[Role: architect]` | Iron Law, 4 phases, data-flow tracing, “3+ failed fixes => architecture discussion” |
| `skills/test-driven-development/SKILL.md` | `coder` | RED-GREEN-REFACTOR `[Role: coder]` → spec/compliance check `[Role: architect]` (when required) | Iron Law, delete-code-before-test rule, fail-first verification, checklist quality gates |
| `skills/subagent-driven-development/SKILL.md` | `router` | Orchestration `[Role: router]` → implementer nodes `[Role: coder]` → spec review `[Role: architect]` → quality review `[Role: architect]` → optional final strict audit `[Role: reviewer-claude]` | Two-stage review order, re-review loops, per-task isolation, required workflow integrations |
| `skills/executing-plans/SKILL.md` | `coder` | Batch execution `[Role: coder]` → checkpoint reviews `[Role: architect]` → branch-finish decision flow `[Role: router]` | Critical-plan review first, batch/checkpoint cadence, stop-and-ask conditions |
| `skills/requesting-code-review/SKILL.md` | `architect` | Review request + criteria `[Role: architect]` → high-signal audit option `[Role: reviewer-claude]` → fix implementation loop `[Role: coder]` | Mandatory review triggers, review input template, issue-severity handling, reasoned pushback |
| `skills/receiving-code-review/SKILL.md` | `coder` | Feedback triage + implementation `[Role: coder]` → standard reviewer source `[Role: architect]` → strict audit reviewer (optional) `[Role: reviewer-claude]` → evidence lookup `[Role: librarian]` (if needed) | Verify-before-implement rule, clarify unclear items first, non-performative technical responses |
| `skills/dispatching-parallel-agents/SKILL.md` | `router` | Domain split + dispatch `[Role: router]` → parallel fixes `[Role: coder]` → integration/conflict review `[Role: architect]` | Independence criteria, one-agent-per-domain prompt discipline, post-merge verification |
| `skills/writing-skills/SKILL.md` | `architect` | Skill design + structure `[Role: architect]` → pressure-scenario execution `[Role: coder]` → orchestration/selection `[Role: router]` | TDD-for-skills framing, CSO rules, checklist and anti-rationalization guidance |
| `skills/verification-before-completion/SKILL.md` | `coder` | Verification gate execution `[Role: coder]` → spec compliance confirmation `[Role: architect]` (complex changes) | Evidence-before-claims gate, command-output proof requirement, red flags/rationalizations |
| `skills/finishing-a-development-branch/SKILL.md` | `router` | Option presentation + branch decision `[Role: router]` → fix loop if tests fail `[Role: coder]` → standard pre-merge audit (optional) `[Role: architect]` → strict pre-merge audit (optional) `[Role: reviewer-claude]` | Tests-first rule, exactly four completion options, destructive-action confirmation |
| `skills/using-git-worktrees/SKILL.md` | `router` | Worktree selection/orchestration `[Role: router]` → setup/test baseline fixes `[Role: coder]` | Directory priority, ignore-check requirement, baseline verification before implementation |
| `skills/using-superpowers/SKILL.md` | `router` | Skill applicability triage `[Role: router]` → route to creative design `[Role: architect]`, implementation `[Role: coder]`, or research `[Role: librarian]` | “Use skill first” rule, skill-priority model, rationalization red flags |

---

## Execution Tasks

### Task 1: Apply shared concise structure to all 14 skills

**Files:**
- Modify: all `skills/*/SKILL.md` files listed in the matrix above

**Steps:**
1. Keep core sections that define behavior and quality gates.
2. Remove duplicated prose/examples that repeat existing rules.
3. Insert `## Role Contract` with one main role and explicit step-role labels.
4. Ensure any non-main step clearly shows `[Role: ...]` using allowed names only.

### Task 2: Update router-led skills

**Files:**
- `skills/subagent-driven-development/SKILL.md`
- `skills/dispatching-parallel-agents/SKILL.md`
- `skills/finishing-a-development-branch/SKILL.md`
- `skills/using-git-worktrees/SKILL.md`
- `skills/using-superpowers/SKILL.md`

**Steps:**
1. Preserve orchestration logic and decision gates.
2. Label delegated implementation/review nodes with specialist roles.
3. Keep workflow order constraints explicit (especially review sequence).

### Task 3: Update architect-led skills

**Files:**
- `skills/brainstorming/SKILL.md`
- `skills/writing-plans/SKILL.md`
- `skills/requesting-code-review/SKILL.md`
- `skills/writing-skills/SKILL.md`

**Steps:**
1. Preserve planning/review methodology and templates.
2. Add role overrides for routing, execution, and strict audit nodes where needed.
3. Compress long narrative examples to compact patterns/checklists.

### Task 4: Update coder-led skills

**Files:**
- `skills/systematic-debugging/SKILL.md`
- `skills/test-driven-development/SKILL.md`
- `skills/executing-plans/SKILL.md`
- `skills/receiving-code-review/SKILL.md`
- `skills/verification-before-completion/SKILL.md`

**Steps:**
1. Preserve discipline rules (Iron Law / gate rules / required verification).
2. Add role overrides for architecture escalation and research support.
3. Keep regression-prevention and verification checkpoints explicit.

### Task 5: Validate role compliance and concision

**Steps:**
1. Verify role labels were added:
   ```bash
   rg -n "Main Role:|\\[Role:" skills/*/SKILL.md
   ```
2. Verify no out-of-config roles are introduced:
   ```bash
   rg -n "\\[Role: [^\\]]+\\]|Main Role: " skills/*/SKILL.md
   ```
   Manually confirm every role token is one of:
   `router`, `architect`, `coder`, `librarian`, `reviewer-claude`.
3. Verify useful content is retained (not over-pruned):
   ```bash
   git diff -- skills/*/SKILL.md
   ```
4. Verify major-change governance expectations still hold for final execution:
   - run required checks appropriate to changed files
   - update related docs if workflow behavior changed

---

## Completion Checklist

- [ ] Each skill has exactly one `Main Role`.
- [ ] Each workflow step that needs specialization has explicit `[Role: ...]`.
- [ ] Every role name comes from `config/roles.default.json`.
- [ ] Content is shorter but still preserves key guidance/checklists/guards.
- [ ] No workflow-order regressions introduced.
- [ ] Verification evidence captured before completion claims.
