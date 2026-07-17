# ovr — Project Context

## What It Is

`ovr` is a CLI tool for layered config merging. Given two or more TOML, JSON, or YAML files, it deep-merges them left-to-right so that later files override earlier ones.

Primary use case: layered configuration management (base + environment + local overrides).

## Design Decisions

### Merge semantics
- **Arrays**: replace (later wins entirely — no append)
- **Type conflicts**: hard error — configs must be structurally compatible patches
- **`null` values**: tombstone — setting a key to `null` in an override deletes it from the output

### CLI interface
```
ovr base.toml dev.toml local.yaml           # stdout
ovr base.toml dev.toml local.yaml -o out.toml  # write to file
```
- Positional args, left-to-right precedence
- No config file (`.ovr.toml` etc.) — intentionally kept simple
- Output: stdout by default, `-o <file>` to write to a file
- Output format precedence: `--output-format` / `-f` if set; else known `-o` extension; else first input's format. Hard error when `-f` and a known `-o` extension disagree.

### Language & distribution
- **Go**
- `go install` + GitHub Releases with pre-built binaries

## What's Not Decided Yet

- Array append opt-in syntax or flag (deferred — replace-by-default covers the common case)
- Error message format for type conflicts
- Whether to support `-` (stdin) as an input filename
