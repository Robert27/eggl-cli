# Directory Aliases

Directory aliases map short names to paths:

```yaml
directories:
  homelab: ~/projects/homelab
  work: /Users/me/code/work
```

List configured aliases:

```bash
eggl cd list
```

Resolve an alias:

```bash
eggl cd homelab
```

The command prints the path but cannot change the parent shell's working directory. Use command substitution:

```bash
cd "$(eggl cd homelab)"
```

For a shorter command, define a shell function:

```bash
egglcd() { cd "$(eggl cd "$1")" || return; }
```
