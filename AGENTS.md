# AGENTS.md

Guidance for agents working in this repository.

## Overview

This is a fork of `golang.org/x/tools/cmd/present` (module `github.com/ppreeper/present`). It runs an HTTP server that renders `.slide` presentations and `.article` documents, plus a live code playground.

## Build, Test, Lint

```sh
go build ./...                 # build everything
go test ./...                  # run all tests
go test ./internal/present/    # parser, rendering, code highlighting, theme/background resolution
go test ./cmd/present/         # template rendering (stylesheet order, background, themes)
go vet ./...
gofmt -l .                     # list unformatted files
```

There is no separate lint config; use `go vet` and `gofmt`.

## Architecture

- `cmd/present/` — the server (`package main`).
  - `main.go` — flags, embeds `static` and `templates` via `//go:embed`, serves `/static/`, starts the server. With `-base` it serves those assets from disk instead of the embedded copies.
  - `dir.go` — request routing for `/`: renders docs, lists directories, or serves arbitrary files out of the content directory (so CSS/fonts placed next to a deck are HTTP-servable). `initTemplates` wires `templates/action.tmpl` plus the format template into `contentTemplate`.
  - `templates/` — `slides.tmpl`, `article.tmpl`, `dir.tmpl`, and shared `action.tmpl` element templates.
  - `static/` — CSS, JavaScript, bundled themes, and self-hosted fonts.
- `internal/present/` — parsing and rendering.
  - `parse.go` — parses headers, author blocks, sections; resolves `Theme:`/`Background:`.
  - `theme.go` — theme/background name and path resolution.
  - `code.go` — `.code`/`.play` parsing and server-side Chroma syntax highlighting.
  - `doc.go` — the authoritative file-format specification (keep it in sync with parser changes).
- `internal/txtar`, `internal/playground/` — supporting packages from upstream.

## Styling and Theming

- CSS is layered: `fonts.css` (self-hosted `@font-face`) → `base.css` (design tokens + reset) → `styles.css` or `article.css` → optional theme.
- `base.css` `:root` is the published token contract. **Adding** tokens is fine; **renaming or removing** them is a breaking change for custom themes. Document new tokens in the theme section of `README.md`.
- Bundled themes live in `cmd/present/static/themes/` and override tokens only. `default.css` is intentionally a no-op because the base layer is the default look.
- Syntax highlighting classes are `chroma-*`; their colors are `--chroma-*` tokens in `base.css`.
- Fonts are regenerated with `tools/fetch-fonts.sh` (downloads WOFF2 subsets + licences from Fontsource). Do not hand-edit `fonts.css` or the files under `static/fonts/`.

## Conventions

- Match the surrounding Go style; keep code gofmt-clean and `go vet`-clean.
- Do not add code comments unless they clarify non-obvious behavior.
- User-visible format changes must be reflected in `internal/present/doc.go` and `README.md`.

## Manual Verification

- The `F` "fill the window" toggle is JavaScript-only and has no automated test. Verify it manually in a browser with slides at a few different window sizes, and confirm the state persists across reloads (localStorage key `present.fillScreen`).
- Verify theme changes by rendering a deck with each bundled theme and in both light and print/PDF modes.
