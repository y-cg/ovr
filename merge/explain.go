package merge

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
)

//go:embed explain.md
var ExplainDoc string

// Example is one executable contract case from an ovr-example fence in explain.md.
type Example struct {
	Format      Format
	Base        string
	Override    string
	Output      string // expected merge result; mutually exclusive with Error
	Error       string // substring that must appear in the merge error
	ArrayAppend bool
}

// ParseExamples extracts ovr-example fences from an explain document.
func ParseExamples(doc string) ([]Example, error) {
	var examples []Example
	rest := doc
	for {
		start := strings.Index(rest, "```ovr-example")
		if start < 0 {
			break
		}
		rest = rest[start+len("```ovr-example"):]
		if len(rest) > 0 && rest[0] == '\n' {
			rest = rest[1:]
		}
		end := strings.Index(rest, "```")
		if end < 0 {
			return nil, fmt.Errorf("unclosed ovr-example fence")
		}
		body := rest[:end]
		rest = rest[end+3:]

		ex, err := parseExampleBody(body)
		if err != nil {
			return nil, err
		}
		examples = append(examples, ex)
	}
	return examples, nil
}

func parseExampleBody(body string) (Example, error) {
	var ex Example
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			return Example{}, fmt.Errorf("ovr-example line missing ':': %q", line)
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		switch key {
		case "format":
			format, err := FormatFromExt(val)
			if err != nil {
				return Example{}, err
			}
			ex.Format = format
		case "base":
			ex.Base = val
		case "override":
			ex.Override = val
		case "output":
			ex.Output = val
		case "error":
			ex.Error = val
		case "options":
			if val != "array-append" {
				return Example{}, fmt.Errorf("unknown options value %q", val)
			}
			ex.ArrayAppend = true
		default:
			return Example{}, fmt.Errorf("unknown ovr-example field %q", key)
		}
	}
	if ex.Format == "" {
		return Example{}, fmt.Errorf("ovr-example missing format")
	}
	if ex.Base == "" || ex.Override == "" {
		return Example{}, fmt.Errorf("ovr-example requires base and override")
	}
	if (ex.Output == "") == (ex.Error == "") {
		return Example{}, fmt.Errorf("ovr-example requires exactly one of output or error")
	}
	return ex, nil
}

// Run executes the example against Merge and checks output or error.
func (ex Example) Run() error {
	opts := Options{}
	if ex.ArrayAppend {
		opts.Arrays = ArrayAppend
	}
	inputs := []Input{
		{Data: []byte(ex.Base), Format: ex.Format},
		{Data: []byte(ex.Override), Format: ex.Format},
	}
	out, err := Merge(inputs, ex.Format, opts)

	if ex.Error != "" {
		if err == nil {
			return fmt.Errorf("expected error containing %q, got success:\n%s", ex.Error, out)
		}
		if !strings.Contains(err.Error(), ex.Error) {
			return fmt.Errorf("error %q does not contain %q", err, ex.Error)
		}
		return nil
	}

	if err != nil {
		return err
	}
	return mapsEqual(ex.Format, out, []byte(ex.Output))
}

func mapsEqual(format Format, got, want []byte) error {
	gotMap, err := parse(Input{Data: got, Format: format})
	if err != nil {
		return fmt.Errorf("parsing merge output: %w\n%s", err, got)
	}
	wantMap, err := parse(Input{Data: want, Format: format})
	if err != nil {
		return fmt.Errorf("parsing expected output: %w\n%s", err, want)
	}
	if !reflect.DeepEqual(gotMap, wantMap) {
		return fmt.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	return nil
}
