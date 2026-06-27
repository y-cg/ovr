package merge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Format string

const (
	TOML Format = "toml"
	JSON Format = "json"
	YAML Format = "yaml"
)

type ArrayMode int

const (
	// ArrayReplace is the default: the override array wins entirely.
	ArrayReplace ArrayMode = iota
	// ArrayAppend concatenates base and override arrays.
	ArrayAppend
)

type Options struct {
	Arrays ArrayMode
}

type Input struct {
	Data   []byte
	Format Format
}

// Merge deep-merges inputs left-to-right and serializes the result as outputFormat.
// Later inputs win on scalar conflicts. Type conflicts are hard errors.
// A null value in an override deletes the key (tombstone).
// Array behavior is controlled by opts.Arrays (replace by default).
func Merge(inputs []Input, outputFormat Format, opts Options) ([]byte, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("no inputs")
	}

	result, err := parse(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("input 1: %w", err)
	}

	for i, input := range inputs[1:] {
		next, err := parse(input)
		if err != nil {
			return nil, fmt.Errorf("input %d: %w", i+2, err)
		}
		result, err = deepMerge(result, next, opts)
		if err != nil {
			return nil, fmt.Errorf("merging input %d: %w", i+2, err)
		}
	}

	return serialize(result, outputFormat)
}

// ==============================================================================
// Parsing
// ==============================================================================

func parse(input Input) (map[string]any, error) {
	var m map[string]any
	switch input.Format {
	case JSON:
		if err := json.Unmarshal(input.Data, &m); err != nil {
			return nil, err
		}
	case TOML:
		if err := toml.Unmarshal(input.Data, &m); err != nil {
			return nil, err
		}
	case YAML:
		if err := yaml.Unmarshal(input.Data, &m); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown format %q", input.Format)
	}
	return m, nil
}

// ==============================================================================
// Deep merge
// ==============================================================================

// deepMerge merges override onto base, returning a new map.
// Keys present only in base are kept. Keys present only in override are added.
// Keys present in both are merged according to their types:
//   - Both maps      → recurse
//   - Map vs non-map → hard error (structural conflict)
//   - nil override   → tombstone, delete the key
//   - Both arrays    → replace or append depending on opts.Arrays
//   - Anything else  → override wins
func deepMerge(base, override map[string]any, opts Options) (map[string]any, error) {
	result := make(map[string]any, len(base))
	for k, v := range base {
		result[k] = v
	}

	for k, ov := range override {
		// nil is the tombstone: remove the key entirely.
		if ov == nil {
			delete(result, k)
			continue
		}

		bv, exists := result[k]
		if !exists {
			result[k] = ov
			continue
		}

		bMap, bIsMap := bv.(map[string]any)
		oMap, oIsMap := ov.(map[string]any)

		switch {
		case bIsMap && oIsMap:
			// Both sides are objects — recurse into them.
			merged, err := deepMerge(bMap, oMap, opts)
			if err != nil {
				return nil, fmt.Errorf(".%s%w", k, err)
			}
			result[k] = merged

		case bIsMap != oIsMap:
			// One side is an object and the other is a scalar/array.
			// This is a structural type conflict that we cannot resolve safely.
			return nil, fmt.Errorf(": type conflict at key %q (base is %T, override is %T)", k, bv, ov)

		default:
			// Both are scalars or both are arrays.
			if opts.Arrays == ArrayAppend {
				bSlice, bIsSlice := toSlice(bv)
				oSlice, oIsSlice := toSlice(ov)
				if bIsSlice && oIsSlice {
					result[k] = append(bSlice, oSlice...)
					continue
				}
			}
			// Scalar, or ArrayReplace mode: override wins.
			result[k] = ov
		}
	}

	return result, nil
}

func toSlice(v any) ([]any, bool) {
	s, ok := v.([]any)
	return s, ok
}

// ==============================================================================
// Serialization
// ==============================================================================

func serialize(m map[string]any, format Format) ([]byte, error) {
	switch format {
	case JSON:
		return json.MarshalIndent(m, "", "  ")

	case TOML:
		var buf bytes.Buffer
		if err := toml.NewEncoder(&buf).Encode(m); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil

	case YAML:
		return yaml.Marshal(m)

	default:
		return nil, fmt.Errorf("unknown output format %q", format)
	}
}

// FormatFromExt infers a Format from a file extension (with or without leading dot).
func FormatFromExt(ext string) (Format, error) {
	switch strings.TrimPrefix(strings.ToLower(ext), ".") {
	case "toml":
		return TOML, nil
	case "json":
		return JSON, nil
	case "yaml", "yml":
		return YAML, nil
	default:
		return "", fmt.Errorf("unsupported file extension %q", ext)
	}
}
