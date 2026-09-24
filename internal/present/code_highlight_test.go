// Copyright 2026 The present Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package present

import (
	"strings"
	"testing"
)

func TestCodeSyntaxHighlighting(t *testing.T) {
	defer func(play bool) { PlayEnabled = play }(PlayEnabled)
	PlayEnabled = true

	src := []byte("package main\n\nfunc main() { // HL\n\tprintln(\"hi\") // HL\n}\n")
	read := func(string) ([]byte, error) { return src, nil }

	t.Run("code is highlighted", func(t *testing.T) {
		e, err := parseCode(&Context{read}, "main.go", 0, ".code main.go")
		if err != nil {
			t.Fatal(err)
		}
		html := strings.TrimSpace(string(e.(Code).Text))
		if !strings.Contains(html, "chroma-") {
			t.Errorf("expected chroma token classes, got:\n%s", html)
		}
		// Line numbers must still be emitted.
		if !strings.Contains(html, `<span num="4">`) {
			t.Errorf("expected line-number span num=4, got:\n%s", html)
		}
		// The highlighted lines are wrapped in <b>.
		if !strings.Contains(html, "<b>") {
			t.Errorf("expected highlighted lines wrapped in <b>, got:\n%s", html)
		}
	})

	t.Run("play blocks stay plain", func(t *testing.T) {
		e, err := parseCode(&Context{read}, "main.go", 0, ".play main.go")
		if err != nil {
			t.Fatal(err)
		}
		html := string(e.(Code).Text)
		if strings.Contains(html, "chroma-") {
			t.Errorf("play blocks must not be tokenized, got:\n%s", html)
		}
	})

	t.Run("editable blocks stay plain", func(t *testing.T) {
		e, err := parseCode(&Context{read}, "main.go", 0, ".code -edit main.go")
		if err != nil {
			t.Fatal(err)
		}
		html := string(e.(Code).Text)
		if strings.Contains(html, "chroma-") {
			t.Errorf("editable blocks must not be tokenized, got:\n%s", html)
		}
	})
}
