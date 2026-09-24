// Copyright 2026 The present Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"
)

func TestHowtoText(t *testing.T) {
	for _, section := range []string{"HEADER", "THEMES", "FONTS", "COMMAND INVOCATIONS", "KEYBOARD SHORTCUTS"} {
		if !strings.Contains(howtoText, section) {
			t.Errorf("howtoText is missing section %q", section)
		}
	}
	for _, theme := range []string{"default", "dark", "mono", "high-contrast"} {
		if !strings.Contains(howtoText, theme) {
			t.Errorf("howtoText does not mention bundled theme %q", theme)
		}
	}
}
