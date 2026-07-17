# Contributing to eggl-cli

## Project layout

```
main.go           # entrypoint
cmd/              # Cobra commands (thin orchestration)
internal/         # business logic by domain
Makefile          # build, install, test
```

## Development

```bash
make fmt && make check   # required before committing
make build VERSION=v0.2.0
```

Enable the git pre-commit hook (runs `make check` on each commit):

```bash
make hooks
```

## Adding a new command

1. Create business logic in `internal/<domain>/` (keep commands thin).
2. Add a command file in `cmd/<name>.go` following the `doctor` pattern:

```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    RunE: func(cmd *cobra.Command, args []string) error {
        slog.Debug("running mycommand")
        // call internal/<domain> logic
        return nil
    },
    SilenceUsage: true,
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

3. Run `make fmt && make check` and `make build`.

Alternatively, scaffold with [cobra-cli](https://github.com/spf13/cobra-cli):

```bash
go install github.com/spf13/cobra-cli@latest
cobra-cli add mycommand
```

## Releasing

Releases are automated with [GoReleaser](https://goreleaser.com/) when a `v*` tag is pushed. The workflow builds cross-platform archives, publishes a GitHub Release, and updates the [homebrew-tap](https://github.com/roberteggl/homebrew-tap) formula.

**One-time setup:** add a `TAP_GITHUB_TOKEN` repository secret in `eggl-cli` with `repo` scope so Goreleaser can push formula updates to `roberteggl/homebrew-tap`. The default `GITHUB_TOKEN` cannot write to other repositories.

Test a release locally without publishing:

```bash
make release-snapshot
```
