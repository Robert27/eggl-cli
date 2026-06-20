# version

Print build version, git commit, and build date.

## Usage

```bash
eggl version [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--short` | `false` | Print only the version string |

## Examples

```bash
eggl version
# version v0.9.1
# commit  abc1234
# built   2024-01-15T12:00:00Z

eggl version --short
# v0.9.1
```

## Notes

- Release builds embed version from git tags via GoReleaser
- Local `make build` defaults to `dev` unless `VERSION` is set
