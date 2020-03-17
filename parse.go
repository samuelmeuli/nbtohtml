package nbtohtml

import (
	"encoding/json"
)

// Documentation of the Jupyter Notebook JSON format: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// (VCS: https://github.com/ipython/ipython-doc/blob/e9c83570cf3dea6d7a6b178ee59869b4f441220f/3/notebook/nbformat.html)

// OutputData can contain the cell output in various data types.
// Source: https://github.com/jupyter/nbconvert/blob/c837a22d44d98f6a58d1934bd85af1506df48f21/nbconvert/utils/base.py#L16
type OutputData struct {
	TextHTML       []string `json:"text/html,omitempty"`
	ApplicationPDF *string  `json:"application/pdf,omitempty"`
	TextLaTeX      *string  `json:"text/latex,omitempty"`
	ImageSVGXML    []string `json:"image/svg+xml,omitempty"`
	ImagePNG       *string  `json:"image/png,omitempty"`
	ImageJPEG      *string  `json:"image/jpeg,omitempty"`
	TextMarkdown   []string `json:"text/markdown,omitempty"`
	TextPlain      []string `json:"text/plain,omitempty"`
}

// Output is the result of a code cell's execution in a Jupyter Notebook.
type Output struct {
	OutputType     string     `json:"output_type"`
	ExecutionCount *int       `json:"execution_count,omitempty"`
	Text           []string   `json:"text,omitempty"`
	Data           OutputData `json:"data,omitempty"`
	Traceback      []string   `json:"traceback,omitempty"`
	// Omitted fields: "ename", "evalue", "name"
}

// Cell is a single Jupyter Notebook cell.
type Cell struct {
	CellType       string   `json:"cell_type"`
	ExecutionCount *int     `json:"execution_count,omitempty"`
	Source         []string `json:"source"`
	Outputs        []Output `json:"outputs,omitempty"`
	// Omitted fields: "metadata"
}

// LanguageInfo provides details about the programming language of the Jupyter Notebook kernel.
type LanguageInfo struct {
	FileExtension string `json:"file_extension"`
	// Omitted fields: codemirror_mode", "mimetype", "name", "nbconvert_exporter", "pygments_lexer",
	// "version"
}

// Metadata contains additional information about the Jupyter Notebook.
type Metadata struct {
	LanguageInfo LanguageInfo `json:"language_info"`
	// Omitted fields: "kernelspec"
}

// Notebook represents the JSON data structure in which a Jupyter Notebook is stored.
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
