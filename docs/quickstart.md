# Quickstart

This walkthrough creates two environment profiles, a named port-forward, and a directory alias.

## 1. Create the starter config

```bash
eggl env init
```

The default path is `~/.config/eggl/config.yaml`. Use `eggl env path` to print the resolved path.

## 2. Edit the config

Replace the example values with contexts and account references that exist on your machine:

```yaml
profiles:
  homelab:
    kube_context: context-a
    tailscale_account: tailnet-a
  work:
    kube_context: context-b
    tailscale_account: work@example.com

port_forwards:
  dashboard:
    namespace: monitoring
    resource: svc/grafana
    ports: ["3000:80"]

directories:
  homelab: ~/projects/homelab
  work: ~/projects/work
```

Account references can be an ID, tailnet name, or email shown by `tailscale switch --list`.

## 3. Switch environments

```bash
eggl env show
eggl env use homelab
```

The command changes the kubectl context and Tailscale account only when either value differs from the target profile.

For exactly two profiles, `toggle` switches to the other one:

```bash
eggl env toggle
```

## 4. Start a port-forward

```bash
eggl pf list
```

The forward uses the active kubectl context and stays attached until you press `Ctrl+C`.

## 5. Use a directory alias

`eggl cd` prints a path so your shell can change directory:

```bash
cd "$(eggl cd homelab)"
```

For repeated use, add a wrapper:

```bash
egglcd() { cd "$(eggl cd "$1")" || return; }
```
