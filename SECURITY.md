# Security Policy

`eggl-cli` is a small personal CLI. It runs locally, does not phone home, and does not handle credentials or remote services.

## Supported versions

Security fixes are applied on the latest release only. Older tags are not maintained.

## Reporting a vulnerability

If you believe you have found a security issue, please open a [private security advisory](https://github.com/Robert27/eggl-cli/security/advisories/new) on GitHub, or email the maintainer via their GitHub profile.

Please include:

- A clear description of the issue
- Steps to reproduce
- Impact (e.g. unexpected file writes, path traversal)
- Your environment (OS, `eggl-cli version`)

I aim to acknowledge reports within a few days. Critical issues in supported releases will be fixed and released when practical.

## Scope

**In scope**

- Issues in this repository’s code (including release artifacts built from it)
- Unsafe file handling (e.g. following symlinks outside the intended tree, corrupting binary files)
- Misleading or dangerous defaults in commands that modify the filesystem

**Out of scope**

- Bugs in third-party tools `eggl doctor` checks (Homebrew, Git, etc.)
- Problems caused by running untrusted binaries not obtained from [GitHub Releases](https://github.com/Robert27/eggl-cli/releases) or the official Homebrew tap
- General hardening requests with no demonstrated exploit

## Safe use

- Install from the official release channel or build from source you trust.
- This project is maintained in spare time; there is no bug bounty program.

Thank you for responsible disclosure.
