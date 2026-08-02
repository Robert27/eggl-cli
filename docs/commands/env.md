# `eggl env`

Switch a Kubernetes context and Tailscale account as one named environment.

```bash
eggl env init
eggl env show
eggl env use homelab
eggl env toggle
```

## Subcommands

| Command | Description |
| --- | --- |
| `init` | Create an example config file if one does not exist |
| `show` | Show active profile and current state |
| `use <profile>` | Apply a named profile |
| `toggle` | Flip between exactly two profiles |
| `path` | Print the resolved config path |

## Flags

`env` has a persistent `--config <path>` flag. It applies to all subcommands and overrides the default config path.

`toggle` requires exactly two profiles. `use` is the explicit option when there are more than two profiles or when the current state is not recognized.
