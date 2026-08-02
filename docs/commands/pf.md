# `eggl pf`

Run a configured Kubernetes port-forward using the active kubectl context.

```bash
eggl pf list
eggl pf grafana
eggl pf grafana --open
```

## Flags

| Flag | Description |
| --- | --- |
| `--config <path>` | Use a specific config file |
| `-o, --open` | Open `http://localhost:<port>` after the tunnel is ready |

`eggl pf list` prints configured names, namespaces, resources, and ports. Calling `eggl pf` without a service prints command help.

The service argument supports shell completion from the names in the config file.
