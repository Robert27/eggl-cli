# AGENTS

## Commands

```bash
make check             # fmt-check + test (required before every commit)
make fmt               # gofmt -w .
make test              # go test ./...
make build             # bin/eggl (VERSION=dev by default)
make build VERSION=v0.2.0
go vet ./...
```

Use `head`/`tail` to limit output tokens.

## Before commit

Always run formatting and tests before creating a commit:

```bash
make fmt && make check
```

A pre-commit hook (`.githooks/pre-commit`) runs `make check` automatically. One-time setup:

```bash
make hooks
```

## Release

Push a `v*` tag; GitHub Actions runs GoReleaser (cross-platform archives, GitHub Release, Homebrew cask update via `TAP_GITHUB_TOKEN`).

```bash
git tag vX.Y.Z && git push origin vX.Y.Z
```

Test locally without publishing:

```bash
make release-snapshot
```

## Project layout

```
main.go           # entrypoint → cmd.Execute()
cmd/              # Cobra commands (thin orchestration)
internal/         # business logic by domain (doctor, dedash, ui, …)
Makefile          # build, install, test, release-snapshot
```

## Adding a command

1. Put logic in `internal/<domain>/`.
2. Add `cmd/<name>.go` following `doctor` / `dedash`: `RunE`, `SilenceUsage: true`, `slog.Debug` at start, register in `init()`.
3. Use `internal/ui` for styled help, version, and doctor output when appropriate.
4. Run `make fmt && make check` and `make build`.

## Rules

### Consistency

- Follow existing Cobra patterns in `cmd/`.
- Reuse `internal/ui` and domain packages; keep commands thin.

### General

- **Before every commit:** run `make fmt && make check` (or rely on the pre-commit hook after `make hooks`).
- No AI-generated comments; only meaningful comments for non-obvious logic.
- No unnecessary Markdown files unless requested.
- Keep agent output compact and token-efficient (comments and prose-not code logic).
