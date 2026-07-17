# ovr contract

Agent-facing merge and CLI contract. Examples in `ovr-example` fences are executed by tests.

## Merge

Layers are merged left to right. Later Layers override earlier ones (**Deep merge**).

### Scalars and maps

Maps recurse. Scalars: the Override wins.

```ovr-example
format: toml
base: database = { host = "localhost", port = 5432 }
override: database = { port = 3306 }
output: database = { host = "localhost", port = 3306 }
```

### Array replace (default)

The Override array replaces the earlier array entirely.

```ovr-example
format: toml
base: plugins = ["a", "b"]
override: plugins = ["c"]
output: plugins = ["c"]
```

### Array append

With `--array-append`, arrays concatenate.

```ovr-example
format: toml
options: array-append
base: plugins = ["a", "b"]
override: plugins = ["c"]
output: plugins = ["a", "b", "c"]
```

### Tombstone

A `null` in an Override deletes that key. Only when the format can express null (JSON/YAML; not TOML).

```ovr-example
format: json
base: {"a": 1, "b": 2}
override: {"b": null}
output: {"a": 1}
```

### Type conflict

One side is a map and the other is not → hard error.

```ovr-example
format: toml
base: database = "sqlite"
override: database = { host = "localhost" }
error: type conflict
```

## CLI

```
ovr <layer> <layer> [<layer>...]
ovr explain
```

- **Layer order**: positional args, left to right. Later Layers win under Merge rules.
- **Stdout** by default; `-o FILE` writes a file.
- **Output format**: `-f` / `--output-format` if set; else a known `-o` extension (`.toml`, `.json`, `.yaml`, `.yml`); else the first Layer's format. If `-f` and a known `-o` extension disagree → hard error.
- **`--array-append`**: use Array append instead of Array replace.
- Cross-format Layers are allowed; serialization format follows the rules above, not "last Layer's format".
