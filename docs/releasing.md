# Releasing

Releases are created by pushing a version tag.

```bash
git tag v0.3.0
git push origin v0.3.0
```

The release workflow runs GoReleaser after generating notes with Git Cliff.

## Published artifacts

- GitHub Release with checksums
- Linux, macOS, and Windows archives
- amd64 and arm64 binaries
- Homebrew formula update in `roberteggl/homebrew-tap`
- Shell completions in the Homebrew formula

The release build uses `CGO_ENABLED=0` and embeds the version, commit, and build date into the binary.

## Required secret

The repository needs `TAP_GITHUB_TOKEN` so GoReleaser can update the external Homebrew tap. The default repository token cannot write to another repository.

## Test a release locally

Build a snapshot without publishing:

```bash
make release-snapshot
```

The regular pre-release checks are:

```bash
make fmt
make check
```

## Documentation deployment

The docs site has its own workflow. Pull requests build the site and catch dead links. Pushes to `main` deploy the generated Pages artifact to [cli.eggl.dev](https://cli.eggl.dev).
