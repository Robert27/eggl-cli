# pf

Run named `kubectl port-forward` entries from your config file.

## Usage

```bash
eggl pf <name> [flags]
eggl pf list [flags]
```

## Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `~/.config/eggl/config.yaml` | Path to config file |

## Subcommands

### `list`

Print configured port-forwards (tab-separated: name, namespace/resource, ports).

```bash
eggl pf list
```

Example output:

```
longhorn    longhorn-system/svc/longhorn-frontend    8080:80
```

### `<name>`

Start a port-forward. Blocks until you press Ctrl+C.

```bash
eggl pf longhorn
```

Running `eggl pf` without arguments shows help.

## Configuration

Port-forwards are defined in the same config file as env profiles:

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
| `resource` | yes | e.g. `svc/name`, `deploy/name` |
| `ports` | no | `local:remote` mappings. Default: `8080:80` |

## Behavior

- Uses the **active** kubectl context (not tied to env profiles)
- Shell completion suggests configured port-forward names
- Requires `kubectl` on `PATH`

## Examples

```bash
eggl pf list
eggl pf longhorn
```

See [Configuration](/guide/configuration) for validation rules and setup.
