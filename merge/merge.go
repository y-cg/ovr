package merge

type Format string

const (
	TOML Format = "toml"
	JSON Format = "json"
	YAML Format = "yaml"
)

type Input struct {
	Data   []byte
	Format Format
}

// Merge deep-merges inputs left-to-right and serializes the result as outputFormat.
// Later inputs win on scalar conflicts. Arrays are replaced, not appended.
// A null value in an override deletes the key (tombstone). Type conflicts are hard errors.
func Merge(inputs []Input, outputFormat Format) ([]byte, error) {
	panic("not implemented")
}
