package merge_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/y-cg/ovr/merge"
)

// fixture holds a parsed input alongside its filename, so the snapshot helper
// can render the original source alongside the merged output.
type fixture struct {
	merge.Input
	name string
}

// load reads testdata/<dir>/<name>, infers the format from the extension, and
// returns a fixture ready to pass to mergeAndSnap.
func load(t *testing.T, dir, name string) fixture {
	t.Helper()
	path := filepath.Join("testdata", dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("load %s: %v", path, err)
	}
	ext := filepath.Ext(name)[1:]
	return fixture{
		Input: merge.Input{Data: data, Format: merge.Format(ext)},
		name:  name,
	}
}

// mergeAndSnap runs merge.Merge on the given fixtures, then compares (or writes)
// a snapshot at testdata/snapshots/<TestName>.snap.
// Run with UPDATE_SNAPS=true to create or refresh snapshots.
func mergeAndSnap(t *testing.T, fixtures []fixture, outputFormat merge.Format) {
	t.Helper()

	inputs := make([]merge.Input, len(fixtures))
	for i, f := range fixtures {
		inputs[i] = f.Input
	}

	out, err := merge.Merge(inputs, outputFormat)
	if err != nil {
		t.Fatal(err)
	}

	snap := renderSnapshot(fixtures, string(out), outputFormat)
	snapPath := filepath.Join("testdata", "snapshots", t.Name()+".snap")

	if os.Getenv("UPDATE_SNAPS") == "true" {
		if err := os.MkdirAll(filepath.Dir(snapPath), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(snapPath), err)
		}
		if err := os.WriteFile(snapPath, []byte(snap), 0o644); err != nil {
			t.Fatalf("write snapshot: %v", err)
		}
		t.Logf("snapshot written: %s", snapPath)
		return
	}

	existing, err := os.ReadFile(snapPath)
	if os.IsNotExist(err) {
		t.Fatalf("snapshot not found — run with UPDATE_SNAPS=true to create it:\n  %s", snapPath)
	}
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	if string(existing) != snap {
		t.Fatalf("snapshot mismatch for %s\n\ngot:\n%s\nwant:\n%s", snapPath, snap, existing)
	}
}

func renderSnapshot(fixtures []fixture, output string, outputFormat merge.Format) string {
	var sb strings.Builder
	for _, f := range fixtures {
		fmt.Fprintf(&sb, "=== %s ===\n%s\n\n", f.name, strings.TrimRight(string(f.Data), "\n"))
	}
	fmt.Fprintf(&sb, "=== output (%s) ===\n%s\n", outputFormat, strings.TrimRight(output, "\n"))
	return sb.String()
}
