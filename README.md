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
eggl completion zsh
```

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
