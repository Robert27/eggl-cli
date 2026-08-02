# Troubleshooting

Start with:

```bash
eggl doctor --verbose
```

## `kubectl` is missing

Install kubectl and verify it is on your `PATH`:

```bash
kubectl version --client
```

The `env` and `pf` commands need kubectl. Other commands do not.

## Tailscale is missing or the account is unavailable

Verify the CLI and account list:

```bash
tailscale version
tailscale switch --list
```

The `tailscale_account` value may be an ID, tailnet name, or email shown in that list.

## The active profile is unknown

`eggl env show` compares the current kubectl context and Tailscale account against every configured profile. Apply the intended profile explicitly:

```bash
eggl env use homelab
```

## `toggle` refuses to run

`toggle` requires exactly two configured profiles and a current state that matches one of them. Use `env use` when you have more than two profiles or the current state is not recognized.

## A config file is rejected

Check the path and parse it through an eggl CLI command:

```bash
eggl env path
eggl env show
```

Common causes are missing required fields, invalid Kubernetes resource syntax, invalid port mappings, or a file larger than 1 MiB.

## A port-forward cannot start

Confirm the active context, namespace, resource, and local port:

```bash
kubectl get svc -n monitoring grafana
```

If the local port is already occupied:

```bash
eggl kill --dry-run 3000
eggl kill --yes 3000
```

## A command changes nothing

For `dedash` and `eol`, check whether the file is excluded by extension, hidden-path rules, Git scope, binary detection, or the 50 MiB size limit. Use `--verbose` to see skipped paths.
