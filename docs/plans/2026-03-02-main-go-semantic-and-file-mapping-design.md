# Main vs Go Semantic and File Mapping Design

Date: 2026-03-02  
Status: Approved

## 1. Context

Current branch-comparison docs are inconsistent and partially stale:
- `docs/architecture/commands_skills_difference.md` is based on older branch state and outdated version/count assumptions.
- `docs/architecture/script-differences.md` already tracks script migration for upstream `v8.31.1`, but needs cross-check against strict parity requirements.
- There is no dedicated document for non-command/non-skill/non-script file-type differences.

User requirement:
- Compare `main` and `go` carefully.
- Evaluate whether `go` reaches equivalent outcomes to `main`.
- Enforce semantic one-to-one mapping first.
- If semantic mapping is non-bijective, add file-level one-to-one supplement mappings.
- Verify and correct both existing architecture diff docs.
- Add additional file-type comparison results to `docs/architecture/other-differences.md`.

## 2. Goal

Produce a consistent, auditable parity analysis across three architecture documents:
1. Commands/skills mapping (`commands_skills_difference.md`)
2. Script mapping (`script-differences.md`)
3. All other file types (`other-differences.md`)

The output must clearly state parity status and precise gap/remediation paths.

## 3. Scope

### In Scope
- Full branch comparison: `main` vs `go`
- Semantic one-to-one mapping for same capability classes
- File-level one-to-one supplement for non-bijective mappings
- Correctness audit and correction of:
  - `docs/architecture/commands_skills_difference.md`
  - `docs/architecture/script-differences.md`
- Creation of:
  - `docs/architecture/other-differences.md`

### Out of Scope
- Runtime implementation changes
- Refactor of production logic
- Rewriting migration strategy beyond comparison and mapping correction

## 4. Design

### 4.1 Document Responsibilities

- `commands_skills_difference.md`
  - Canonical mapping for command/skill capabilities.
  - Semantic one-to-one mapping as primary table.
  - File-level one-to-one supplement for non-bijective cases.

- `script-differences.md`
  - Script inventory and migration strategy ledger.
  - Validates strategy/status counts and target ownership.
  - Cross-references command/skill parity where needed.

- `other-differences.md`
  - Non command/skill/script categories:
    - `go`, `json`, `yml/yaml`, `toml`, `ts/js`, docs, assets, configs, metadata, etc.
  - Semantic mapping first, file-level supplement second.

### 4.2 Mapping Rules

1. Semantic mapping is mandatory first.
2. If semantic mapping is not one-to-one by files:
   - add explicit file-level supplement rows.
3. Each mapping row must include status:
   - `equivalent`
   - `partial`
   - `missing`
   - `intentional-diff`
4. Each `partial/missing` row must include remediation target.
5. Final parity conclusion is based on semantic equivalence plus traceable file mapping, not path identity alone.

### 4.3 Output Model

Each doc should include:
- summary counts
- semantic mapping table
- file-level supplement table
- gap/remediation section
- parity conclusion

## 5. Execution Flow

1. Collect complete file inventories from `main` and `go`.
2. Build capability buckets by type (commands/skills/scripts/others).
3. Fill semantic mapping tables first.
4. Generate file-level supplement for all non-bijective semantic rows.
5. Cross-check count consistency across docs.
6. Publish final parity result with explicit unresolved gaps.

## 6. Validation Strategy

- Coverage validation:
  - every relevant capability in `main` appears in mapping output.
- Consistency validation:
  - no conflicting counts/status among the three docs.
- Mapping completeness:
  - no non-bijective semantic mapping without file-level supplement.
- Conclusion integrity:
  - parity verdict must be evidence-backed and reproducible.

## 7. Acceptance Criteria

- `commands_skills_difference.md` updated to current branch reality.
- `script-differences.md` verified/corrected with consistent totals and statuses.
- `other-differences.md` created and populated with non-script categories.
- Semantic one-to-one mapping available for each analyzed class.
- File one-to-one supplements present for all non-bijective cases.
- A clear answer is provided on whether `go` reaches `main`-equivalent outcomes.

## 8. Risks and Mitigations

- Risk: path refactors mistaken as missing capability.
  - Mitigation: semantic-first mapping plus supplement tables.
- Risk: stale doc baselines.
  - Mitigation: regenerate counts from current branch trees.
- Risk: oversized docs become unreadable.
  - Mitigation: compact summary + structured tables + explicit cross-references.

## 9. Next Step

Invoke `writing-plans` to produce an implementation plan for executing this comparison and doc update work.
