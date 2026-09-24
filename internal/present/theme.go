// Copyright 2026 The present Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package present

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// bundledThemes maps the case-insensitive name of a theme that ships with the
// present binary to the URL of its stylesheet.
var bundledThemes = map[string]string{
	"default":       "/static/themes/default.css",
	"dark":          "/static/themes/dark.css",
	"mono":          "/static/themes/mono.css",
	"high-contrast": "/static/themes/high-contrast.css",
}

// bundledThemeNames returns the recognized bundled theme names in a stable
// order, for use in error messages and documentation.
func bundledThemeNames() []string {
	return []string{"default", "dark", "mono", "high-contrast"}
}

// isRemoteOrRooted reports whether v is a URL or a server-absolute path that
// should be emitted verbatim rather than resolved against the source file.
func isRemoteOrRooted(v string) bool {
	return strings.HasPrefix(v, "/") ||
		strings.HasPrefix(v, "//") ||
		strings.Contains(v, "://")
}

// resolveTheme maps a Theme header value to the URL of the stylesheet to load.
// A bundled theme name (matched case-insensitively) maps to its stylesheet;
// anything ending in ".css" or containing a path separator is treated as a
// stylesheet path relative to the document's directory and checked for
// existence. name is the path of the document being parsed.
func resolveTheme(value, name string) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", errors.New("empty theme value")
	}
	if href, ok := bundledThemes[strings.ToLower(v)]; ok {
		return href, nil
	}
	if strings.HasSuffix(strings.ToLower(v), ".css") || strings.ContainsAny(v, "/\\") {
		if !isRemoteOrRooted(v) {
			p := filepath.Join(filepath.Dir(name), v)
			if _, err := os.Stat(p); err != nil {
				return "", fmt.Errorf("theme %q: %v", v, err)
			}
		}
		return v, nil
	}
	return "", fmt.Errorf("unknown theme %q (built-in themes: %s)",
		v, strings.Join(bundledThemeNames(), ", "))
}

// resolveBackground maps a Background header value to a CSS declaration that
// sets the deck-wide slide background image. A relative path is checked for
// existence against the document's directory; URLs and absolute paths are
// emitted verbatim. name is the path of the document being parsed.
func resolveBackground(value, name string) (template.CSS, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", errors.New("empty background value")
	}
	if !isRemoteOrRooted(v) {
		p := filepath.Join(filepath.Dir(name), v)
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf("background %q: %v", v, err)
		}
	}
	return template.CSS("--slide-background-image:url(" + cssURL(v) + ")"), nil
}

// cssURL formats v as a quoted CSS url() token.
func cssURL(v string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range v {
		switch r {
		case '"', '\\':
			b.WriteByte('\\')
			b.WriteRune(r)
		case '\n', '\r':
			b.WriteByte(' ')
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
