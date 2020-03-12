package main

import (
	"fmt"
	"io/ioutil"
)

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

	// Build HTML string from converted cells
	html := ""
	for _, cell := range notebook.Cells {
		html += fmt.Sprintf("<p>%s/p>", cell.CellType) // TODO: Replace with conversion logic
	}

	return html, nil
}
