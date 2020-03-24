package nbtohtml

// Documentation of the Jupyter Notebook JSON format: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// (VCS: https://github.com/ipython/ipython-doc/blob/e9c83570cf3dea6d7a6b178ee59869b4f441220f/3/notebook/nbformat.html)
const testNotebookString = `{
	"metadata": {
		"kernelspec": {
			"display_name": "Python 3",
			"language": "python",
			"name": "python3"
		},
		"language_info": {
			"codemirror_mode": {
				"name": "ipython",
				"version": 3
			},
			"file_extension": ".py",
			"mimetype": "text/x-python",
			"name": "python",
			"nbconvert_exporter": "python",
			"pygments_lexer": "ipython3",
			"version": "3.7.6"
		}
	},
	"nbformat": 4,
	"nbformat_minor": 4,
	"cells": [
		{
			"cell_type": "markdown",
			"metadata": {},
			"source": [
				"# Hello World\n",
				"\n",
				"This is **bold** and _italic_"
			]
		},
		{
			"cell_type": "code",
			"execution_count": 2,
			"metadata": {
				"collapsed": true,
				"autoscroll": false
			},
			"source": [
				"print(\"Hello\")\n",
				"print(\"World\")"
			],
			"outputs": [
				{
					"output_type": "display_data",
					"data": {
						"text/html": [
							"<div>\n",
							"<p>Hello world</p>\n",
							"</div>"
						]
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"application/pdf": "base64-encoded-pdf-data"
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"text/latex": "latex-data"
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"image/svg+xml": [
							"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"255\" height=\"240\">\n",
							"<path d=\"M25 1l6 17h18L35 29l5 17-15-10-15 10 5-17L1 18h18z\"/>\n",
							"</svg>"
						]
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"image/png": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVQYV2NgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"image/jpeg": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVQYV2NgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"text/markdown": [
							"# Hello World\n",
							"\n",
							"This is **bold** and _italic_"
						]
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"text/plain": [
							"multiline\n",
							"text\n",
							"data"
						]
					}
				},
				{
          "output_type": "error",
					"ename": "Some error name",
					"evalue": "Some error value",
					"traceback": [
						"Error message",
						"With \u001b[0;31mANSI colors\u001b[0m"
					]
				},
				{
					"output_type": "execute_result",
					"execution_count": 1,
					"data": {
						"text/plain": [
							"multiline\n",
							"text\n",
							"data"
						]
					}
				},
				{
					"output_type": "stream",
					"name": "stdout",
					"text": [
						"multiline\n",
						"stream\n",
						"text\n"
					]
				}
			]
		},
		{
			"cell_type": "raw",
			"metadata": {
				"format": "mime/type"
			},
			"source": [
				"This is a raw section, without formatting.\n",
				"This is the second line."
			]
		}
	]
}`

var testMarkdownCell = cell{
	CellType: "markdown",
	Source: []string{
		"# Hello World\n",
		"\n",
		"This is **bold** and _italic_",
	},
}
var testMarkdownCellCodeInjection = cell{
	CellType: "markdown",
	Source: []string{
		"# Hello World",
		"\n",
		"This is **bold** and _italic_",
		"\n",
		"<script>window.alert('I'm evil!');</script>",
	},
}

var testHTMLOutput = output{
	OutputType: "display_data",
	Data: outputData{
		TextHTML: []string{
			"<div>\n",
			"<p>Hello world</p>\n",
			"</div>",
		},
	},
}
var testHTMLOutputCodeInjection = output{
	OutputType: "display_data",
	Data: outputData{
		TextHTML: []string{
			"<div>\n",
			"<p>Hello world</p>\n",
			"<script>window.alert('I'm evil!');</script>",
			"</div>",
		},
	},
}

var testPDFString = "base64-encoded-pdf-data"
var testPDFOutput = output{
	OutputType: "display_data",
	Data: outputData{
		ApplicationPDF: &testPDFString,
	},
}

var testLaTeXString = "latex-data"
var testLaTeXOutput = output{
	OutputType: "display_data",
	Data: outputData{
		TextLaTeX: &testLaTeXString,
	},
}

var testSVGOutput = output{
	OutputType: "display_data",
	Data: outputData{
		ImageSVGXML: []string{
			"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"255\" height=\"240\">\n",
			"<path d=\"M25 1l6 17h18L35 29l5 17-15-10-15 10 5-17L1 18h18z\"/>\n",
			"</svg>",
		},
	},
}

var testPNGString = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVQYV2NgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="
var testPNGOutput = output{
	OutputType: "display_data",
	Data: outputData{
		ImagePNG: &testPNGString,
	},
}

var testJPEGString = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVQYV2NgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="
var testJPEGOutput = output{
	OutputType: "display_data",
	Data: outputData{
		ImageJPEG: &testJPEGString,
	},
}

var testMarkdownOutput = output{
	OutputType: "display_data",
	Data: outputData{
		TextMarkdown: []string{
			"# Hello World\n",
			"\n",
			"This is **bold** and _italic_",
		},
	},
}
var testMarkdownOutputCodeInjection = output{
	OutputType: "display_data",
	Data: outputData{
		TextMarkdown: []string{
			"# Hello World",
			"\n",
			"This is **bold** and _italic_",
			"\n",
			"<script>window.alert('I'm evil!');</script>",
		},
	},
}

var testPlainTextOutput = output{
	OutputType: "display_data",
	Data: outputData{
		TextPlain: []string{
			"multiline\n",
			"text\n",
			"data",
		},
	},
}
var testPlainTextOutputCodeInjection = output{
	OutputType: "display_data",
	Data: outputData{
		TextPlain: []string{
			"multiline\n",
			"text\n",
			"data\n",
			"<script>window.alert('I'm evil!');</script>\n",
		},
	},
}

var testErrorOutput = output{
	OutputType: "error",
	Traceback: []string{
		"Error message",
		"With \u001b[0;31mANSI colors\u001b[0m",
	},
}
var testErrorOutputCodeInjection = output{
	OutputType: "error",
	Traceback: []string{
		"Error message",
		"With \u001b[0;31mANSI colors\u001b[0m",
		"<script>window.alert('I'm evil!');</script>",
	},
}

var testStreamOutput = output{
	OutputType: "stream",
	Text: []string{
		"multiline\n",
		"stream\n",
		"text\n",
	},
}
var testStreamOutputCodeInjection = output{
	OutputType: "stream",
	Text: []string{
		"multiline\n",
		"stream\n",
		"text\n",
		"<script>window.alert('I'm evil!');</script>\n",
	},
}
var testStreamOutputMissingKey = output{
	OutputType: "stream",
}

var testExecutionCount1 = 1
var testExecuteResultOutput = output{
	OutputType:     "execute_result",
	ExecutionCount: &testExecutionCount1,
	Data: outputData{
		TextPlain: []string{
			"multiline\n",
			"text\n",
			"data",
		},
	},
}

var testExecutionCount2 = 2
var testCodeCellOutputs = []output{
	testHTMLOutput,
	testPDFOutput,
	testLaTeXOutput,
	testSVGOutput,
	testPNGOutput,
	testJPEGOutput,
	testMarkdownOutput,
	testPlainTextOutput,
	testErrorOutput,
	testExecuteResultOutput,
	testStreamOutput,
}
var testCodeCell = cell{
	CellType:       "code",
	ExecutionCount: &testExecutionCount2,
	Source: []string{
		"print(\"Hello\")\n",
		"print(\"World\")",
	},
	Outputs: testCodeCellOutputs,
}
var testCodeCellCodeInjection = cell{
	CellType:       "code",
	ExecutionCount: &testExecutionCount2,
	Source: []string{
		"print(\"Hello\")\n",
		"print(\"World\")\n",
		"<script>window.alert('I'm evil!');</script>",
	},
	Outputs: testCodeCellOutputs,
}

var testRawCell = cell{
	CellType: "raw",
	Source: []string{
		"This is a raw section, without formatting.\n",
		"This is the second line.",
	},
}
var testRawCellCodeInjection = cell{
	CellType: "raw",
	Source: []string{
		"This is a raw section, without formatting.\n",
		"This is the second line.\n",
		"<script>window.alert('I'm evil!');</script>",
	},
}

var testMetadata = metadata{
	LanguageInfo: languageInfo{
		FileExtension: ".py",
	},
}

var testParsedNotebook = notebook{
	Cells: []cell{
		testMarkdownCell,
		testCodeCell,
		testRawCell,
	},
	Metadata:      testMetadata,
	NBFormat:      4,
	NBFormatMinor: 4,
}
