package nbtohtml

import (
	"encoding/json"
)

// Output is the result of a code cell's execution in a Jupyter Notebook
type Output struct {
	Name       string   `json:"name"`
	OutputType string   `json:"output_type"`
	Text       []string `json:"text,omitempty"`
}

// Cell is a single Jupyter Notebook cell
type Cell struct {
	CellType       string   `json:"cell_type"`
	ExecutionCount *int     `json:"execution_count,omitempty"`
	Source         []string `json:"source"`
	Outputs        []Output `json:"outputs,omitempty"`
	// Omitted fields: "metadata"
}

// LanguageInfo provides details about the programming language of the Jupyter Notebook kernel
type LanguageInfo struct {
	FileExtension string `json:"file_extension"`
	// Omitted fields: codemirror_mode", "mimetype", "name", "nbconvert_exporter", "pygments_lexer",
	// "version"
}

// Metadata contains additional information about the Jupyter Notebook
type Metadata struct {
	LanguageInfo LanguageInfo `json:"language_info"`
	// Omitted fields: "kernelspec"
}

// Notebook represents the JSON data structure in which a Jupyter Notebook is stored
type Notebook struct {
	Cells    []Cell   `json:"cells"`
	Metadata Metadata `json:"metadata"`
	// Omitted fields: "nbformat", "nbformat_minor"
}

// parseNotebook takes the provided Jupyter Notebook JSON string and parses it into the
// corresponding structs.
func parseNotebook(notebookString string) (Notebook, error) {
	var notebookParsed Notebook
	err := json.Unmarshal([]byte(notebookString), &notebookParsed)
	return notebookParsed, err
}
