# doctor

Sanity-check your local environment and the tools eggl depends on.

## Usage

```bash
eggl doctor [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--check-path` | `$HOME` | Directory to validate instead of home |

## Checks performed

| Check | Description |
|-------|-------------|
| `go` | Go runtime version |
| `os` | `GOOS/GOARCH` |
| `home` | Home directory exists and is readable |
| `kubectl` | Available on `PATH` |
| `git` | Available on `PATH` |
| `tailscale` | Available on `PATH` |
| `config` | `~/.config/eggl/config.yaml` is valid if present |

A missing config file is OK — the check suggests running `eggl env init`.

## Examples

```bash
eggl doctor
eggl doctor --verbose
eggl doctor --check-path /tmp
```

## Exit code

Returns non-zero if any check fails. Missing `kubectl` or `tailscale` fails the check but only matters if you use `env` or `pf` commands.

## Output

Styled with checkmarks (`✓`) and crosses (`✗`) when stdout is a TTY. Plain text otherwise.
