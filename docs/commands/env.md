# env

Switch kubectl context and Tailscale account together from paired profiles in your config file.

## Usage

```bash
eggl env <subcommand> [flags]
```

## Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `~/.config/eggl/config.yaml` | Path to config file |

## Subcommands

### `init`

Create an example `config.yaml` if it does not exist.

```bash
eggl env init
```

### `show`

Display the active profile, current kubectl context, Tailscale account, and all configured profiles.

```bash
eggl env show
```

### `toggle`

Flip between exactly two profiles. Requires precisely two profiles in config.

```bash
eggl env toggle
```

### `use <profile>`

Switch to a named profile. Runs `kubectl config use-context` and `tailscale switch`.

```bash
eggl env use homelab
```

### `path`

Print the config file path.

```bash
eggl env path
```

## Configuration

Profiles are defined in `~/.config/eggl/config.yaml`:

```yaml
profiles:
  alpha:
    kube_context: context-a
    tailscale_account: b3e1
  beta:
    kube_context: context-b
    tailscale_account: a7f2
```

See [Configuration](/guide/configuration) for full details.

## Behavior

- Active profile is detected by matching current kube context + Tailscale account ID
- If kube switches successfully but Tailscale fails, the error notes that kube was already switched
- Requires `kubectl` and `tailscale` on `PATH`

## Examples

```bash
eggl env init
eggl env show
eggl env toggle
eggl env use homelab
eggl env path
```
