# ovr

Deep-merge TOML, JSON, and YAML config files. Later files override earlier ones.

```
ovr base.toml dev.toml local.yaml
ovr base.toml dev.toml local.yaml -o out.toml
```

## Merge semantics

- **Scalars and maps**: later file wins
- **Arrays**: replace entirely — later file wins
- **`null`**: tombstone — removes the key from output
- **Type conflicts**: hard error

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
ovr [--output-format FORMAT] [-o FILE] <files>...
```

| Flag | Default | Description |
|---|---|---|
| `--output-format` | first input's format | Output format: `toml`, `json`, or `yaml` |
| `-o` | stdout | Write output to a file instead |
