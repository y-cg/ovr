package merge_test

import (
	"testing"

	"github.com/y-cg/ovr/merge"
)

func TestDeepMerge(t *testing.T) {
	mergeAndSnap(t, []fixture{
		load(t, "deep_merge", "base.toml"),
		load(t, "deep_merge", "override.toml"),
	}, merge.TOML, merge.Options{})
}

func TestArrayReplacement(t *testing.T) {
	mergeAndSnap(t, []fixture{
		load(t, "array_replacement", "base.toml"),
		load(t, "array_replacement", "override.toml"),
	}, merge.TOML, merge.Options{})
}

func TestArrayAppend(t *testing.T) {
	mergeAndSnap(t, []fixture{
		load(t, "array_append", "base.toml"),
		load(t, "array_append", "override.toml"),
	}, merge.TOML, merge.Options{Arrays: merge.ArrayAppend})
}

func TestNullTombstone(t *testing.T) {
	mergeAndSnap(t, []fixture{
		load(t, "null_tombstone", "base.json"),
		load(t, "null_tombstone", "override.json"),
	}, merge.JSON, merge.Options{})
}

func TestTypeConflictError(t *testing.T) {
	_, err := merge.Merge([]merge.Input{
		load(t, "type_conflict", "base.toml").Input,
		load(t, "type_conflict", "override.toml").Input,
	}, merge.TOML, merge.Options{})
	if err == nil {
		t.Fatal("expected type conflict error, got nil")
	}
}

func TestCrossFormatMerge(t *testing.T) {
	mergeAndSnap(t, []fixture{
		load(t, "cross_format", "base.toml"),
		load(t, "cross_format", "override.json"),
	}, merge.YAML, merge.Options{})
}

func TestThreeLayerMerge(t *testing.T) {
	mergeAndSnap(t, []fixture{
		load(t, "three_layer", "01_base.toml"),
		load(t, "three_layer", "02_env.toml"),
		load(t, "three_layer", "03_local.toml"),
	}, merge.TOML, merge.Options{})
}
