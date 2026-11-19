// Copyright 2026 The present Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const howtoText = `PRESENT - HOW TO WRITE A DECK

A .slide (presentation) or .article (document) file has three parts:

  1. A header (title and metadata)
  2. Zero or more author blocks
  3. The body: slides/sections containing text and command invocations


================================================================================
HEADER
================================================================================

The header is the first block of the file, terminated by a blank line:

  # Title of document
  Subtitle of document
  15:04 2 Jan 2006
  Tags: foo, bar, baz
  Summary: A short summary used in feeds.
  OldURL: former-path-for-this-doc
  Theme: dark
  Background: images/cover.jpg

Only the title is required. The leading "# " marks a Markdown-enabled file
(it uses Markdown for body text). Without it, the file uses legacy present
markup. The date may omit the time ("2 Jan 2006", interpreted as 10am UTC).

  Theme:      Stylesheet to use (see THEMES). Case-insensitive name, or a
              path/URL to a custom CSS file.
  Background: Image shown behind every slide (slides only). Relative paths
              are resolved against this file's directory; "/..." and full
              URLs are used as-is. Override per slide with ".background".


================================================================================
AUTHOR BLOCKS
================================================================================

After the header, before the body:

  Author Name
  Job title, Company
  joe@example.com
  https://url/
  @twitter_name

For slides, only plain text lines appear on the title slide. Separate multiple
authors with a blank line.


================================================================================
SECTIONS AND TEXT
================================================================================

Markdown-enabled (header starts with "# "):

  ## Slide or section title
  ### Subsection {#anchor}
  #### Sub-subsection

  Some text with _italic_, **bold**, and ` + "`program`" + ` font.

  - bullets
  - more bullets

      Indented text (one tab or four spaces) is a code block.

  Visit [the Go home page](https://golang.org/).

  // A line beginning with // is a comment and is ignored.
  : A line beginning with ": " is a presenter note.

Legacy (no leading "# "):

  * Slide or section title
  ** Subsection
  *** Sub-subsection

  _italic_, *bold*, and ` + "`program`" + ` font.
  Visit [[https://golang.org][the Go home page]].


================================================================================
COMMAND INVOCATIONS
================================================================================

Each command starts with a period at the beginning of the line:

  .code x.go /^func main/,/^}/     Include source from a file.
  .play y.go                       Like .code, but runnable in the browser.
  .image images/betsy.jpg 100 200  Image (optional height and width).
  .background images/susan.jpg     Background image for this slide only.
  .iframe https://foo              Embed a page.
  .link https://foo label          Hyperlink with optional label.
  .html file.html                  Include a file as raw, unescaped HTML.
  .caption _Gopher_ by [[...][Renee French]]   Figure caption.

.code and .play:

  .code [-numbers] [-edit] file.go [address] [HLword]

  -numbers   Show line numbers.
  -edit      Show an editable text area (not syntax-highlighted).
  address    A sam/ed-style range, e.g. /^func main/,/^}/ . If omitted the
             whole file is shown.
  HLword     Highlight lines whose comment ends with "// HLword". A plain
             "// HL" matches an invocation with no HLword.

  Lines ending in "OMIT" are removed before display, so you can bracket a
  snippet with "// START OMIT" / "// END OMIT" markers.

.image / .iframe: 1 or 3 args (file [height] [width]); use "_" to preserve
the aspect ratio, e.g. ".image images/janet.jpg _ 300".

.video: 2 or 4 args (file content-type [height] [width]),
e.g. ".video videos/evangeline.mp4 video/mp4 400 600".


================================================================================
THEMES
================================================================================

Select a bundled theme with the header line "Theme: <name>":

  default         Built-in light look; also used when no Theme line is present.
  dark            Dark background with light text.
  mono            Grayscale and understated.
  high-contrast   Strong borders and maximum contrast, for projectors.

Names are case-insensitive; "Theme: default" is a no-op. Themes apply to both
.slide and .article files. Print/PDF output is always forced to black-on-white.

Custom theme:

  Theme: my-theme.css

Relative paths resolve against this file's directory and must exist; a path
starting with "/" or a full URL is linked verbatim. Your CSS is loaded after
the built-in stylesheets, so it can override anything. The easiest approach is
to redefine the design-token custom properties declared in base.css, for
example:

  :root {
    --font-sans: "My Font", system-ui, sans-serif;
    --color-accent: #b03a2e;
    --card-shadow: 0 10px 30px rgba(0,0,0,.25);
  }

Token groups: fonts, slide typography, page/card colors, code block, tables,
slide background, and syntax-highlighting colors (--chroma-*). Adding tokens
is safe; renaming them breaks custom themes.


================================================================================
FONTS
================================================================================

Self-hosted, open-source families ship with present:

  Sans:  Inter, Source Sans 3, IBM Plex Sans
  Serif: Source Serif 4, Merriweather, EB Garamond
  Mono:  JetBrains Mono, Source Code Pro, IBM Plex Mono

They are packaged as WOFF2 subsets with unicode-range, so only the scripts a
deck actually uses are downloaded. All are SIL OFL 1.1; licences live in
static/fonts/. Defaults: Inter for body text, JetBrains Mono for code.

To use your own font, place the font file in your content directory and add an
@font-face rule in a custom theme, then point a font token at it:

  @font-face {
    font-family: "My Font";
    src: url("MyFont.woff2") format("woff2");
  }
  :root { --font-sans: "My Font", system-ui, sans-serif; }


================================================================================
PRESENTER NOTES
================================================================================

Lines beginning with ": " are presenter notes (Markdown and legacy). They are
hidden by default. Start the server with -notes and press N in the browser to
open a synced second window showing the notes.


================================================================================
KEYBOARD SHORTCUTS
================================================================================

  ->  Space  PgDn    Next slide
  <-  PgUp           Previous slide
  F                  Fill the window with the slide (toggle)
  N                  Toggle presenter notes (with -notes)
  H  Esc             Toggle the help overlay


================================================================================
EXAMPLE
================================================================================

  # My Talk
  A short subtitle
  2 Jan 2006
  Theme: dark
  Background: images/cover.jpg

  My Name
  me@example.com

  ## Hello

  Some text.

  ## Code

  .code -numbers hello.go /^func main/,/^}/ HL

See internal/present/doc.go for the complete, authoritative format reference.
`
