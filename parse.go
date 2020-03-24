package nbtohtml

import (
	"encoding/json"
)

// Documentation of the Jupyter Notebook JSON format:
// - https://nbformat.readthedocs.io
// - https://ipython.org/ipython-doc/3/notebook/nbformat.html

// outputData can contain the cell output in various data types. Source:
// https://github.com/jupyter/nbconvert/blob/c837a22d44d98f6a58d1934bd85af1506df48f21/nbconvert/utils/base.py#L16.
type outputData struct {
	TextHTML       []string `json:"text/html,omitempty"`
	ApplicationPDF *string  `json:"application/pdf,omitempty"`
	TextLaTeX      *string  `json:"text/latex,omitempty"`
	ImageSVGXML    []string `json:"image/svg+xml,omitempty"`
	ImagePNG       *string  `json:"image/png,omitempty"`
	ImageJPEG      *string  `json:"image/jpeg,omitempty"`
	TextMarkdown   []string `json:"text/markdown,omitempty"`
	TextPlain      []string `json:"text/plain,omitempty"`
}

// output is the result of a code cell's execution in a Jupyter Notebook.
type output struct {
	OutputType     string     `json:"output_type"`
	ExecutionCount *int       `json:"execution_count,omitempty"`
	Text           []string   `json:"text,omitempty"`
	Data           outputData `json:"data,omitempty"`
	Traceback      []string   `json:"traceback,omitempty"`
	// Omitted fields: "ename", "evalue", "name"
}

// cell is a single Jupyter Notebook cell.
type cell struct {
	CellType       string   `json:"cell_type"`
	ExecutionCount *int     `json:"execution_count,omitempty"`
	Source         []string `json:"source"`
	Outputs        []output `json:"outputs,omitempty"`
	// Omitted fields: "metadata"
}

// languageInfo provides details about the programming language of the Jupyter Notebook kernel.
type languageInfo struct {
	FileExtension *string `json:"file_extension,omitempty"`
	// Omitted fields: "codemirror_mode", "mimetype", "name", "nbconvert_exporter", "pygments_lexer",
	// "version"
}

// kernelSpec provides details about the Jupyter Notebook kernel.
type kernelSpec struct {
	DisplayName *string `json:"display_name,omitempty"`
	Language    *string `json:"language,omitempty"`
	Name        *string `json:"name,omitempty"`
}

// metadata contains additional information about the Jupyter Notebook.
type metadata struct {
	LanguageInfo languageInfo `json:"language_info"`
	KernelSpec   kernelSpec   `json:"kernelspec"`
}

// notebook represents the JSON data structure in which a Jupyter Notebook is stored.
type notebook struct {
	Cells         []cell   `json:"cells"`
	Metadata      metadata `json:"metadata"`
	NBFormat      int      `json:"nbformat"`
	NBFormatMinor int      `json:"nbformat_minor"`
}

// parseNotebook takes the provided Jupyter Notebook JSON string and parses it into the
// corresponding structs.
func parseNotebook(notebookString string) (notebook, error) {
	var notebookParsed notebook
	err := json.Unmarshal([]byte(notebookString), &notebookParsed)
	return notebookParsed, err
}
