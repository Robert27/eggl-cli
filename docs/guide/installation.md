# Installation

## Homebrew (recommended)

```bash
brew tap Robert27/tap
brew install eggl-cli
```

Shell completions for bash, zsh, and fish are installed automatically.

## GitHub Releases

Download the archive for your OS/arch from [GitHub Releases](https://github.com/Robert27/eggl-cli/releases), extract `eggl` (or `eggl.exe` on Windows), and put it on your `PATH`.

Supported platforms:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

## Go install

With Go 1.26+ installed:

```bash
go install github.com/Robert27/eggl-cli@latest
```

## From source

```bash
git clone https://github.com/Robert27/eggl-cli.git
cd eggl-cli
make install
```

Or build a local binary without installing:

```bash
make build
./bin/eggl --help
```

## Verify installation

```bash
eggl version
eggl doctor
```

`doctor` checks for optional tools like `kubectl` and `tailscale`. Missing tools are reported but only matter for commands that use them.
