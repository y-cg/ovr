package merge

import (
	"testing"
)

func TestParseExamples_outputAndError(t *testing.T) {
	doc := `
## Merge

` + "```ovr-example\n" + `format: json
base: {"a":1,"b":2}
override: {"b":null}
output: {"a":1}
` + "```\n\n" + "```ovr-example\n" + `format: toml
base: database = "sqlite"
override: database = { host = "localhost" }
error: type conflict
` + "```\n"

	examples, err := ParseExamples(doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(examples) != 2 {
		t.Fatalf("got %d examples, want 2", len(examples))
	}
	if examples[0].Format != JSON || examples[0].Output == "" || examples[0].Error != "" {
		t.Fatalf("first example: %+v", examples[0])
	}
	if examples[1].Format != TOML || examples[1].Error != "type conflict" || examples[1].Output != "" {
		t.Fatalf("second example: %+v", examples[1])
	}
}

func TestExampleRun_tombstone(t *testing.T) {
	ex := Example{
		Format:   JSON,
		Base:     `{"a":1,"b":2}`,
		Override: `{"b":null}`,
		Output:   `{"a":1}`,
	}
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestExampleRun_typeConflict(t *testing.T) {
	ex := Example{
		Format:   TOML,
		Base:     `database = "sqlite"`,
		Override: `database = { host = "localhost" }`,
		Error:    "type conflict",
	}
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestExampleRun_arrayAppend(t *testing.T) {
	ex := Example{
		Format:      TOML,
		Base:        `plugins = ["a", "b"]`,
		Override:    `plugins = ["c"]`,
		Output:      `plugins = ["a", "b", "c"]`,
		ArrayAppend: true,
	}
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestExplainDoc_examples(t *testing.T) {
	examples, err := ParseExamples(ExplainDoc)
	if err != nil {
		t.Fatal(err)
	}
	if len(examples) < 5 {
		t.Fatalf("ExplainDoc has %d examples, want at least 5 contract cases", len(examples))
	}
	for i, ex := range examples {
		if err := ex.Run(); err != nil {
			t.Errorf("example %d: %v", i+1, err)
		}
	}
}
