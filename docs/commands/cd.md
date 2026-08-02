# `eggl cd`

Resolve a configured directory alias and print its path.

```bash
cd "$(eggl cd homelab)"
```

## Flags

| Flag | Description |
| --- | --- |
| `--config <path>` | Use a specific config file |

The `eggl cd` command does not change the caller's working directory because child processes cannot modify the parent shell. Use command substitution or the `egglcd` shell function shown in [Directory aliases](/configuration/directories).
