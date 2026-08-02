# Environment Profiles

An environment profile pairs a Kubernetes context with a Tailscale account:

```yaml
profiles:
  homelab:
    kube_context: homelab-cluster
    tailscale_account: homelab.example.com
  work:
    kube_context: work-cluster
    tailscale_account: work@example.com
```

The account value can be an account ID, tailnet name, or email from:

```bash
tailscale switch --list
```

## Show state

```bash
eggl env show
```

The report includes the active profile, current kubectl context, current Tailscale account, config path, and all configured profiles.

If the current context/account pair does not match a profile, the active profile is shown as `unknown`.

## Apply a profile

```bash
eggl env use homelab
```

Only the values that differ are changed. If the Tailscale switch fails after the Kubernetes context changed, the error reports that partial state so it can be corrected explicitly.

## Toggle two profiles

```bash
eggl env toggle
```

`toggle` requires exactly two profiles and requires the current Kubernetes context and Tailscale account to match one of them. If the current state is unknown, use `eggl env use <profile>` first.
