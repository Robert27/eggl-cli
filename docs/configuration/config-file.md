# Config File

The complete shape of the file is:

```yaml
profiles:
  homelab:
    kube_context: context-a
    tailscale_account: tailnet-a

port_forwards:
  dashboard:
    namespace: monitoring
    resource: svc/grafana
    ports: ["3000:80"]

directories:
  homelab: ~/projects/homelab
  work: /Users/me/code/work
```

## Path resolution

`~` is expanded only at the beginning of a directory path. Both `~/projects/app` and absolute paths are valid. Paths such as `~other-user/app` are rejected.

## Validation

Configuration is validated before a command uses it:

- Profile names and aliases must not be empty.
- Profiles require both `kube_context` and `tailscale_account`.
- Profile kube contexts cannot contain whitespace or start with `-`.
- Kubernetes namespaces must use lowercase DNS-label syntax.
- Resources must use `kind/name`, such as `svc/my-app` or `deployment/api`.
- Port mappings must use `local:remote` with ports from `1` through `65535`.
- Duplicate profile targets are rejected.
- Files larger than 1 MiB are rejected.

Use the starter file as a safe baseline:

```bash
eggl env init
```

eggl-cli refuses to overwrite an existing file.
