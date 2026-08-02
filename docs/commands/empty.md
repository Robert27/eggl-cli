# `eggl empty`

Create an empty Git commit. This is useful for retriggering a CI pipeline without changing tracked files.

```bash
eggl empty
eggl empty --push
eggl empty -p
```

| Flag | Description |
| --- | --- |
| `-p, --push` | Push the new commit after creating it |

The command runs in the current repository and reports the commit hash. With `--push`, the configured Git remote and credentials are used.
