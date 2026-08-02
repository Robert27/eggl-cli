# `eggl kill`

Find processes listening on a local TCP port and terminate them.

```bash
eggl kill 8080
eggl kill --dry-run 3000
eggl kill --yes --force 8080
```

## Flags

| Flag | Description |
| --- | --- |
| `--dry-run` | List matching processes without killing them |
| `-y, --yes` | Skip the confirmation prompt |
| `-f, --force` | Send `SIGKILL` instead of `SIGTERM` |

Use `--dry-run` before `--force` when investigating a stale port-forward or development server. Without `--force`, eggl CLI requests graceful termination.
