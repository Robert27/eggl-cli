# Shell Completion

eggl supports shell completion for bash, zsh, fish, and PowerShell.

## Homebrew

Completions are installed automatically when you install via Homebrew.

## Manual setup

### Zsh

```bash
eggl completion zsh > "${fpath[1]}/_eggl"
```

Restart your shell or run `compinit`.

### Bash

```bash
eggl completion bash > /etc/bash_completion.d/eggl
```

For the current session only:

```bash
source <(eggl completion bash)
```

### Fish

```bash
eggl completion fish > ~/.config/fish/completions/eggl.fish
```

### PowerShell

```powershell
eggl completion powershell | Out-String | Invoke-Expression
```

## What gets completed?

- All commands and subcommands
- Configured port-forward names (`eggl pf <tab>`)
- Profile names where applicable
