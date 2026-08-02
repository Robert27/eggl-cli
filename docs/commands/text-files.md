# `eggl dedash` and `eggl eol`

Both commands recursively scan text files and provide the same safety controls.

| Command | Change |
| --- | --- |
| `eggl dedash` | Replace Unicode em-dashes with ASCII hyphens |
| `eggl eol` | Convert CRLF or CR line endings to LF |

## Preview first

```bash
eggl dedash --dry-run
eggl eol --dry-run
```

Without `--dry-run`, the command asks for confirmation before writing. Use `--yes` in automation or when stdin is not a terminal.

## Scope controls

```bash
eggl dedash --path ./docs --ext md,txt --dry-run
eggl eol --diff --dry-run
eggl dedash --diff-base main --dry-run
```

| Flag | Description |
| --- | --- |
| `--path <dir>` | Directory tree to scan; defaults to `.` |
| `--dry-run` | Report changes without writing |
| `-y, --yes` | Skip the write confirmation |
| `--ext <list>` | Limit processing to extensions such as `md,txt` |
| `--include-hidden` | Include dotfiles and dot-directories |
| `--diff` | Process staged or unstaged Git changes |
| `--diff-base <ref>` | Process files changed since a Git ref |

`--diff` and `--diff-base` are mutually exclusive. Binaries, common dependency directories, `.git`, and files larger than 50 MiB are skipped.
