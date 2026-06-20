# Getting Started

**eggl** is a personal helper CLI for everyday development tasks — file hygiene, environment switching, Kubernetes port-forwards, and more.

## Quick start

```bash
# Install (Homebrew)
brew tap Robert27/tap
brew install eggl-cli

# Verify
eggl version
eggl doctor
```

## What can eggl do?

| Area | Commands | Description |
|------|----------|-------------|
| File hygiene | `dedash`, `eol` | Replace em-dashes and normalize line endings |
| Git | `empty` | Create empty commits to retrigger CI |
| Cloud | `env`, `pf` | Switch kube + Tailscale profiles; run port-forwards |
| System | `kill`, `doctor` | Free stuck ports; check your environment |

## Global flags

Every command inherits these root flags:

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Log operation details to stderr (skipped paths, progress, etc.) |

## Next steps

- [Installation](/guide/installation) — Homebrew, releases, and building from source
- [Configuration](/guide/configuration) — Set up `env` profiles and `pf` port-forwards
- [All commands](/commands/) — Full command reference
