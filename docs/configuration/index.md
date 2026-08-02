# Configuration

eggl-cli keeps profiles, port-forwards, and directory aliases in one YAML file.

The default path is:

```text
~/.config/eggl/config.yaml
```

When `XDG_CONFIG_HOME` is set, eggl-cli uses `$XDG_CONFIG_HOME/eggl/config.yaml` instead. Every configuration-aware command also accepts `--config` to use a different file.

```bash
eggl env --config ./eggl.yaml show
eggl pf --config ./eggl.yaml list
eggl cd --config ./eggl.yaml list
```

The file must contain at least one profile or directory. Port-forwards are optional.

## Sections

| Section | Used by | Required fields |
| --- | --- | --- |
| `profiles` | `eggl env` | `kube_context`, `tailscale_account` |
| `port_forwards` | `eggl pf` | `namespace`, `resource` |
| `directories` | `eggl cd` | Alias mapped to a path |

See the individual pages for examples and validation rules:

- [Config file](/configuration/config-file)
- [Environment profiles](/configuration/profiles)
- [Port-forwards](/configuration/port-forwards)
- [Directory aliases](/configuration/directories)
