// Copyright 2026 The present Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ppreeper/present/internal/present"
)

func render(t *testing.T, src, name, ext string) string {
	t.Helper()
	if err := initTemplates(embedFS); err != nil {
		t.Fatalf("initTemplates: %v", err)
	}
	doc, err := present.Parse(strings.NewReader(src), name, 0)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var buf bytes.Buffer
	if err := doc.Render(&buf, contentTemplate[ext]); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return buf.String()
}

// assertStylesheets checks that the given stylesheet hrefs appear in the
// document, in the specified order.
func assertStylesheets(t *testing.T, html string, hrefs ...string) {
	t.Helper()
	last := -1
	for _, href := range hrefs {
		i := strings.Index(html, href)
		if i < 0 {
			t.Fatalf("stylesheet %q not found in rendered document", href)
		}
		if i < last {
			t.Fatalf("stylesheet %q appears out of order", href)
		}
		last = i
	}
}

func TestRenderStylesheets(t *testing.T) {
	const body = "\nName\n\n## Section\n\nhi\n"

	tests := []struct {
		name  string
		ext   string
		src   string
		wants []string
	}{
		{
			name: "slides unthemed",
			ext:  ".slide",
			src:  "# Title\nSubtitle\n\n" + body,
			wants: []string{
				"/static/fonts.css",
				"/static/base.css",
				"/static/styles.css",
			},
		},
		{
			name: "slides themed",
			ext:  ".slide",
			src:  "# Title\nSubtitle\nTheme: dark\n\n" + body,
			wants: []string{
				"/static/fonts.css",
				"/static/base.css",
				"/static/styles.css",
				"/static/themes/dark.css",
			},
		},
		{
			name: "article themed",
			ext:  ".article",
			src:  "# Title\nSubtitle\nTheme: mono\n\n" + body,
			wants: []string{
				"/static/fonts.css",
				"/static/base.css",
				"/static/article.css",
				"/static/themes/mono.css",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := render(t, tt.src, "deck"+tt.ext, tt.ext)
			assertStylesheets(t, html, tt.wants...)
			if tt.name == "slides unthemed" && strings.Contains(html, "/static/themes/") {
				t.Error("unthemed deck must not link a theme stylesheet")
			}
		})
	}
}

func TestRenderBackground(t *testing.T) {
	const body = "\nName\n\n## Section\n\nhi\n"
	src := "# Title\nSubtitle\nBackground: https://example.com/bg.png\n\n" + body
	html := render(t, src, "deck.slide", ".slide")
	if !strings.Contains(html, "--slide-background-image:url") {
		t.Errorf("expected deck background custom property, got:\n%s", html)
	}
	if !strings.Contains(html, "https://example.com/bg.png") {
		t.Errorf("expected background URL in rendered output")
	}
}

func TestRenderCustomTheme(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "custom.css"), []byte("/* x */"), 0644); err != nil {
		t.Fatal(err)
	}
	const body = "\nName\n\n## Section\n\nhi\n"
	src := "# Title\nSubtitle\nTheme: custom.css\n\n" + body
	html := render(t, src, filepath.Join(dir, "deck.slide"), ".slide")
	assertStylesheets(t, html,
		"/static/fonts.css",
		"/static/base.css",
		"/static/styles.css",
		"custom.css",
	)
}
