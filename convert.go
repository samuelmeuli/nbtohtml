package nbtohtml

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma"
	htmlFormatter "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/buildkite/terminal-to-html"
	"github.com/russross/blackfriday/v2"
	"html"
	"io"
	"io/ioutil"
	"strings"
)

// 3rd party renderers

// highlightCode uses the Chroma library to convert the provided source code string to HTML. Instead
// of inline styles, HTML classes are used for syntax highlighting, which allows the users to style
// source code according to their needs.
func highlightCode(writer io.Writer, source string, lexer string) error {
	l := lexers.Get(lexer)
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
		return err
	}
	return formatter.Format(writer, styles.GitHub, iterator)
}

// renderMarkdown uses the Blackfriday library to convert the provided Markdown lines to HTML.
func renderMarkdown(markdownLines []string) string {
	markdownString := strings.Join(markdownLines, "")
	mdHTML := blackfriday.Run([]byte(markdownString))
	return string(mdHTML)
}

// Output renderers

func convertDataOutputToHTML(output Output) (string, error) {
	if output.Data.TextHTML != nil {
		return fmt.Sprintf(
			`<div class="output output-data-html">%s</div>`,
			strings.Join(output.Data.TextHTML, ""),
		), nil
	}
	if output.Data.ApplicationPDF != nil {
		return "", fmt.Errorf("missing conversion logic for `application/pdf` data type")
	}
	if output.Data.TextLaTeX != nil {
		return "", fmt.Errorf("missing conversion logic for `text/latex` data type")
	}
	if output.Data.ImageSVGXML != nil {
		return fmt.Sprintf(
			`<div class="output output-data-svg">%s</div>`,
			strings.Join(output.Data.ImageSVGXML, ""),
		), nil
	}
	if output.Data.ImagePNG != nil {
		return fmt.Sprintf(
			`<div class="output output-data-png"><img src="data:image/png;base64,%s"></div>`,
			*output.Data.ImagePNG,
		), nil
	}
	if output.Data.ImageJPEG != nil {
		return fmt.Sprintf(
			`<div class="output output-data-jpeg"><img src="data:image/jpeg;base64,%s"></div>`,
			*output.Data.ImageJPEG,
		), nil
	}
	if output.Data.TextMarkdown != nil {
		return fmt.Sprintf(
			`<div class="output output-data-markdown">%s</div>`,
			renderMarkdown(output.Data.TextMarkdown),
		), nil
	}
	if output.Data.TextPlain != nil {
		return fmt.Sprintf(
			`<div class="output output-data-plain-text"><pre>%s</pre></div>`,
			html.EscapeString(strings.Join(output.Data.TextPlain, "")),
		), nil
	}

	return "", fmt.Errorf(
		"missing `execute_result` data type in output of type `%s`",
		output.OutputType,
	)
}

func convertErrorOutputToHTML(output Output) (string, error) {
	if output.Traceback == nil {
		return "", fmt.Errorf("missing `traceback` key in output of type `error`")
	}

	// Convert ANSI colors to HTML
	var linesHTML []string
	for _, tracebackLine := range output.Traceback {
		lineHTML := terminal.Render([]byte(tracebackLine))
		linesHTML = append(linesHTML, string(lineHTML))
	}

	return fmt.Sprintf(
		`<div class="output output-error"><pre>%s</pre></div>`,
		strings.Join(linesHTML, "\n"),
	), nil
}

func convertStreamOutputToHTML(output Output) (string, error) {
	if output.Text == nil {
		return "", fmt.Errorf("missing `text` key in output of type `stream`")
	}

	return fmt.Sprintf(
		`<div class="output output-stream"><pre>%s</pre></div>`,
		strings.Join(output.Text, ""),
	), nil
}

func convertOutputToHTML(output Output) (string, error) {
	switch output.OutputType {
	case "display_data":
		return convertDataOutputToHTML(output)
	case "error":
		return convertErrorOutputToHTML(output)
	case "execute_result":
		return convertDataOutputToHTML(output)
	case "stream":
		return convertStreamOutputToHTML(output)
	default:
		return "", fmt.Errorf("missing conversion logic for output type `%s`", output.OutputType)
	}
}

// Cell renderers

// convertMarkdownCellToHTML converts a Markdown cell to HTML.
func convertMarkdownCellToHTML(cell Cell) string {
	return fmt.Sprintf(`<div class="cell cell-markdown">%s</div>`, renderMarkdown(cell.Source))
}

// convertCodeCellToHTML converts a code cell to HTML with classes for syntax highlighting.
func convertCodeCellToHTML(cell Cell, fileExtension string) (string, error) {
	cellHTML := `<div class="cell cell-code">`

	// Source
	codeString := strings.Join(cell.Source, "")
	codeBuffer := new(bytes.Buffer)
	err := highlightCode(codeBuffer, codeString, fileExtension)
	if err != nil {
		return "", err
	}
	cellHTML += fmt.Sprintf(`<div class="input">%s</div>`, codeBuffer.String())

	// Outputs
	if cell.Outputs != nil {
		for _, output := range cell.Outputs {
			var outputHTML, err = convertOutputToHTML(output)
			if err == nil {
				cellHTML += outputHTML
			} else {
				fmt.Printf("skipping output: %s\n", err)
			}
		}
	}

	return cellHTML + "</div>", nil
}

// convertRawCellToHTML returns a simple HTML element for the raw notebook cell.
func convertRawCellToHTML(cell Cell) string {
	return fmt.Sprintf(
		`<div class="cell cell-raw"><pre>%s</pre></div>`,
		html.EscapeString(strings.Join(cell.Source, "")),
	)
}

// Notebook renderers

// ConvertFileToHTML reads the file at the provided path and converts its content (the Jupyter
// Notebook JSON) to HTML.
func ConvertFileToHTML(notebookPath string) (string, error) {
	// Read file
	fileContent, err := ioutil.ReadFile(notebookPath)
	if err != nil {
		return "", fmt.Errorf("could not read Jupyter Notebook file at %s", notebookPath)
	}

	// Convert file content
	return ConvertStringToHTML(string(fileContent))
}

// ConvertStringToHTML converts the provided Jupyter Notebook JSON string to HTML.
func ConvertStringToHTML(notebookString string) (string, error) {
	notebook, err := parseNotebook(notebookString)
	if err != nil {
		return "", err
	}

	// Get format extension of Jupyter Kernel language (e.g. "py")
	fileExtension := notebook.Metadata.LanguageInfo.FileExtension[1:]

	// Build HTML string from converted cells
	htmlString := `<div class="jupyter-notebook">`
	for _, cell := range notebook.Cells {
		switch cell.CellType {
		case "markdown":
			htmlString += convertMarkdownCellToHTML(cell)
		case "code":
			codeHTMLString, err := convertCodeCellToHTML(cell, fileExtension)
			if err == nil {
				htmlString += codeHTMLString
			} else {
				fmt.Printf("skipping cell (syntax highlighting error: %d)\n", err)
			}
		case "raw":
			htmlString += convertRawCellToHTML(cell)
		default:
			fmt.Printf("skipping cell (unrecognized cell type \"%s\")\n", cell.CellType)
		}
	}

	return htmlString + "</div>", nil
}
