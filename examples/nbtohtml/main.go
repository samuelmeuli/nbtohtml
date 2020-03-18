package main

import (
	"bytes"
	"html/template"
	"io"
	"os"

	"github.com/samuelmeuli/nbtohtml"
)

const notebookPath = "examples/nbtohtml/sample-notebook.ipynb"
const templatePath = "examples/nbtohtml/template.html"
const outputPath = "examples/nbtohtml/index.html"

type templateContent struct {
	NotebookHTML template.HTML
}

// checkError logs the error message and quits the program if an error is encountered.
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// fillTemplate converts the demo Jupyter Notebook to HTML and injects it into the template page.
func fillTemplate(writer io.Writer) error {
	tmpl := template.Must(template.ParseFiles(templatePath))
	notebookHTML := new(bytes.Buffer)
	err := nbtohtml.ConvertFile(notebookHTML, notebookPath)
	if err != nil {
		return err
	}
	templateContent := templateContent{NotebookHTML: template.HTML(notebookHTML.String())}
	return tmpl.Execute(writer, templateContent)
}

// main is a script for generating a demo page from the HTML template and sample Jupyter Notebook.
func main() {
	file, err := os.Create(outputPath)
	checkError(err)

	err = fillTemplate(file)
	checkError(err)

	err = file.Close()
	checkError(err)
}
