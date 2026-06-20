# empty

Create an empty git commit. Useful for retriggering CI pipelines without code changes.

## Usage

```bash
eggl empty [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--push` | `-p` | `false` | Push to remote after creating the commit |

## Behavior

- Must be run inside a git work tree
- Default commit message: `chore: empty commit`
- Prints the new commit hash

## Examples

```bash
# Create empty commit locally
eggl empty

# Create and push
eggl empty --push
eggl empty -p
```

## Output

```
empty commit abc1234
pushed          # only with --push
```
