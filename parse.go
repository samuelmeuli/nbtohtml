package main

import (
	"encoding/json"
)

type output struct {
	Name       string   `json:"name"`
	OutputType string   `json:"output_type"`
	Text       []string `json:"text,omitempty"`
}

type cell struct {
	CellType       string   `json:"cell_type"`
	ExecutionCount *int     `json:"execution_count,omitempty"`
	Source         []string `json:"source"`
	Outputs        []output `json:"outputs,omitempty"`
	// Omitted fields: "metadata"
}

type languageInfo struct {
	FileExtension string `json:"file_extension"`
	// Omitted fields: codemirror_mode", "mimetype", "name", "nbconvert_exporter", "pygments_lexer",
	// "version"
}

type metadata struct {
	LanguageInfo languageInfo `json:"language_info"`
	// Omitted fields: "kernelspec"
}

// Jupyter Notebook JSON data structure
type notebook struct {
	Cells    []cell   `json:"cells"`
	Metadata metadata `json:"metadata"`
	// Omitted fields: "nbformat", "nbformat_minor"
}

// parseNotebook takes the provided Jupyter Notebook JSON string and parses it into the
// corresponding structs.
func parseNotebook(notebookString string) (notebook, error) {
	var notebookParsed notebook
	err := json.Unmarshal([]byte(notebookString), &notebookParsed)
	return notebookParsed, err
}
