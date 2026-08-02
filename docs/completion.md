# Shell Completion

eggl-cli generates completion scripts for Bash, Zsh, Fish, and PowerShell.

## Bash

Load for the current session:

```bash
source <(eggl completion bash)
```

Persist it system-wide or in your shell configuration:

```bash
eggl completion bash > /etc/bash_completion.d/eggl
```

## Zsh

```bash
eggl completion zsh > "${fpath[1]}/_eggl"
```

Restart the shell or run `compinit` after writing the file.

## Fish

```bash
eggl completion fish > ~/.config/fish/completions/eggl.fish
```

## PowerShell

Load for the current session:

```powershell
eggl completion powershell | Out-String | Invoke-Expression
```

Write a script for later use:

```powershell
eggl completion powershell > eggl.ps1
```
