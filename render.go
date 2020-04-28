package nbtohtml

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/alecthomas/chroma"
	htmlFormatter "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/buildkite/terminal-to-html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

// Enable syntax highlighting in Markdown
var markdownParser = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				htmlFormatter.WithClasses(true),
			),
		),
	),
)

// renderMarkdown uses the goldmark library to convert the provided Markdown lines to HTML.
func renderMarkdown(markdown string) template.HTML {
	var htmlBuffer bytes.Buffer
	if err := markdownParser.Convert([]byte(markdown), &htmlBuffer); err != nil {
		panic(err)
	}
	// goldmark does not render raw HTML or potentially-dangerous URLs, so HTML should be safe from
	// code injection
	return template.HTML(htmlBuffer.String()) // nolint:gosec
}

// renderSourceCode uses the Chroma library to convert the provided source code string to HTML.
// Instead of inline styles, HTML classes are used for syntax highlighting, which allows the users
// to style source code according to their needs.
func renderSourceCode(source string, languageID string) (template.HTML, error) {
	htmlBuffer := new(bytes.Buffer)

	// Set up lexer for programming language
	var l chroma.Lexer
	if languageID != "" {
		l = lexers.Get(languageID)
	}
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Configure Chroma to use classes instead of inline styles
	formatter := htmlFormatter.New(htmlFormatter.WithClasses(true))

	iterator, err := l.Tokenise(nil, source)
	if err != nil {
		return "", fmt.Errorf("could not render source code (tokenization error): %d", err)
	}

	err = formatter.Format(htmlBuffer, styles.GitHub, iterator)
	if err != nil {
		return "", fmt.Errorf("could not render source code (formatting error): %d", err)
	}

	// Chroma escapes tags, so HTML should be safe from code injection
	return template.HTML(htmlBuffer.String()), nil // nolint:gosec
}

// renderMarkdown uses the `terminal-to-html` library to convert the provided Terminal output to
// HTML with classes for styling ANSI colors.
func renderTerminalOutput(source []string) template.HTML {
	linesHTML := []string{}
	for _, tracebackLine := range source {
		lineHTML := terminal.Render([]byte(tracebackLine))
		linesHTML = append(linesHTML, string(lineHTML))
	}
	htmlString := fmt.Sprintf("<pre>%s</pre>", strings.Join(linesHTML, "\n"))
	// `terminal-to-html` escapes tags, so HTML should be safe from code injection
	return template.HTML(htmlString) // nolint:gosec
}
