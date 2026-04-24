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

## File Format

Present files begin with a header:

```
# Title of document
Subtitle of document
15:04 2 Jan 2006
Tags: foo, bar, baz
Summary: A short summary.
```

Followed by optional author blocks, then slides/sections starting with `##` (Markdown) or `*` (legacy).

See [doc.go](internal/present/doc.go) for the full format specification, including commands like `.code`, `.play`, `.image`, `.video`, `.iframe`, `.link`, `.html`, `.caption`, and `.background`.
