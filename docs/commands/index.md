# Command Reference

Run `eggl --help` for the top-level list or `eggl <command> --help` for the flags and examples for a command. The executable is `eggl`; the product is eggl CLI.

| Command | Purpose | External dependency |
| --- | --- | --- |
| `env` | Switch paired Kubernetes and Tailscale profiles | `kubectl`, `tailscale` |
| `pf` | Run named Kubernetes port-forwards | `kubectl` |
| `cd` | Resolve configured directory aliases | None |
| `dedash` | Replace em-dashes in text files | Optional Git integration |
| `eol` | Normalize text files to LF | Optional Git integration |
| `kill` | Find and terminate processes on a TCP port | OS process tools |
| `empty` | Create an empty Git commit | `git` |
| `doctor` | Check local tools and paths | Checks available tools |
| `version` | Print build information | None |
| `completion` | Generate shell completion scripts | None |

All eggl CLI commands support the global `--verbose` flag. It sends operation details to stderr. Set `NO_COLOR` to disable interactive styling.

## Command groups

- [Environment profiles](/commands/env)
- [Kubernetes port-forwards](/commands/pf)
- [Directory aliases](/commands/cd)
- [Repository text files](/commands/text-files)
- [Process cleanup](/commands/kill)
- [Git automation](/commands/empty)
- [Diagnostics and metadata](/commands/diagnostics)
