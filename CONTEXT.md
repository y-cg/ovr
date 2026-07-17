# ovr

CLI for layered config merging across TOML, JSON, and YAML. Later layers override earlier ones.

## Language

**Layer**:
One input file in a merge, ordered by CLI argument position (left to right).
_Avoid_: config file (as a domain term), document

**Override**:
A later Layer relative to an earlier one; on conflict, the override wins under merge rules.
_Avoid_: patch, update file

**Deep merge**:
Recursive merge of maps; non-map values follow Array replace/append, Tombstone, and Type conflict rules—not a blind combine.
_Avoid_: mix, combine, shallow merge

**Tombstone**:
A `null` value in an Override that deletes that key from the result. Only available when the Layer's format can express null (JSON/YAML; not TOML).
_Avoid_: unset, delete key (as the formal name)

**Type conflict**:
At the same key, one side is a map and the other is not. This is a hard error.
_Avoid_: schema mismatch, shape error

**Array replace**:
Default array rule: the Override array replaces the earlier array entirely.
_Avoid_: overwrite (informal only)

**Array append**:
Opt-in via `--array-append`: concatenate earlier and Override arrays.
_Avoid_: merge arrays, extend (as the formal name)

## Design Decisions

### CLI interface
```
ovr base.toml dev.toml local.yaml           # stdout
ovr base.toml dev.toml local.yaml -o out.toml  # write to file
ovr explain                                 # merge + CLI contract (tested)
```
- Positional args, left-to-right Layer order
- No config file (`.ovr.toml` etc.) — intentionally kept simple
- Output: stdout by default, `-o <file>` to write to a file
- Output format precedence: `--output-format` / `-f` if set; else known `-o` extension; else first Layer's format. Hard error when `-f` and a known `-o` extension disagree.
- Agent-facing contract: `merge` embeds `explain.md`; `ovr explain` prints it; README points at it only — see [ADR 0001](docs/adr/0001-explain-as-tested-contract.md)

### Language & distribution
- **Go**
- `go install` + GitHub Releases with pre-built binaries

## What's Not Decided Yet

- Error message format for type conflicts (stable substring `type conflict` is enough for contract tests)
- Whether to support `-` (stdin) as an input Layer
