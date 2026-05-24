# eggl-cli

[![Release](https://img.shields.io/github/v/release/Robert27/eggl-cli?style=flat-square)](https://github.com/Robert27/eggl-cli/releases)
[![CI](https://github.com/Robert27/eggl-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Robert27/eggl-cli/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Robert27/eggl-cli/graph/badge.svg)](https://codecov.io/gh/Robert27/eggl-cli)
[![Go](https://img.shields.io/github/go-mod/go-version/Robert27/eggl-cli?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/github/license/Robert27/eggl-cli?style=flat-square)](LICENSE)

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
eggl pf list
eggl pf longhorn
eggl completion zsh
```

### Environment profiles (`eggl env`)

Switch kubectl context and a mesh VPN target together from `~/.config/eggl/config.yaml`. Each profile uses Tailscale or NetBird:

```yaml
profiles:
  work:
    kube_context: context-work
    tailscale_account: b3e1
  homelab:
    kube_context: context-home
    vpn: netbird
    netbird_profile: homelab
```

`tailscale_account` can be an account id, tailnet name, or email from `tailscale switch --list`. `netbird_profile` is a name from `netbird profile list`. Omit `vpn` for Tailscale (default). `toggle` requires exactly two profiles.

### Port-forwards (`eggl pf`)

Named kubectl port-forwards in the same config file (uses the active kubectl context):

```yaml
port_forwards:
  longhorn:
    namespace: longhorn-system
    resource: svc/longhorn-frontend
    ports: ["8080:80"]
```

`ports` is optional (defaults to `8080:80`). Run `eggl pf list` or `eggl pf longhorn` (blocks until Ctrl+C).

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
