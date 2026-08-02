# Workflows

These recipes combine eggl with the tools it already uses.

## Switch into a cluster environment

```bash
eggl env show
```

The profile makes the kubectl context and Tailscale account agree before you work with private cluster services.

## Open a local service

```bash
eggl env use homelab
eggl pf dashboard --open
```

The named port-forward is repeatable and keeps the resource definition in versionable YAML rather than in shell history.

## Clean changed repository files

Preview only the files changed against `main`:

```bash
eggl dedash --diff-base main --dry-run
eggl eol --diff-base main --dry-run
```

Apply the same changes after reviewing the preview:

```bash
eggl dedash --diff-base main --yes
eggl eol --diff-base main --yes
```

## Recover a stale development port

First inspect matching processes:

```bash
eggl kill --dry-run 8080
```

Then terminate them gracefully:

```bash
eggl kill --yes 8080
```

Use `--force` only when a process does not respond to a normal termination signal.

## Retrigger CI

```bash
eggl empty --push
```

This creates and pushes an empty commit in the current Git repository.
