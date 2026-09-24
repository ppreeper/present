// Copyright 2026 The present Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package present

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveTheme(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "custom.css"), []byte("/* x */"), 0644); err != nil {
		t.Fatal(err)
	}
	doc := filepath.Join(dir, "deck.slide")

	tests := []struct {
		value string
		want  string
		err   string
	}{
		{value: "dark", want: "/static/themes/dark.css"},
		{value: "DARK", want: "/static/themes/dark.css"},
		{value: "  default  ", want: "/static/themes/default.css"},
		{value: "high-contrast", want: "/static/themes/high-contrast.css"},
		{value: "custom.css", want: "custom.css"},
		{value: "sub/custom.css", want: "sub/custom.css", err: "no such file"},
		{value: "/static/themes/mono.css", want: "/static/themes/mono.css"},
		{value: "https://example.com/theme.css", want: "https://example.com/theme.css"},
		{value: "nope", err: `unknown theme "nope"`},
		{value: "missing.css", err: "no such file"},
		{value: "", err: "empty theme value"},
	}
	for _, tt := range tests {
		got, err := resolveTheme(tt.value, doc)
		if tt.err != "" {
			if err == nil {
				t.Errorf("resolveTheme(%q) = %q, want error %q", tt.value, got, tt.err)
			} else if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("resolveTheme(%q) error = %v, want %q", tt.value, err, tt.err)
			}
			continue
		}
		if err != nil {
			t.Errorf("resolveTheme(%q): %v", tt.value, err)
			continue
		}
		if got != tt.want {
			t.Errorf("resolveTheme(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestResolveBackground(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bg.png"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	doc := filepath.Join(dir, "deck.slide")

	tests := []struct {
		value string
		want  string
		err   string
	}{
		{value: "bg.png", want: `--slide-background-image:url("bg.png")`},
		{value: "/images/bg.png", want: `--slide-background-image:url("/images/bg.png")`},
		{value: "https://example.com/bg.png", want: `--slide-background-image:url("https://example.com/bg.png")`},
		{value: "missing.png", err: "no such file"},
		{value: "", err: "empty background value"},
	}
	for _, tt := range tests {
		got, err := resolveBackground(tt.value, doc)
		if tt.err != "" {
			if err == nil {
				t.Errorf("resolveBackground(%q) = %q, want error %q", tt.value, got, tt.err)
			} else if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("resolveBackground(%q) error = %v, want %q", tt.value, err, tt.err)
			}
			continue
		}
		if err != nil {
			t.Errorf("resolveBackground(%q): %v", tt.value, err)
			continue
		}
		if string(got) != tt.want {
			t.Errorf("resolveBackground(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestThemeHeaderParse(t *testing.T) {
	const body = "\nName\n\n## Section\n\nhi\n"

	t.Run("resolves after parse", func(t *testing.T) {
		src := "# Title\nSubtitle\nTheme: dark\n\n" + body
		doc, err := Parse(strings.NewReader(src), "deck.slide", 0)
		if err != nil {
			t.Fatal(err)
		}
		if doc.Theme != "dark" {
			t.Errorf("Theme = %q, want %q", doc.Theme, "dark")
		}
		if doc.ThemeHref != "/static/themes/dark.css" {
			t.Errorf("ThemeHref = %q, want %q", doc.ThemeHref, "/static/themes/dark.css")
		}
	})

	t.Run("unknown theme errors", func(t *testing.T) {
		src := "# Title\nSubtitle\nTheme: nope\n\n" + body
		if _, err := Parse(strings.NewReader(src), "deck.slide", 0); err == nil {
			t.Fatal("expected an error for an unknown theme")
		}
	})

	t.Run("empty theme errors", func(t *testing.T) {
		src := "# Title\nSubtitle\nTheme:\n\n" + body
		if _, err := Parse(strings.NewReader(src), "deck.slide", 0); err == nil {
			t.Fatal("expected an error for an empty Theme header")
		}
	})

	t.Run("TitlesOnly skips resolution", func(t *testing.T) {
		src := "# Title\nSubtitle\nTheme: nope\n\n" + body
		doc, err := Parse(strings.NewReader(src), "deck.slide", TitlesOnly)
		if err != nil {
			t.Fatalf("TitlesOnly should not resolve the theme: %v", err)
		}
		if doc.Theme != "nope" {
			t.Errorf("Theme = %q, want %q", doc.Theme, "nope")
		}
		if doc.ThemeHref != "" {
			t.Errorf("ThemeHref = %q, want empty in TitlesOnly mode", doc.ThemeHref)
		}
	})

	t.Run("empty background errors", func(t *testing.T) {
		src := "# Title\nSubtitle\nBackground:\n\n" + body
		if _, err := Parse(strings.NewReader(src), "deck.slide", 0); err == nil {
			t.Fatal("expected an error for an empty Background header")
		}
	})
}
