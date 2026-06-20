# dedash

Replace Unicode em-dashes (`—`) with ASCII hyphens (`-`) in text files.

## Usage

```bash
eggl dedash [flags]
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

- Recursively scans the given directory
- Skips binaries, `node_modules`, `.git`, `vendor`, and similar directories
- Skips dotfiles and dot-directories unless `--include-hidden`
- Skips files larger than 50 MiB
- Prints a summary plus per-file replacement counts

## Examples

```bash
# Preview changes in current directory
eggl dedash --dry-run

# Process docs folder without prompting
eggl dedash --path ./docs --yes

# Only markdown and text files
eggl dedash --ext md,txt --dry-run

# Only files with uncommitted git changes
eggl dedash --diff --dry-run

# Files changed since main
eggl dedash --diff-base main --dry-run
```

## Notes

- Prompts for confirmation before modifying files unless `--dry-run` or `--yes`
- Use `--verbose` to see skipped paths and each modified file on stderr
