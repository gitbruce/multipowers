# Codex Sandbox Configuration

This guide explains how to configure Codex sandbox mode for advanced use cases like mounted filesystems.

## Overview

By default, Codex agents run in `workspace-write` sandbox mode, which restricts filesystem access to the current workspace.

**Added in:** v7.13.1 (addressing Issue #9)

## When You Need This

You may need to configure sandbox mode if:
- Working with repositories on mounted filesystems (SSHFS, NFS, FUSE)
- Running code audits on remote repositories
- Getting `Sandbox(LandlockRestrict)` errors from Codex

## Configuration

### Environment Variable (Recommended)

Set the **`MP_CODEX_SANDBOX`** environment variable:

```bash
# Temporary (current session only)
export MP_CODEX_SANDBOX=danger-full-access
mp discover "audit code in mounted repo"

# Permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export MP_CODEX_SANDBOX=danger-full-access' >> ~/.bashrc
```

### Per-Command Override

```bash
# One-time override
MP_CODEX_SANDBOX=danger-full-access mp discover "audit /mnt/nas/repo"
```

## Security Considerations

### ⚠️ Risks of `danger-full-access`

- **Full filesystem access**: Codex can read any file the user can read.
- **Data exfiltration risk**: Malicious prompts could leak sensitive data.

### ✅ Mitigation Strategies

1. **Use for read-only tasks only**: Code audits, research, analysis.
2. **Temporary override**: Use per-command override instead of permanent export.
3. **Output isolation**: Multipowers only writes to `~/.multipowers/results/` by design.

---

**Last Updated:** March 2026
