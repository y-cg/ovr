package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/y-cg/ovr/merge"
)

// fixture reads a file from testdata/<dir>/<name> and wraps it as a merge.Input,
// inferring the format from the file extension.
func fixture(t *testing.T, dir, name string) merge.Input {
	t.Helper()
	path := filepath.Join("testdata", dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture %s: %v", path, err)
	}
	ext := filepath.Ext(name)[1:] // strip leading dot
	return merge.Input{Data: data, Format: merge.Format(ext)}
}

// ==============================================================================
// Basic deep merge — nested objects, later file wins on scalar conflicts
// ==============================================================================

func TestDeepMerge(t *testing.T) {
	out, err := merge.Merge(
		[]merge.Input{
			fixture(t, "deep_merge", "base.toml"),
			fixture(t, "deep_merge", "override.toml"),
		},
		merge.TOML,
	)
	if err != nil {
		t.Fatal(err)
	}
	snaps.MatchSnapshot(t, string(out))
}

// ==============================================================================
// Array replacement — later file wins entirely, no appending
// ==============================================================================

func TestArrayReplacement(t *testing.T) {
	out, err := merge.Merge(
		[]merge.Input{
			fixture(t, "array_replacement", "base.toml"),
			fixture(t, "array_replacement", "override.toml"),
		},
		merge.TOML,
	)
	if err != nil {
		t.Fatal(err)
	}
	snaps.MatchSnapshot(t, string(out))
}

// ==============================================================================
// Null tombstone — null in override deletes the key from output
// ==============================================================================

func TestNullTombstone(t *testing.T) {
	out, err := merge.Merge(
		[]merge.Input{
			fixture(t, "null_tombstone", "base.json"),
			fixture(t, "null_tombstone", "override.json"),
		},
		merge.JSON,
	)
	if err != nil {
		t.Fatal(err)
	}
	snaps.MatchSnapshot(t, string(out))
}

// ==============================================================================
// Type conflict — hard error when override changes the type at a key
// ==============================================================================

func TestTypeConflictError(t *testing.T) {
	_, err := merge.Merge(
		[]merge.Input{
			fixture(t, "type_conflict", "base.toml"),
			fixture(t, "type_conflict", "override.toml"),
		},
		merge.TOML,
	)
	if err == nil {
		t.Fatal("expected type conflict error, got nil")
	}
}

// ==============================================================================
// Cross-format merge — TOML base, JSON override, YAML output
// ==============================================================================

func TestCrossFormatMerge(t *testing.T) {
	out, err := merge.Merge(
		[]merge.Input{
			fixture(t, "cross_format", "base.toml"),
			fixture(t, "cross_format", "override.json"),
		},
		merge.YAML,
	)
	if err != nil {
		t.Fatal(err)
	}
	snaps.MatchSnapshot(t, string(out))
}

// ==============================================================================
// Three-layer merge — each layer narrows closer to the final config
// ==============================================================================

func TestThreeLayerMerge(t *testing.T) {
	out, err := merge.Merge(
		[]merge.Input{
			fixture(t, "three_layer", "01_base.toml"),
			fixture(t, "three_layer", "02_env.toml"),
			fixture(t, "three_layer", "03_local.toml"),
		},
		merge.TOML,
	)
	if err != nil {
		t.Fatal(err)
	}
	snaps.MatchSnapshot(t, string(out))
}
