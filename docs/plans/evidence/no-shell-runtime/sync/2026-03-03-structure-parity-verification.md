# Structure Parity Verification Evidence

Date: 2026-03-03 (UTC)
Branch: `go-sync-automation`

## 1) Validate structure parity (dry-run)

- Command: `./scripts/validate-claude-structure.sh -dry-run`
- Timestamp (UTC): `2026-03-03T05:49:54Z`
- Exit code: `0`

```text
$ ./scripts/validate-claude-structure.sh -dry-run
dry-run: structure parity ok (config/sync/claude-structure-rules.json)
```

## 2) End-to-end sync dry-run

- Command: `./scripts/sync-all.sh -dry-run`
- Timestamp (UTC): `2026-03-03T05:49:55Z`
- Exit code: `0`

```text
$ ./scripts/sync-all.sh -dry-run
dry-run: would run sync-upstream-main; would sync 6 COPY_FROM_MAIN paths from config/sync/main-to-go-rules.json
```

## 3) Targeted tests

- Command: `go test ./internal/devx ./cmd/mp-devx -v`
- Timestamp (UTC): `2026-03-03T05:49:55Z`
- Exit code: `0`

```text
$ go test ./internal/devx ./cmd/mp-devx -v
=== RUN   TestListTreeNames_UsesGitLsTree
--- PASS: TestListTreeNames_UsesGitLsTree (0.00s)
=== RUN   TestRunSyncMainToGo_CopiesOnlyCopyFromMainRules
--- PASS: TestRunSyncMainToGo_CopiesOnlyCopyFromMainRules (0.00s)
=== RUN   TestRunSyncUpstreamMain_PlansFetchAndFastForward
--- PASS: TestRunSyncUpstreamMain_PlansFetchAndFastForward (0.00s)
=== RUN   TestDevxRunner_SuiteUnitRunsGoTest
--- PASS: TestDevxRunner_SuiteUnitRunsGoTest (0.00s)
=== RUN   TestWithTempWorktree_UsesIsolatedPath
--- PASS: TestWithTempWorktree_UsesIsolatedPath (0.00s)
=== RUN   TestCompareStructure_MustHomomorphicDetectsMissingAndExtra
--- PASS: TestCompareStructure_MustHomomorphicDetectsMissingAndExtra (0.00s)
=== RUN   TestLoadStructureRules_ValidAndInvalid
--- PASS: TestLoadStructureRules_ValidAndInvalid (0.00s)
=== RUN   TestValidateStructureParity_PassesWithExplicitIgnores
--- PASS: TestValidateStructureParity_PassesWithExplicitIgnores (0.00s)
=== RUN   TestValidateStructureParity_FailsOnUnignoredDifferences
--- PASS: TestValidateStructureParity_FailsOnUnignoredDifferences (0.00s)
=== RUN   TestLoadSyncRules_ValidAndInvalid
=== RUN   TestLoadSyncRules_ValidAndInvalid/loads_valid_rules
=== RUN   TestLoadSyncRules_ValidAndInvalid/rejects_unknown_decision
--- PASS: TestLoadSyncRules_ValidAndInvalid (0.00s)
    --- PASS: TestLoadSyncRules_ValidAndInvalid/loads_valid_rules (0.00s)
    --- PASS: TestLoadSyncRules_ValidAndInvalid/rejects_unknown_decision (0.00s)
PASS
ok   github.com/gitbruce/claude-octopus/internal/devx  0.008s
=== RUN   TestRun_ActionSyncAll
--- PASS: TestRun_ActionSyncAll (0.00s)
=== RUN   TestRun_ActionValidateStructureParity
--- PASS: TestRun_ActionValidateStructureParity (0.39s)
PASS
ok   github.com/gitbruce/claude-octopus/cmd/mp-devx  0.394s
```

## 4) Architecture diff docs consistency

- Command: `scripts/verify-architecture-diff-docs.sh`
- Timestamp (UTC): `2026-03-03T05:49:56Z`
- Exit code: `0`

```text
$ scripts/verify-architecture-diff-docs.sh
Checking baseline hash consistency...
Checking evidence level legends...
Checking decision tokens...
Checking hook lifecycle events...
Checking mcp/openclaw decision tags...
Checking tracker structure...

verify-architecture-diff-docs: PASS
```

## 5) Freshness re-check

- Command: `./scripts/validate-claude-structure.sh -dry-run`
- Timestamp (UTC): `2026-03-03T05:49:56Z`
- Exit code: `0`

```text
$ ./scripts/validate-claude-structure.sh -dry-run
dry-run: structure parity ok (config/sync/claude-structure-rules.json)
```
