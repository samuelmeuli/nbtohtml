package nbtohtml

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"strings"

	"github.com/alecthomas/chroma"
	htmlFormatter "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/buildkite/terminal-to-html"
	"github.com/russross/blackfriday/v2"
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
	return string(blackfriday.Run([]byte(markdownString)))
}

// Output renderers

// convertDataOutput converts data output (e.g. a base64-encoded plot image) to HTML.
func convertDataOutput(output Output) template.HTML {
	var htmlString string

	if output.Data.TextHTML != nil {
		htmlString = strings.Join(output.Data.TextHTML, "")
	} else if output.Data.ApplicationPDF != nil {
		// TODO: Implement PDF conversion
		fmt.Printf("missing conversion logic for `application/pdf` data type\n")
		htmlString = "<pre>PDF output</pre>"
	} else if output.Data.TextLaTeX != nil {
		// TODO: Implement LaTeX conversion
		fmt.Printf("missing conversion logic for `text/latex` data type\n")
		htmlString = "<pre>LaTeX output</pre>"
	} else if output.Data.ImageSVGXML != nil {
		htmlString = strings.Join(output.Data.ImageSVGXML, "")
	} else if output.Data.ImagePNG != nil {
		htmlString = fmt.Sprintf(`<img src="data:image/png;base64,%s">`, *output.Data.ImagePNG)
	} else if output.Data.ImageJPEG != nil {
		htmlString = fmt.Sprintf(`<img src="data:image/jpeg;base64,%s">`, *output.Data.ImageJPEG)
	} else if output.Data.TextMarkdown != nil {
		htmlString = renderMarkdown(output.Data.TextMarkdown)
	} else if output.Data.TextPlain != nil {
		htmlString = fmt.Sprintf(
			`<pre>%s</pre>`,
			html.EscapeString(strings.Join(output.Data.TextPlain, "")),
		)
	} else {
		htmlString = ""
		fmt.Printf("missing `execute_result` data type in output of type `%s`\n", output.OutputType)
	}

	return template.HTML(htmlString)
}

// convertErrorOutput converts error output (e.g. generated by a Python exception) to HTML.
func convertErrorOutput(output Output) template.HTML {
	if output.Traceback == nil {
		fmt.Printf("missing `traceback` key in output of type `error`\n")
		return "<pre>An unknown error occurred</pre>"
	}

	// Convert ANSI colors to HTML
	var linesHTML []string
	for _, tracebackLine := range output.Traceback {
		lineHTML := terminal.Render([]byte(tracebackLine))
		linesHTML = append(linesHTML, string(lineHTML))
	}
	htmlString := fmt.Sprintf(`<pre>%s</pre>`, strings.Join(linesHTML, "\n"))
	return template.HTML(htmlString)
}

// convertStreamOutput converts stream output (e.g. stdout written by a Python program) to HTML.
func convertStreamOutput(output Output) template.HTML {
	if output.Text == nil {
		fmt.Printf("missing `text` key in output of type `stream`\n")
		return ""
	}

	htmlString := fmt.Sprintf(`<pre>%s</pre>`, strings.Join(output.Text, ""))
	return template.HTML(htmlString)
}

// Cell renderers

// convertMarkdownCell converts a Markdown cell to HTML.
func convertMarkdownCell(cell Cell) template.HTML {
	return template.HTML(renderMarkdown(cell.Source))
}

// convertCodeCell converts a code cell to HTML with classes for syntax highlighting.
func convertCodeCell(cell Cell, fileExtension string) template.HTML {
	codeString := strings.Join(cell.Source, "")
	codeBuffer := new(bytes.Buffer)
	err := highlightCode(codeBuffer, codeString, fileExtension)
	if err != nil {
		fmt.Printf("skipping syntax highlighting: %d\n", err)
		return template.HTML(fmt.Sprintf("<pre>%s</pre>", codeString))
	}
	return template.HTML(codeBuffer.String())
}

// convertRawCell returns a simple HTML element for the raw notebook cell.
func convertRawCell(cell Cell) template.HTML {
	htmlString := fmt.Sprintf(
		`<pre>%s</pre>`,
		html.EscapeString(strings.Join(cell.Source, "")),
	)
	return template.HTML(htmlString)
}

// Input/output renderers

// convertPrompt returns an HTML string which indicates the input/output's execution count.
func convertPrompt(executionCount *int) template.HTML {
	if executionCount == nil {
		return ""
	}
	return template.HTML(fmt.Sprintf("[%d]:", *executionCount))
}

// convertOutput converts the provided cell input to HTML.
func convertInput(fileExtension string, cell Cell) template.HTML {
	switch cell.CellType {
	case "markdown":
		return convertMarkdownCell(cell)
	case "code":
		return convertCodeCell(cell, fileExtension)
	case "raw":
		return convertRawCell(cell)
	default:
		fmt.Printf("skipping cell (unrecognized cell type \"%s\")\n", cell.CellType)
		return ""
	}
}

// convertOutput converts the provided output of a cell execution to HTML.
func convertOutput(output Output) template.HTML {
	switch output.OutputType {
	case "display_data":
		return convertDataOutput(output)
	case "error":
		return convertErrorOutput(output)
	case "execute_result":
		return convertDataOutput(output)
	case "stream":
		return convertStreamOutput(output)
	default:
		fmt.Printf("missing conversion logic for output type `%s`\n", output.OutputType)
		return ""
	}
}

// Notebook renderers

// ConvertFile reads the file at the provided path and converts its content (the Jupyter Notebook
// JSON) to HTML.
func ConvertFile(writer io.Writer, notebookPath string) error {
	// Read file
	fileContent, err := ioutil.ReadFile(notebookPath)
	if err != nil {
		return fmt.Errorf("could not read Jupyter Notebook file at %s", notebookPath)
	}

	// Convert file content
	return ConvertString(writer, string(fileContent))
}

// ConvertString converts the provided Jupyter Notebook JSON string to HTML.
func ConvertString(writer io.Writer, notebookString string) error {
	notebook, err := parseNotebook(notebookString)
	if err != nil {
		return err
	}

	// Get format extension of Jupyter Kernel language (e.g. "py")
	fileExtension := notebook.Metadata.LanguageInfo.FileExtension[1:]

	t := template.New("notebook")
	t = t.Funcs(template.FuncMap{
		"convertPrompt": convertPrompt,
		"convertInput":  convertInput,
		"convertOutput": convertOutput,
	})
	t, err = t.Parse(`
		<div class="notebook">
			{{ $fileExtension := .FileExtension }}
			{{ range .Notebook.Cells }}
				<div class="cell cell-{{ .CellType }}">
					<div class="input-wrapper">
						<div class="input-prompt">
							{{ .ExecutionCount | convertPrompt }}
						</div>
						<div class="input">
							{{ . | convertInput $fileExtension }}
						</div>
					</div>
					{{ range .Outputs }}
						<div class="output-wrapper">
							<div class="output-prompt">
								{{ .ExecutionCount | convertPrompt }}
							</div>
							<div class="output output-{{ .OutputType }}">
								{{ . | convertOutput }}
							</div>
						</div>
					{{ end }}
				</div>
			{{ end }}
		</div>
	`)
	if err != nil {
		return err
	}

	templateVars := map[string]interface{}{
		"FileExtension": fileExtension,
		"Notebook":      notebook,
	}
	return t.Execute(writer, templateVars)
}
