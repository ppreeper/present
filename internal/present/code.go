// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package present

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"html/template"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

// PlayEnabled specifies whether runnable playground snippets should be
// displayed in the present user interface.
var PlayEnabled = false

// TODO(adg): replace the PlayEnabled flag with something less spaghetti-like.
// Instead this will probably be determined by a template execution Context
// value that contains various global metadata required when rendering
// templates.

// NotesEnabled specifies whether presenter notes should be displayed in the
// present user interface.
var NotesEnabled = false

func init() {
	Register("code", parseCode)
	Register("play", parseCode)
}

type Code struct {
	Cmd      string // original command from present source
	Text     template.HTML
	Play     bool   // runnable code
	Edit     bool   // editable code
	FileName string // file name
	Ext      string // file extension
	Raw      []byte // content of the file
}

func (c Code) PresentCmd() string   { return c.Cmd }
func (c Code) TemplateName() string { return "code" }

// The input line is a .code or .play entry with a file name and an optional HLfoo marker on the end.
// Anything between the file and HL (if any) is an address expression, which we treat as a string here.
// We pick off the HL first, for easy parsing.
var (
	highlightRE = regexp.MustCompile(`\s+HL([a-zA-Z0-9_]+)?$`)
	hlCommentRE = regexp.MustCompile(`(.+) // HL(.*)$`)
	codeRE      = regexp.MustCompile(`\.(code|play)\s+((?:(?:-edit|-numbers)\s+)*)([^\s]+)(?:\s+(.*))?$`)
)

// parseCode parses a code present directive. Its syntax:
//
//	.code [-numbers] [-edit] <filename> [address] [highlight]
//
// The directive may also be ".play" if the snippet is executable.
func parseCode(ctx *Context, sourceFile string, sourceLine int, cmd string) (Elem, error) {
	cmd = strings.TrimSpace(cmd)
	origCmd := cmd

	// Pull off the HL, if any, from the end of the input line.
	highlight := ""
	if hl := highlightRE.FindStringSubmatchIndex(cmd); len(hl) == 4 {
		if hl[2] < 0 || hl[3] < 0 {
			return nil, fmt.Errorf("%s:%d invalid highlight syntax", sourceFile, sourceLine)
		}
		highlight = cmd[hl[2]:hl[3]]
		cmd = cmd[:hl[2]-2]
	}

	// Parse the remaining command line.
	// Arguments:
	// args[0]: whole match
	// args[1]:  .code/.play
	// args[2]: flags ("-edit -numbers")
	// args[3]: file name
	// args[4]: optional address
	args := codeRE.FindStringSubmatch(cmd)
	if len(args) != 5 {
		return nil, fmt.Errorf("%s:%d: syntax error for .code/.play invocation", sourceFile, sourceLine)
	}
	command, flags, file, addr := args[1], args[2], args[3], strings.TrimSpace(args[4])
	play := command == "play" && PlayEnabled

	// Read in code file and (optionally) match address.
	filename := filepath.Join(filepath.Dir(sourceFile), file)
	textBytes, err := ctx.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s:%d: %v", sourceFile, sourceLine, err)
	}
	lo, hi, err := addrToByteRange(addr, 0, textBytes)
	if err != nil {
		return nil, fmt.Errorf("%s:%d: %v", sourceFile, sourceLine, err)
	}
	if lo > hi {
		// The search in addrToByteRange can wrap around so we might
		// end up with the range ending before its starting point
		hi, lo = lo, hi
	}

	// Acme pattern matches can stop mid-line,
	// so run to end of line in both directions if not at line start/end.
	for lo > 0 && textBytes[lo-1] != '\n' {
		lo--
	}
	if hi > 0 {
		for hi < len(textBytes) && textBytes[hi-1] != '\n' {
			hi++
		}
	}

	lines := codeLines(textBytes, lo, hi)

	data := &codeTemplateData{
		Lines:   formatLines(lines, highlight),
		Edit:    strings.Contains(flags, "-edit"),
		Numbers: strings.Contains(flags, "-numbers"),
	}

	// Syntax-highlight non-interactive code blocks. Runnable and editable
	// blocks stay as plain text so the browser can edit them safely.
	if !play && !data.Edit {
		data.Lines = highlightLines(data.Lines, filename)
	}

	// Include before and after in a hidden span for playground code.
	if play {
		data.Prefix = textBytes[:lo]
		data.Suffix = textBytes[hi:]
	}

	var buf bytes.Buffer
	if err := codeTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}
	return Code{
		Cmd:      origCmd,
		Text:     template.HTML(buf.String()),
		Play:     play,
		Edit:     data.Edit,
		FileName: filepath.Base(filename),
		Ext:      filepath.Ext(filename),
		Raw:      rawCode(lines),
	}, nil
}

// formatLines returns a new slice of codeLine with the given lines
// replacing tabs with spaces and adding highlighting where needed.
func formatLines(lines []codeLine, highlight string) []codeLine {
	formatted := make([]codeLine, len(lines))
	for i, line := range lines {
		// Replace tabs with spaces, which work better in HTML.
		line.L = strings.Replace(line.L, "\t", "    ", -1)

		// Highlight lines that end with "// HL[highlight]"
		// and strip the magic comment.
		if m := hlCommentRE.FindStringSubmatch(line.L); m != nil {
			line.L = m[1]
			line.HL = m[2] == highlight
		}

		formatted[i] = line
	}
	return formatted
}

// rawCode returns the code represented by the given codeLines without any kind
// of formatting.
func rawCode(lines []codeLine) []byte {
	b := new(bytes.Buffer)
	for _, line := range lines {
		b.WriteString(line.L)
		b.WriteByte('\n')
	}
	return b.Bytes()
}

// highlightLines syntax-highlights each line of code, using the lexer matched
// by filename. Unknown file types fall back to plain text. The returned lines
// carry their highlighted HTML in the HTML field; L is left untouched.
func highlightLines(lines []codeLine, filename string) []codeLine {
	if len(lines) == 0 {
		return lines
	}

	lexer := lexers.Match(filepath.Base(filename))
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	var src strings.Builder
	for i, l := range lines {
		if i > 0 {
			src.WriteByte('\n')
		}
		src.WriteString(l.L)
	}

	tokens, err := chroma.Tokenise(lexer, nil, src.String())
	if err != nil {
		return lines
	}

	split := chroma.SplitTokensIntoLines(tokens)
	for i := range lines {
		if i >= len(split) {
			break
		}
		line := split[i]
		// SplitTokensIntoLines keeps the line's trailing newline on the last
		// token; the template emits its own newline between line spans.
		if n := len(line); n > 0 {
			line[n-1].Value = strings.TrimRight(line[n-1].Value, "\r\n")
		}
		h := renderTokens(line)
		if lines[i].HL {
			h = template.HTML("<b>" + string(h) + "</b>")
		}
		lines[i].HTML = h
	}
	return lines
}

// renderTokens renders chroma tokens as HTML spans carrying "chroma-" class
// names, HTML-escaping the token text.
func renderTokens(tokens []chroma.Token) template.HTML {
	var b strings.Builder
	for _, t := range tokens {
		if t.Value == "" {
			continue
		}
		cls := chromaClass(t.Type)
		if cls == "" {
			b.WriteString(html.EscapeString(t.Value))
			continue
		}
		b.WriteString(`<span class="`)
		b.WriteString(cls)
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(t.Value))
		b.WriteString(`</span>`)
	}
	return template.HTML(b.String())
}

// chromaClass returns the prefixed CSS class for a token type, walking up to
// the nearest parent type that has one. It mirrors the class naming used by
// chroma's own HTML formatter.
func chromaClass(t chroma.TokenType) string {
	for t != 0 {
		if cls, ok := chroma.StandardTypes[t]; ok {
			if cls == "" {
				return ""
			}
			return "chroma-" + cls
		}
		t = t.Parent()
	}
	return ""
}

type codeTemplateData struct {
	Lines          []codeLine
	Prefix, Suffix []byte
	Edit, Numbers  bool
}

var leadingSpaceRE = regexp.MustCompile(`^[ \t]*`)

var codeTemplate = template.Must(template.New("code").Funcs(template.FuncMap{
	"trimSpace":    strings.TrimSpace,
	"leadingSpace": leadingSpaceRE.FindString,
}).Parse(codeTemplateHTML))

const codeTemplateHTML = `
{{with .Prefix}}<pre style="display: none"><span>{{printf "%s" .}}</span></pre>{{end -}}

<pre{{if .Edit}} contenteditable="true" spellcheck="false"{{end}}{{if .Numbers}} class="numbers"{{end}}>{{/*
	*/}}{{range .Lines}}<span num="{{.N}}">{{/*
	*/}}{{if .HTML}}{{.HTML}}{{/*
	*/}}{{else if .HL}}<b>{{.L}}</b>{{/*
	*/}}{{else}}{{.L}}{{end}}{{/*
*/}}</span>
{{end}}</pre>
{{with .Suffix}}<pre style="display: none"><span>{{printf "%s" .}}</span></pre>{{end -}}
`

// codeLine represents a line of code extracted from a source file.
type codeLine struct {
	L    string        // The line of code, as plain text.
	N    int           // The line number from the source file.
	HL   bool          // Whether the line should be highlighted.
	HTML template.HTML // Syntax-highlighted HTML for the line, if any.
}

// codeLines takes a source file and returns the lines that
// span the byte range specified by start and end.
// It discards lines that end in "OMIT".
func codeLines(src []byte, start, end int) (lines []codeLine) {
	startLine := 1
	for i, b := range src {
		if i == start {
			break
		}
		if b == '\n' {
			startLine++
		}
	}
	s := bufio.NewScanner(bytes.NewReader(src[start:end]))
	for n := startLine; s.Scan(); n++ {
		l := s.Text()
		if strings.HasSuffix(l, "OMIT") {
			continue
		}
		lines = append(lines, codeLine{L: l, N: n})
	}
	// Trim leading and trailing blank lines.
	for len(lines) > 0 && len(lines[0].L) == 0 {
		lines = lines[1:]
	}
	for len(lines) > 0 && len(lines[len(lines)-1].L) == 0 {
		lines = lines[:len(lines)-1]
	}
	return
}

func parseArgs(name string, line int, args []string) (res []any, err error) {
	res = make([]any, len(args))
	for i, v := range args {
		if len(v) == 0 {
			return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
		}
		switch v[0] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
			}
			res[i] = n
		case '/':
			if len(v) < 2 || v[len(v)-1] != '/' {
				return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
			}
			res[i] = v
		case '$':
			res[i] = "$"
		case '_':
			if len(v) == 1 {
				// Do nothing; "_" indicates an intentionally empty parameter.
				break
			}
			fallthrough
		default:
			return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
		}
	}
	return
}
