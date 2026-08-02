# Install

Choose the distribution channel that fits your machine.

## Homebrew

```bash
brew tap roberteggl/tap
brew install eggl-cli
```

Homebrew installs Bash, Zsh, and Fish completions automatically.

## GitHub Releases

Download the archive for your operating system and architecture from the [GitHub Releases](https://github.com/roberteggl/eggl-cli/releases) page.

Extract `eggl` on macOS or Linux, or `eggl.exe` on Windows, and put the binary on your `PATH`.

The published targets are:

| Operating system | Architectures |
| --- | --- |
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64, arm64 |

## Go install

```bash
go install github.com/roberteggl/eggl-cli@latest
```

## Build from source

```bash
git clone https://github.com/roberteggl/eggl-cli.git
cd eggl-cli
make install
```

Or build a local binary without installing it:

```bash
make build
./bin/eggl --help
```

## Verify the installation

```bash
eggl version
eggl doctor
```

`doctor` reports missing optional integrations such as `kubectl` or Tailscale. Those tools are only required by the commands that use them.
