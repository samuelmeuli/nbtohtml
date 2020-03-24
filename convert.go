// Package nbtohtml is a library for converting Jupyter Notebook files to HTML.
package nbtohtml

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Output renderers

// convertDataOutput converts data output (e.g. a base64-encoded plot image) to HTML.
func convertDataOutput(output output) template.HTML {
	var outputHTML template.HTML = ""

	switch {
	case output.Data.TextHTML != nil:
		htmlString := strings.Join(output.Data.TextHTML, "")
		// Remove unnecessary wrapper <div>
		if strings.HasPrefix(htmlString, "<div>") && strings.HasSuffix(htmlString, "</div>") {
			htmlString = htmlString[5 : len(htmlString)-6]
		}
		outputHTML = sanitizeHTML(htmlString)
	case output.Data.ApplicationPDF != nil:
		// TODO: Implement PDF conversion
		fmt.Printf("missing conversion logic for `application/pdf` data type\n")
		outputHTML = "<pre>PDF output</pre>"
	case output.Data.TextLaTeX != nil:
		// TODO: Implement LaTeX conversion
		fmt.Printf("missing conversion logic for `text/latex` data type\n")
		outputHTML = "<pre>LaTeX output</pre>"
	case output.Data.ImageSVGXML != nil:
		// TODO: Implement SVG conversion
		fmt.Printf("missing conversion logic for `image/svg+xml` data type\n")
		outputHTML = "<pre>SVG output</pre>"
	case output.Data.ImagePNG != nil:
		htmlString := fmt.Sprintf(`<img src="data:image/png;base64,%s">`, *output.Data.ImagePNG)
		outputHTML = sanitizeHTML(htmlString)
	case output.Data.ImageJPEG != nil:
		htmlString := fmt.Sprintf(`<img src="data:image/jpeg;base64,%s">`, *output.Data.ImageJPEG)
		outputHTML = sanitizeHTML(htmlString)
	case output.Data.TextMarkdown != nil:
		outputHTML = renderMarkdown(output.Data.TextMarkdown)
	case output.Data.TextPlain != nil:
		escapedHTML := escapeHTML(strings.Join(output.Data.TextPlain, ""))
		outputHTML = "<pre>" + escapedHTML + "</pre>"
	default:
		fmt.Printf("missing `execute_result` data type in output of type `%s`\n", output.OutputType)
	}

	return outputHTML
}

// convertErrorOutput converts error output (e.g. generated by a Python exception) to HTML.
func convertErrorOutput(output output) template.HTML {
	if output.Traceback == nil {
		fmt.Printf("missing `traceback` key in output of type `error`\n")
		return "<pre>An unknown error occurred</pre>"
	}

	// Convert ANSI colors to HTML
	return renderTerminalOutput(output.Traceback)
}

// convertStreamOutput converts stream output (e.g. stdout written by a Python program) to HTML.
func convertStreamOutput(output output) template.HTML {
	if output.Text == nil {
		fmt.Printf("missing `text` key in output of type `stream`\n")
		return ""
	}

	escapedHTML := escapeHTML(strings.Join(output.Text, ""))
	return "<pre>" + escapedHTML + "</pre>"
}

// Cell renderers

// convertMarkdownCell converts a Markdown cell to HTML.
func convertMarkdownCell(cell cell) template.HTML {
	return renderMarkdown(cell.Source)
}

// convertCodeCell converts a code cell to HTML with classes for syntax highlighting.
func convertCodeCell(cell cell, languageID string) template.HTML {
	sourceString := strings.Join(cell.Source, "")
	cellHTML, err := renderSourceCode(sourceString, languageID)

	// Render code without syntax highlighting if an error occurred
	if err != nil {
		fmt.Printf("skipping syntax highlighting: %d\n", err)
		escapedHTML := escapeHTML(sourceString)
		return "<pre>" + escapedHTML + "</pre>"
	}

	return cellHTML
}

// convertRawCell returns a simple HTML element for the raw notebook cell.
func convertRawCell(cell cell) template.HTML {
	escapedHTML := escapeHTML(strings.Join(cell.Source, ""))
	return "<pre>" + escapedHTML + "</pre>"
}

// Input/output renderers

// convertPrompt returns an HTML string which indicates the input/output's execution count.
func convertPrompt(executionCount *int) template.HTML {
	if executionCount == nil {
		return ""
	}
	// Execution count is an integer, so HTML should be safe from code injection
	return template.HTML(fmt.Sprintf("[%d]:", *executionCount)) // nolint:gosec
}

// convertOutput converts the provided cell input to HTML.
func convertInput(languageID string, cell cell) template.HTML {
	switch cell.CellType {
	case "markdown":
		return convertMarkdownCell(cell)
	case "code":
		return convertCodeCell(cell, languageID)
	case "raw":
		return convertRawCell(cell)
	default:
		fmt.Printf("skipping cell (unrecognized cell type \"%s\")\n", cell.CellType)
		return ""
	}
}

// convertOutput converts the provided output of a cell execution to HTML.
func convertOutput(output output) template.HTML {
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
//
// For example, the function can be called the following way:
//
//  notebookHTML := new(bytes.Buffer)
//  notebookPath := "/path/to/your/notebook.ipynb"
//  err := nbtohtml.ConvertFile(notebookHTML, notebookPath)
func ConvertFile(writer io.Writer, notebookPath string) error {
	// Read file
	fileContent, err := ioutil.ReadFile(filepath.Clean(notebookPath))
	if err != nil {
		return fmt.Errorf("could not read Jupyter Notebook file at %s", notebookPath)
	}

	// Convert file content
	return ConvertString(writer, string(fileContent))
}

// ConvertString converts the provided Jupyter Notebook JSON string to HTML.
//
// For example, the function can be called the following way:
//
//  notebookHTML := new(bytes.Buffer)
//  notebookString := `{ "cells": ... }`
//  err := nbtohtml.ConvertString(notebookHTML, notebookString)
func ConvertString(writer io.Writer, notebookString string) error {
	notebook, err := parseNotebook(notebookString)
	if err != nil {
		return err
	}

	if notebook.NBFormat < 4 {
		return fmt.Errorf(
			"the provided Jupyter Notebook uses an old version of the Notebook file format (version 4 or higher is required)",
		)
	}

	// Try to find information about programming language used by the notebook kernel. Metadata fields
	// in the Jupyter Notebook JSON are optional, so multiple fields are checked
	languageID := ""
	if fileExtensionPtr := notebook.Metadata.LanguageInfo.FileExtension; fileExtensionPtr != nil {
		languageID = (*fileExtensionPtr)[1:]
	} else if kernelLanguagePtr := notebook.Metadata.KernelSpec.Language; kernelLanguagePtr != nil {
		languageID = *kernelLanguagePtr
	} else if kernelNamePtr := notebook.Metadata.KernelSpec.Name; kernelNamePtr != nil {
		languageID = *kernelNamePtr
	}

	t := template.New("notebook")
	t = t.Funcs(template.FuncMap{
		"convertPrompt": convertPrompt,
		"convertInput":  convertInput,
		"convertOutput": convertOutput,
		"getCellClasses": func(cell cell) string {
			return "cell cell-" + cell.CellType
		},
		"getOutputClasses": func(output output) string {
			outputTypeClass := strings.ReplaceAll(output.OutputType, "_", "-")
			return "output output-" + outputTypeClass
		},
	})
	t, err = t.Parse(`
		<div class="notebook">
			{{ $languageID := .languageID }}
			{{ range .notebook.Cells }}
				<div class="{{ . | getCellClasses }}">
					<div class="input-wrapper">
						<div class="input-prompt">
							{{ .ExecutionCount | convertPrompt }}
						</div>
						<div class="input">
							{{ . | convertInput $languageID }}
						</div>
					</div>
					{{ range .Outputs }}
						<div class="output-wrapper">
							<div class="output-prompt">
								{{ .ExecutionCount | convertPrompt }}
							</div>
							<div class="{{ . | getOutputClasses }}">
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
		"languageID": languageID,
		"notebook":   notebook,
	}
	return t.Execute(writer, templateVars)
}
