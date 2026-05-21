# eggl-cli

A general-purpose helper CLI built with Go and [Cobra](https://github.com/spf13/cobra).

## Install

**Homebrew:**

```bash
brew tap Robert27/tap
brew install eggl-cli
```

Shell completions for bash, zsh, and fish are installed automatically.

**Without Homebrew:** download the archive for your OS/arch from [GitHub Releases](https://github.com/Robert27/eggl-cli/releases), extract `eggl` (or `eggl.exe` on Windows), and put it on your `PATH`. Or, with Go installed:

```bash
go install github.com/Robert27/eggl-cli@latest
```

From source:

```bash
make install
```

Or build a local binary:

```bash
make build
./bin/eggl --help
```

## Usage

```bash
eggl --help
eggl version
eggl version --short
eggl doctor
eggl doctor --verbose
eggl env init
eggl env show
eggl env toggle
eggl completion zsh
```

### Environment profiles (`eggl env`)

Switch kubectl context and Tailscale account together from `~/.config/eggl/config.yaml`:

```yaml
profiles:
  alpha:
    kube_context: context-a
    tailscale_account: b3e1
  beta:
    kube_context: context-b
    tailscale_account: a7f2
```

`tailscale_account` can be an account id, tailnet name, or email from `tailscale switch --list`. `toggle` requires exactly two profiles.

## Shell completion

If you installed via Homebrew, completions are set up for you. Otherwise:

**Zsh:**

```bash
eggl completion zsh > "${fpath[1]}/_eggl"
```

**Bash:**

```bash
eggl completion bash > /etc/bash_completion.d/eggl
# or for the current session:
source <(eggl completion bash)
```

**Fish:**

```bash
eggl completion fish > ~/.config/fish/completions/eggl.fish
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, project layout, and release instructions.
