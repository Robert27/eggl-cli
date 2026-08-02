# Diagnostics and Metadata

## `eggl doctor`

Run sanity checks for the local environment:

```bash
eggl doctor
eggl doctor --verbose
eggl doctor --check-path /tmp
```

The checks cover the Go runtime, platform, home directory, `kubectl`, Git, Tailscale, and config validity when a config file exists. The command exits non-zero if any check fails.

## `eggl version`

```bash
eggl version
eggl version --short
```

The full output includes the installed version, commit, and build date. `--short` prints only the version string for scripts.

## `eggl completion`

Generate Bash, Zsh, Fish, or PowerShell completion scripts:

```bash
eggl completion zsh
eggl completion bash
```

See [Shell completion](/completion) for installation instructions.
