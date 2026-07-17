package main

import (
	"strings"
	"testing"

	"github.com/y-cg/ovr/merge"
)

func TestResolveOutputFormat(t *testing.T) {
	tests := []struct {
		name      string
		explicit  string
		output    string
		firstFile string
		want      merge.Format
		wantErr   string
	}{
		{
			name:      "infer from -o over first input",
			output:    "out.json",
			firstFile: "base.toml",
			want:      merge.JSON,
		},
		{
			name:      "explicit -f without -o",
			explicit:  "yaml",
			firstFile: "base.toml",
			want:       merge.YAML,
		},
		{
			name:      "explicit -f agrees with -o",
			explicit:  "toml",
			output:    "out.toml",
			firstFile: "base.json",
			want:       merge.TOML,
		},
		{
			name:      "explicit -f conflicts with -o",
			explicit:  "toml",
			output:    "out.json",
			firstFile: "base.toml",
			wantErr:   `output format "toml" conflicts with output file extension ".json"`,
		},
		{
			name:      "extensionless -o falls back to first input",
			output:    "out",
			firstFile: "base.toml",
			want:      merge.TOML,
		},
		{
			name:      "no -o no -f uses first input",
			firstFile: "base.yaml",
			want:      merge.YAML,
		},
		{
			name:      "unknown -o extension falls back to first input",
			output:    "out.txt",
			firstFile: "base.toml",
			want:      merge.TOML,
		},
		{
			name:      "explicit -f with unknown -o extension",
			explicit:  "json",
			output:    "out",
			firstFile: "base.toml",
			want:      merge.JSON,
		},
		{
			name:      "yml extension from -o",
			output:    "out.yml",
			firstFile: "base.toml",
			want:      merge.YAML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveOutputFormat(tt.explicit, tt.output, tt.firstFile)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %q, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("format = %q, want %q", got, tt.want)
			}
		})
	}
}
