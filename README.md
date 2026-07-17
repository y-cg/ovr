# ovr

Deep-merge TOML, JSON, and YAML config Layers. Later Layers override earlier ones.

```
ovr base.toml dev.toml local.yaml
ovr base.toml dev.toml local.yaml -o out.toml
ovr explain
```

Merge rules, tombstones, type conflicts, and CLI format resolution are the tested contract printed by **`ovr explain`** (source: [`merge/explain.md`](merge/explain.md)).

## Install

```sh
go install github.com/y-cg/ovr/cmd/ovr@latest
```

Or with Nix:

```sh
nix run github:y-cg/ovr
```

## Usage

```
ovr explain
ovr [--output-format FORMAT] [-o FILE] [--array-append] <files>...
```

| Flag | Default | Description |
|---|---|---|
| `--output-format` / `-f` | `-o` extension, else first Layer | `toml`, `json`, or `yaml`. Errors if it disagrees with a known `-o` extension. |
| `-o` | stdout | Write to a file. A known extension sets format when `-f` is omitted. |
| `--array-append` | off | Concatenate arrays instead of replacing them. |

Full contract: `ovr explain`.
