package main

import (
	"errors"
	"strings"
	"testing"
)

func TestFormatExplain(t *testing.T) {
	const doc = "# Title\n\nHello **world**.\n"

	t.Run("raw flag forces source even when interactive", func(t *testing.T) {
		got := formatExplain(doc, true, true)
		if got != doc {
			t.Fatalf("got %q, want raw doc", got)
		}
	})

	t.Run("non-interactive returns source", func(t *testing.T) {
		got := formatExplain(doc, false, false)
		if got != doc {
			t.Fatalf("got %q, want raw doc", got)
		}
	})

	t.Run("interactive renders markdown", func(t *testing.T) {
		got := formatExplain(doc, true, false)
		if got == doc {
			t.Fatal("expected rendered output, got raw doc")
		}
		if !strings.Contains(got, "Title") {
			t.Fatalf("rendered output missing title: %q", got)
		}
	})

	t.Run("render failure falls back to source", func(t *testing.T) {
		orig := renderMarkdown
		t.Cleanup(func() { renderMarkdown = orig })
		renderMarkdown = func(string) (string, error) {
			return "", errors.New("boom")
		}
		got := formatExplain(doc, true, false)
		if got != doc {
			t.Fatalf("got %q, want raw doc on render failure", got)
		}
	})
}
