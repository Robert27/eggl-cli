# Port-forwards

Named port-forwards use the active Kubernetes context:

```yaml
port_forwards:
  grafana:
    namespace: monitoring
    resource: svc/grafana
    ports: ["3000:80"]
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
```

`resource` accepts Kubernetes kinds such as `svc`, `service`, `deployment`, `pod`, `statefulset`, `daemonset`, `job`, and `cronjob`.

If `ports` is omitted, eggl-cli uses `8080:80`.

## List forwards

```bash
eggl pf list
```

## Start a forward

```bash
eggl pf grafana
```

The command remains attached to the port-forward. Press `Ctrl+C` to stop it.

Open the local endpoint after the tunnel is ready:

```bash
eggl pf grafana --open
```

The `--open` flag opens `http://localhost:<local-port>` in the default browser.
