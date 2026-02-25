# Upstream Sync Playbook

## Routine Sync Sequence

```bash
git fetch upstream origin
git switch main
git merge --ff-only upstream/main
git switch multipowers
git merge main -m "chore(sync): merge main into multipowers"
./custom/scripts/apply-custom-overlay.sh
bash tests/integration/test-sync-overlay.sh
```

## Conflict SLA and Fallback

- Target: resolve sync conflicts within 30 minutes.
- If not resolved within SLA:

```bash
git rebase --abort
# or git merge --abort, depending on operation
```

Then restart from clean baseline using this playbook.

## Example Sync Transcript

See: `custom/docs/sync/verification-transcript.md`

## Expected Result

- `main` matches `upstream/main`
- overlay reapplied successfully
- sync/registration tests pass
