# eol

Normalize line endings to LF by converting CRLF or CR to LF.

## Usage

```bash
eggl eol [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--path` | | `.` | Directory tree to scan |
| `--dry-run` | | `false` | Preview changes without writing |
| `--yes` | `-y` | `false` | Skip confirmation prompt |
| `--ext` | | (all text) | Limit to extensions, e.g. `md,txt` |
| `--include-hidden` | | `false` | Process dotfiles and dot-directories |
| `--diff` | | `false` | Only files with staged/unstaged git changes |
| `--diff-base` | | | Only files changed on current branch since ref |

`--diff` and `--diff-base` are mutually exclusive.

## Behavior

Same scan and skip logic as [`dedash`](/commands/dedash):

- Skips binaries, large files (> 50 MiB), and common non-text directories
- Skips dotfiles unless `--include-hidden`
- Reports files that were or would be modified

## Examples

```bash
eggl eol --dry-run
eggl eol --path ./src --yes
eggl eol --ext go,md --dry-run
eggl eol --diff --dry-run
eggl eol --diff-base main --yes
```

## Notes

- Prompts for confirmation before modifying files unless `--dry-run` or `--yes`
- Useful before committing on Windows or after merging branches with mixed line endings
