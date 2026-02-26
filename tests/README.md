# Claude Octopus Test Suite

This directory contains automated tests for the Claude Octopus plugin.

## Running Tests

**Run all tests:**
```bash
go test ./...
```

**Run individual tests:**
```bash
go test ./...
go test ./...
go test ./...
# etc.
```

## Test Descriptions

### Core Functionality Tests

- **`go test suite`** - Validates enforcement pattern documentation structure
  - ⚠️ **Important:** Tests documentation only, NOT runtime enforcement
  - Verifies all mp runtime skills have consistent Validation Gate Pattern docs
  - Does NOT verify mp runtime is actually executed at runtime
  - See: [Issue #TBD](https://github.com/anthropics/claude-code/issues/) for runtime enforcement tracking

- **`go test suite`** - Ensures version numbers match across all files
  - Checks: plugin.json, marketplace.json, package.json, README.md

- **`go test suite`** - Validates slash commands are properly registered
  - Checks: plugin.json commands match .claude-plugin/.claude/commands/*.md files

### Feature-Specific Tests

- **`test-intent-contract-skill.sh`** - Validates intent contract skill structure
- **`test-intent-questions.sh`** - Tests interactive question functionality
- **`test-plan-command.sh`** - Validates /plan command integration
- **`test-multi-command.sh`** - Tests multi-command workflows
- **`test-v2.1.12-integration.sh`** - Validates Claude Code v2.1.12+ feature integration

### Validation Tests

- **`go test ./...`** - Ensures plugin name consistency

## Important Limitations

### Enforcement Pattern Tests (v7.15.0)

As of v7.15.0, the enforcement pattern is **documentation-only**. The `go test suite` suite verifies that:

✅ **What IS tested:**
- Skills have `execution_mode: enforced` in frontmatter
- Skills contain EXECUTION CONTRACT sections with numbered steps
- Skills use imperative language ("MUST", "PROHIBITED", "CANNOT SKIP")
- Skills document validation gates for artifact checking
- Skills have consistent multi-AI attribution

❌ **What is NOT tested:**
- Whether mp runtime is actually executed when skill is invoked
- Whether AskUserQuestion is called before proceeding
- Whether validation gates are checked at runtime
- Whether Claude follows the EXECUTION CONTRACT steps

**Why:** Claude Code does not currently support skill lifecycle hooks. Skills are passive markdown documentation that Claude interprets as guidance, not enforceable requirements.

**Tracking:** See `scratchpad/github-issue-skill-lifecycle-hooks.md` for the feature request to Anthropic for programmatic enforcement.

**Implication:** The enforcement pattern (v7.15.0) provides consistent documentation structure but does not guarantee runtime behavior. Users should invoke mp runtime manually if needed until lifecycle hooks are implemented.

## Test Coverage

| Test Suite | Tests | Coverage |
|------------|-------|----------|
| Enforcement Pattern | 20 | Documentation structure, frontmatter, execution contracts |
| Version Consistency | 4 | All version references across files |
| Command Registration | Variable | All slash commands |
| Intent Contract | Variable | Intent contract skill structure |
| v2.1.12 Integration | Variable | Claude Code v2.1.12+ features |

## Adding New Tests

When adding new tests:

1. **Create test file:** `go test ./...`
2. **Make executable:** `chmod +x go test ./...`
3. **Follow conventions:**
   - Use `set -euo pipefail` for safety
   - Use colored output (RED, GREEN, YELLOW, BLUE, NC)
   - Count passes/fails with `pass()` and `fail()` functions
   - Exit 0 on success, 1 on failure

4. **Add to run-all-tests.sh:**
   ```bash
   run_test "Your Feature Description" "./go test ./..."
   ```

5. **Document in this README** with clear description of what is and isn't tested

## Test Philosophy

**Focus on verifiable behavior:**
- ✅ File structure and content presence
- ✅ Version number consistency
- ✅ Command registration
- ✅ Documentation completeness

**Avoid testing runtime AI behavior:**
- ❌ Whether Claude follows skill instructions (non-deterministic)
- ❌ Quality of AI responses (subjective)
- ❌ Multi-AI orchestration results (dependent on external CLIs)

**Exception:** Integration tests can verify external CLI execution (codex, gemini) when invoked directly, but cannot verify Claude's decision to invoke them.

## CI/CD Integration

These tests are designed to run in CI/CD pipelines:

```bash
# Run all tests, exit non-zero on any failure
go test ./...

# Individual test exit codes
go test ./... && echo "Pass" || echo "Fail"
```

## Known Issues

1. **Enforcement pattern runtime gap** - Documentation exists but not enforced (v7.15.0)
   - Tests verify docs, not behavior
   - Tracking: GitHub issue (pending submission)

2. **Version test fragility** - May fail if versions updated without running tests
   - Always run `go test ./...` after version bumps

## Future Tests (Pending Claude Code Features)

- **Runtime enforcement tests** - Verify mp runtime execution (requires lifecycle hooks)
- **Interactive question tests** - Verify AskUserQuestion is called (requires hooks)
- **Validation gate tests** - Verify artifact checks happen (requires hooks)
- **Fork context tests** - Verify memory-optimized skill execution

---

**Questions or issues?** See main README or file an issue at https://github.com/nyldn/claude-octopus/issues
