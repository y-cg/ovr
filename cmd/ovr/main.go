package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/y-cg/ovr/merge"
)

var cli struct {
	Files        []string `arg:"" name:"file" help:"Input files to merge left-to-right. Later files win." min:"2"`
	OutputFormat string   `short:"f" name:"output-format" help:"Output format (toml, json, yaml). Defaults to the first file's format."`
	Output       string   `short:"o" name:"output" help:"Write output to file instead of stdout."`
	ArrayAppend  bool     `name:"array-append" help:"Append arrays instead of replacing them."`
}

func main() {
	kong.Parse(&cli,
		kong.Name("ovr"),
		kong.Description("Deep-merge TOML, JSON, and YAML config files left-to-right."),
	)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ovr: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	inputs, err := readInputs(cli.Files)
	if err != nil {
		return err
	}

	outputFormat, err := resolveOutputFormat(cli.OutputFormat, cli.Files[0])
	if err != nil {
		return err
	}

	opts := merge.Options{}
	if cli.ArrayAppend {
		opts.Arrays = merge.ArrayAppend
	}

	out, err := merge.Merge(inputs, outputFormat, opts)
	if err != nil {
		return err
	}

	return writeOutput(out, cli.Output)
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

// resolveOutputFormat returns the explicit format if given, otherwise infers
// it from the first input file's extension.
func resolveOutputFormat(explicit, firstFile string) (merge.Format, error) {
	if explicit != "" {
		return merge.FormatFromExt(explicit)
	}
	return merge.FormatFromExt(filepath.Ext(firstFile))
}

// writeOutput writes to stdout or to a file if a path was given.
func writeOutput(data []byte, path string) error {
	if path == "" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
