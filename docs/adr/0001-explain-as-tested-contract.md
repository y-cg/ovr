# Explain command as the tested agent contract

Agent-facing usage and merge semantics live in `merge/explain.md`, printed by `ovr explain` (with a short pointer from `--help`). README only links to that entry point—it does not restate the rules. Examples in the Merge section use `ovr-example` fences and are executed in tests so the prose cannot drift from `merge.Merge`.

## Considered Options

- **Thicken README only** — agents read it, but nothing ties bullets to runtime behavior.
- **Hand-written `--help` essay** — drifts from `deepMerge` within months.
- **Explain by referencing `testdata/` only** — zero duplication, but `ovr explain` is not self-contained for agents.
- **Embed + `ovr-example` contract (chosen)** — one file, CLI discovery, CI-enforced examples (`output` or `error`).

## Consequences

- `merge` embeds `explain.md` and owns example tests; `cmd/ovr` prints it via `ovr explain` (raw for non-TTY/`--raw`; glamour-rendered for interactive TTY). Presentation must not become a second untested semantics source.
- Array replace is the default contract; array append is documented with an explicit example tied to `--array-append`.
- Tombstones require a format that can express `null` (JSON/YAML, not TOML).
