# kill

Find and terminate processes listening on a local TCP port.

## Usage

```bash
eggl kill <port> [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `port` | TCP port number (1–65535) |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dry-run` | | `false` | List matching processes without killing |
| `--yes` | `-y` | `false` | Skip confirmation prompt |
| `--force` | `-f` | `false` | Send SIGKILL instead of SIGTERM |

## Behavior

- Finds processes listening on the specified port
- Sends SIGTERM by default; use `--force` for SIGKILL
- Works on Linux, macOS, and Windows (platform-specific process discovery)

## Examples

```bash
# Kill process on port 8080 (with confirmation)
eggl kill 8080

# Preview what would be killed
eggl kill --dry-run 3000

# Non-interactive kill with force
eggl kill --yes --force 8080
```

## Output

```
would kill pid 12345 (node)    # --dry-run
killed pid 12345               # success
no process listening on port 8080
```

## Notes

- Prompts before killing unless `--dry-run` or `--yes`
- In CI or piped environments, `--yes` is required
