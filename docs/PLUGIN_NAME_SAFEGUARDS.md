# Plugin Name Safeguards - Quick Reference

## What Was Broken

Commands were breaking because the plugin name kept getting changed:
- Changed: `"multipowers"` → `"multipowers"` ❌
- Commands became: `/multipowers:discover` (too long, broke workflows)
- Should be: `/mp:discover` ✅

## What's Protected Now

✅ **Plugin name locked to `"multipowers"` with 4 layers of protection:**

### Layer 1: Documentation Warnings
```
.claude-plugin/plugin.json        ← In-file comment
.claude-plugin/PLUGIN_NAME_LOCK.md ← Detailed explanation
.claude-plugin/README.md           ← Quick warning
SAFEGUARDS.md                      ← Comprehensive reference
```

### Layer 2: Automated Tests
```bash
make test-plugin-name              # Runs validation
./tests/validate-plugin-name.sh    # Direct validation
```

### Layer 3: CI/CD Integration
- ✅ GitHub Actions validates on every PR
- ✅ Smoke tests include plugin name validation
- ✅ Pre-commit hook validates before commit

### Layer 4: Make Target Integration
```makefile
test-smoke: test-plugin-name       # Smoke tests depend on validation
```

## Quick Validation

Run this to verify everything is correct:

```bash
make test-plugin-name
```

Expected output:
```
🔍 Validating plugin name...
✅ Plugin name is correct: "multipowers"
   Commands will work as: /mp:discover, /mp:debate, etc.
```

## If It Breaks Again

1. Check the plugin name:
   ```bash
   grep '"name"' .claude-plugin/plugin.json
   # Should show: "name": "mp"
   ```

2. If wrong, fix it immediately:
   ```json
   {
     "name": "mp"  // ← Must be exactly this
   }
   ```

3. Run validation:
   ```bash
   make test-plugin-name
   ```

## Why Plugin Name ≠ Package Name

| Purpose | File | Name |
|---------|------|------|
| Command prefix | `.claude-plugin/plugin.json` | `"multipowers"` |
| Marketplace ID | `package.json` | `"multipowers"` |

Both are correct and serve different purposes.

---

**Status:** ✅ All safeguards active
**Last Verified:** 2026-01-21
**Commands Working:** `/mp:discover`, `/mp:debate`, `/mp:embrace`, etc.
