# Commands

Complete reference for every `eggl` command. All commands support the global `--verbose` / `-v` flag.

## Command tree

```
eggl
├── completion [bash|zsh|fish|powershell]
├── dedash
├── doctor
├── empty
├── env
│   ├── init
│   ├── show
│   ├── toggle
│   ├── use <profile>
│   └── path
├── eol
├── kill <port>
├── pf [service]
│   └── list
└── version
```

## Quick reference

| Command | Description |
|---------|-------------|
| [`dedash`](/commands/dedash) | Replace em-dashes with hyphens |
| [`eol`](/commands/eol) | Normalize line endings to LF |
| [`empty`](/commands/empty) | Create an empty git commit |
| [`env`](/commands/env) | Switch kubectl + Tailscale profiles |
| [`pf`](/commands/pf) | Kubernetes port-forwards |
| [`kill`](/commands/kill) | Terminate processes on a port |
| [`doctor`](/commands/doctor) | Environment sanity checks |
| [`version`](/commands/version) | Print build info |
| [`completion`](/commands/completion) | Generate shell completions |

## External dependencies

| Command | Requires |
|---------|----------|
| `env show`, `env use`, `env toggle` | `kubectl`, `tailscale` |
| `pf` | `kubectl` |
| `empty` | `git` |
| `dedash --diff`, `eol --diff` | `git` |
| `doctor` | Reports on `go`, `git`, `kubectl`, `tailscale` (optional) |

## Confirmation prompts

`dedash`, `eol`, and `kill` prompt before making changes unless you pass `--yes` or `--dry-run`. In non-interactive environments (CI, pipes), `--yes` is required.
