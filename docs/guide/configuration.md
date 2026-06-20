# Configuration

eggl stores configuration at:

```
~/.config/eggl/config.yaml
```

Or `$XDG_CONFIG_HOME/eggl/config.yaml` when `XDG_CONFIG_HOME` is set.

Create the file with:

```bash
eggl env init
```

## Environment profiles

Profiles pair a kubectl context with a Tailscale account so you can switch both at once.

```yaml
profiles:
  alpha:
    kube_context: context-a
    tailscale_account: b3e1
  beta:
    kube_context: context-b
    tailscale_account: a7f2
```

`tailscale_account` can be:

- An account ID from `tailscale switch --list`
- A tailnet slug
- An email address

### Commands

```bash
eggl env show      # Active profile and current state
eggl env use alpha # Switch to a named profile
eggl env toggle    # Flip between two profiles (requires exactly 2)
eggl env path      # Print config file path
```

Override the config path with `--config /path/to/config.yaml` on any `env` subcommand.

## Port-forwards

Named kubectl port-forwards live in the same config file:

```yaml
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
    ports: ["8080:80"]
```

| Field | Required | Description |
|-------|----------|-------------|
| `namespace` | yes | Kubernetes namespace |
| `resource` | yes | Resource spec, e.g. `svc/name`, `deploy/name` |
| `ports` | no | Port mappings (`local:remote`). Defaults to `8080:80` |

```bash
eggl pf list       # List configured forwards
eggl pf longhorn   # Start forwarding (blocks until Ctrl+C)
```

Port-forwards use the **active** kubectl context — they are not tied to env profiles directly.

## Validation rules

- At least one profile is required
- Each profile needs `kube_context` and `tailscale_account`
- No duplicate profile targets (same kube context + Tailscale account pair)
- Port-forward namespaces and resources must be valid Kubernetes identifiers
