# ⚠️ PLUGIN NAME LOCK

## CRITICAL: DO NOT CHANGE THE PLUGIN NAME

The plugin name in `plugin.json` **MUST remain "multipowers"**.

### Why?

```json
// ✅ CORRECT - plugin.json
{
  "name": "multipowers"  // This produces /mp:discover, /mp:debate, etc.
}
```

```json
// ❌ WRONG - DO NOT DO THIS
{
  "name": "claude-octopus"  // This produces /claude-octopus:discover (too long!)
}
```

### Package vs Plugin Name

These are **different** and serve **different purposes**:

| File | Name | Purpose |
|------|------|---------|
| `package.json` | `"claude-octopus"` | Marketplace/repository identity |
| `.claude-plugin/plugin.json` | `"multipowers"` | Command prefix (`/mp:*`) |

### Command Path Formation

Command paths are formed as: `/[plugin-name]:[command-name]`

- Plugin name: `"multipowers"` + Command: `discover` = `/mp:discover` ✅
- Plugin name: `"claude-octopus"` + Command: `discover` = `/claude-octopus:discover` ❌

### Historical Context

**Commits that fixed this:**
- `d9e8354` - Reverted plugin name to 'multipowers' for correct command prefixes
- `57ce38c` - Removed namespace prefix from command frontmatter

**Why it broke:**
Someone changed the plugin name thinking it should match the package name. It shouldn't.

### Tests

Run `make test-plugin-name` to verify the plugin name is correct.

### If You Need to Change It

**Don't.** But if you absolutely must:
1. Update all documentation showing `/mp:*` commands
2. Update README.md examples
3. Update all skill files with command references
4. Notify all users about the breaking change
5. Consider providing migration script
6. Update this documentation

**Estimated impact:** 100+ command references across docs, skills, and user workflows.

---

**Last verified:** 2026-01-21
**Status:** ✅ Plugin name is "multipowers" and LOCKED
