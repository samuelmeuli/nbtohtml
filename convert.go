package nbtohtml

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday/v2"
	"io"
	"io/ioutil"
	"strings"
)

// highlightCode uses Chroma to convert the provided source code string to HTML with tags and
// classes for syntax highlighting.
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
	formatter := html.New(html.WithClasses(true))

	iterator, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return formatter.Format(writer, styles.GitHub, iterator)
}

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

// convertMarkdownCellToHTML converts a Markdown cell to HTML with Blackfriday
func convertMarkdownCellToHTML(cell Cell) string {
	md := strings.Join(cell.Source, "")
	mdHTML := blackfriday.Run([]byte(md))
	return fmt.Sprintf("<div class=\"cell markdown-cell\">%s</div>", string(mdHTML))
}

// convertCodeCellToHTML converts a code cell to HTML with classes for syntax highlighting using Chroma
func convertCodeCellToHTML(cell Cell, fileExtension string) (string, error) {
	codeString := strings.Join(cell.Source, "")
	codeBuffer := new(bytes.Buffer)
	err := highlightCode(codeBuffer, codeString, fileExtension)
	return fmt.Sprintf("<div class=\"cell code-cell\">%s</div>", codeBuffer.String()), err
}

// convertRawCellToHTML returns a simple HTML element for the raw notebook cell
func convertRawCellToHTML(cell Cell) string {
	return fmt.Sprintf(
		"<div class=\"cell raw-cell\"><pre>%s</pre></div>",
		strings.Join(cell.Source, ""),
	)
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
	htmlString := "<div class=\"jupyter-notebook\">"
	for _, cell := range notebook.Cells {
		switch cell.CellType {
		case "markdown":
			htmlString += convertMarkdownCellToHTML(cell)
		case "code":
			codeHTMLString, err := convertCodeCellToHTML(cell, fileExtension)
			if err == nil {
				htmlString += codeHTMLString
			} else {
				fmt.Printf("skipping cell (syntax highlighting error: %d)", err)
			}
		case "raw":
			htmlString += convertRawCellToHTML(cell)
		default:
			fmt.Printf("skipping cell (unrecognized cell type \"%s\")", cell.CellType)
		}
	}

	return htmlString + "</div>", nil
}
