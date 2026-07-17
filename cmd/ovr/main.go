package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/y-cg/ovr/merge"
)

var cli struct {
	Explain struct{} `cmd:"" help:"Print the tested merge and CLI contract."`

	Merge struct {
		Files        []string `arg:"" name:"file" help:"Input Layers to merge left-to-right. Later Layers win." min:"2"`
		OutputFormat string   `short:"f" name:"output-format" help:"Output format (toml, json, yaml). Defaults to -o extension, else the first Layer's format."`
		Output       string   `short:"o" name:"output" help:"Write output to file instead of stdout. Extension sets format when -f is omitted."`
		ArrayAppend  bool     `name:"array-append" help:"Append arrays instead of replacing them."`
	} `cmd:"" default:"withargs" help:"Deep-merge config Layers left-to-right."`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("ovr"),
		kong.Description("Deep-merge TOML, JSON, and YAML config Layers left-to-right.\n\nSee: ovr explain"),
	)

	var err error
	switch ctx.Command() {
	case "explain":
		_, err = os.Stdout.WriteString(merge.ExplainDoc)
	default:
		err = runMerge()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ovr: %v\n", err)
		os.Exit(1)
	}
}

func runMerge() error {
	m := cli.Merge
	inputs, err := readInputs(m.Files)
	if err != nil {
		return err
	}

	outputFormat, err := resolveOutputFormat(m.OutputFormat, m.Output, m.Files[0])
	if err != nil {
		return err
	}

	opts := merge.Options{}
	if m.ArrayAppend {
		opts.Arrays = merge.ArrayAppend
	}

	out, err := merge.Merge(inputs, outputFormat, opts)
	if err != nil {
		return err
	}

	return writeOutput(out, m.Output)
}

// readInputs reads each file from disk and infers its format from the extension.
func readInputs(paths []string) ([]merge.Input, error) {
	inputs := make([]merge.Input, len(paths))
	for i, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		format, err := merge.FormatFromExt(filepath.Ext(path))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		inputs[i] = merge.Input{Data: data, Format: format}
	}
	return inputs, nil
}

// resolveOutputFormat picks the output format: -f if set, else -o extension
// if known, else the first input file's extension. Errors when -f and a known
// -o extension disagree.
func resolveOutputFormat(explicit, outputPath, firstFile string) (merge.Format, error) {
	fromOutput, outputExt, hasOutputExt := formatFromOutputPath(outputPath)

	if explicit != "" {
		fromFlag, err := merge.FormatFromExt(explicit)
		if err != nil {
			return "", err
		}
		if hasOutputExt && fromFlag != fromOutput {
			return "", fmt.Errorf("output format %q conflicts with output file extension %q", fromFlag, outputExt)
		}
		return fromFlag, nil
	}

	if hasOutputExt {
		return fromOutput, nil
	}

	return merge.FormatFromExt(filepath.Ext(firstFile))
}

// formatFromOutputPath returns a format when path has a known extension.
// Unknown or missing extensions are ignored (ok=false), not errors.
func formatFromOutputPath(path string) (format merge.Format, ext string, ok bool) {
	if path == "" {
		return "", "", false
	}
	ext = filepath.Ext(path)
	if ext == "" {
		return "", "", false
	}
	format, err := merge.FormatFromExt(ext)
	if err != nil {
		return "", ext, false
	}
	return format, ext, true
}

// writeOutput writes to stdout or to a file if a path was given.
func writeOutput(data []byte, path string) error {
	if path == "" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
