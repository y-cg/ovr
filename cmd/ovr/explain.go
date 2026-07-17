package main

import (
	"io"

	"github.com/charmbracelet/glamour"
)

// renderMarkdown is the interactive renderer; swapped in tests for failure paths.
var renderMarkdown = func(doc string) (string, error) {
	return glamour.Render(doc, "auto")
}

// formatExplain returns the explain document for stdout.
// Interactive TTY gets glamour-rendered markdown; non-TTY, --raw, or
// render failure gets the raw contract source.
func formatExplain(doc string, interactive, raw bool) string {
	if raw || !interactive {
		return doc
	}
	rendered, err := renderMarkdown(doc)
	if err != nil {
		return doc
	}
	return rendered
}

func writeExplain(w io.Writer, doc string, interactive, raw bool) error {
	_, err := io.WriteString(w, formatExplain(doc, interactive, raw))
	return err
}
