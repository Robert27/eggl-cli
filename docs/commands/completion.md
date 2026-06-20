# completion

Generate shell completion scripts for eggl.

## Usage

```bash
eggl completion [bash|zsh|fish|powershell]
```

## Arguments

| Shell | Description |
|-------|-------------|
| `bash` | Bash completion script |
| `zsh` | Zsh completion script |
| `fish` | Fish completion script |
| `powershell` | PowerShell completion script |

## Examples

```bash
eggl completion zsh > "${fpath[1]}/_eggl"
source <(eggl completion bash)
eggl completion fish > ~/.config/fish/completions/eggl.fish
```

See the [Shell Completion guide](/guide/completion) for full setup instructions.

## Completions include

- All commands and subcommands
- Configured port-forward names for `eggl pf`
- Profile names where applicable
