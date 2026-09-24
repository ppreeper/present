# present

A tool for displaying slide presentations and articles. It runs a web server that serves slide and article files from the current directory.

## Install

```
go install github.com/ppreeper/present/cmd/present@latest
```

## Usage

```
present [flags] [directory]
```

Run `present` in a directory containing `.slide` or `.article` files and open your browser to the displayed URL (default `http://127.0.0.1:3999`).

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-http` | `127.0.0.1:3999` | HTTP service address |
| `-play` | `true` | Enable playground (run code snippets) |
| `-notes` | `false` | Enable presenter notes (press `N` in browser) |
| `-content` | `.` | Base path for presentation content |
| `-use_playground` | `false` | Use play.golang.org for code execution |
| `-howto` | `false` | Print a guide to the header, themes, fonts, and elements, then exit |

Run `present -howto` for a built-in guide covering the header, themes, fonts, commands, and keyboard shortcuts.

## File Format

Present files begin with a header:

```
# Title of document
Subtitle of document
15:04 2 Jan 2006
Tags: foo, bar, baz
Summary: A short summary.
Theme: dark
Background: images/cover.jpg
```

Followed by optional author blocks, then slides/sections starting with `##` (Markdown) or `*` (legacy).

See [doc.go](internal/present/doc.go) for the full format specification, including commands like `.code`, `.play`, `.image`, `.video`, `.iframe`, `.link`, `.html`, `.caption`, and `.background`.

## Themes

A deck can choose a stylesheet with a `Theme:` header line. Four themes are bundled:

| Theme | Description |
|-------|-------------|
| `default` | The built-in light look. Used automatically when no `Theme:` line is present. |
| `dark` | Dark background with light text. |
| `mono` | Grayscale, understated. |
| `high-contrast` | Maximum contrast with stronger borders, for projection and accessibility. |

`Theme:` names are case-insensitive. `Theme: default` is a no-op. Themes apply to both `.slide` and `.article` files. Print/PDF output is always forced to black-on-white so dark themes remain ink-friendly.

### Custom themes

`Theme:` may instead point at a CSS file:

```
Theme: my-theme.css
```

Relative paths are resolved against the directory containing the present file, and the file must exist. A path starting with `/` or a full URL is emitted verbatim (and not checked). The referenced CSS is linked after the bundled base and format stylesheets, so it can override any selector; the simplest custom themes just set the design-token custom properties defined in [`base.css`](cmd/present/static/base.css).

### Backgrounds

A `Background:` header line displays an image behind the content of every slide:

```
Background: images/cover.jpg
```

As with themes, relative paths are resolved against the present file's directory and must exist, while absolute/URL values are used verbatim. To override the deck-wide background on a single slide, use the existing `.background image.jpg` command on that slide.

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `→` / `Space` / `PgDn` | Next slide |
| `←` / `PgUp` | Previous slide |
| `F` | Fill the window with the slide (toggle; scales up or down to fit) |
| `N` | Toggle presenter notes (when started with `-notes`) |
| `H` / `Esc` | Toggle the help overlay |

## Fonts

The bundled themes use self-hosted open-source fonts (Inter, Source Sans 3, IBM Plex Sans, Source Serif 4, Merriweather, EB Garamond, JetBrains Mono, Source Code Pro, IBM Plex Mono) shipped as WOFF2 subsets with `unicode-range`, so only the scripts actually used are downloaded. All are licensed under the SIL Open Font License; licence files live alongside the fonts in [`cmd/present/static/fonts/`](cmd/present/static/fonts/). The downloading/subsetting helper is [`tools/fetch-fonts.sh`](tools/fetch-fonts.sh).

To use your own font, drop the font file into your content directory and reference it from a custom theme with an `@font-face` rule, then point `--font-sans`, `--font-serif`, or `--font-mono` at the family name.
