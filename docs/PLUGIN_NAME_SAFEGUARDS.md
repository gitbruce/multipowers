# Plugin Name Safeguards - Quick Reference

## The Core Constraint

Commands depend on a consistent plugin namespace.
- **Plugin Name**: `"mp"` ✅
- **Commands**: `/mp:discover`, `/mp:debate`, `/mp:status`, etc.

## Safeguards in Place

✅ **Plugin name is locked to `"mp"` with multiple layers of protection:**

### Layer 1: Manifest Protection
The `.claude-plugin/plugin.json` file contains explicit comments warning against changing the `name` field.

### Layer 2: Go-Native Validation
Validation is now integrated into the **`mp-devx`** toolchain.
```bash
# Verify plugin structure and name parity
./mp-devx -action parity
```

### Layer 3: CI/CD Integration
The GitHub Actions workflow runs the `mp-devx` parity check on every Pull Request to ensure the command namespace hasn't drifted.

## If It Breaks

1. **Verify `plugin.json`**:
   Ensure `"name": "mp"` is set in `.claude-plugin/plugin.json`.

2. **Run Parity Check**:
   ```bash
   ./mp-devx -action parity
   ```

---

**Status:** ✅ All safeguards active (Go-native parity)
**Last Verified:** March 2026
