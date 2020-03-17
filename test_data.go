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
							"<svg id=\"star\" xmlns=\"http://www.w3.org/2000/svg\" width=\"255\" height=\"240\" viewBox=\"0 0 51 48\">\n",
							"<path d=\"M25 1l6 17h18L35 29l5 17-15-10-15 10 5-17L1 18h18z\"/>\n",
							"</svg>"
						]
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"image/png": "base64-encoded-png-data"
					}
				},
				{
					"output_type": "display_data",
					"data": {
						"image/jpeg": "base64-encoded-jpeg-data"
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

var testMarkdownCell = Cell{
	CellType: "markdown",
	Source: []string{
		"# Hello World\n",
		"\n",
		"This is **bold** and _italic_",
	},
}

var testDisplayDataHTMLOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		TextHTML: []string{
			"<div>\n",
			"<p>Hello world</p>\n",
			"</div>",
		},
	},
}

var testPDFString = "base64-encoded-pdf-data"
var testDisplayDataPDFOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		ApplicationPDF: &testPDFString,
	},
}

var testLaTeXString = "latex-data"
var testDisplayDataLaTeXOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		TextLaTeX: &testLaTeXString,
	},
}

var testDisplayDataSVGOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		ImageSVGXML: []string{
			"<svg id=\"star\" xmlns=\"http://www.w3.org/2000/svg\" width=\"255\" height=\"240\" viewBox=\"0 0 51 48\">\n",
			"<path d=\"M25 1l6 17h18L35 29l5 17-15-10-15 10 5-17L1 18h18z\"/>\n",
			"</svg>",
		},
	},
}

var testPNGString = "base64-encoded-png-data"
var testDisplayDataPNGOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		ImagePNG: &testPNGString,
	},
}

var testJPEGString = "base64-encoded-jpeg-data"
var testDisplayDataJPEGOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		ImageJPEG: &testJPEGString,
	},
}

var testDisplayDataMarkdownOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		TextMarkdown: []string{
			"# Hello World\n",
			"\n",
			"This is **bold** and _italic_",
		},
	},
}

var testDisplayDataPlainTextOutput = Output{
	OutputType: "display_data",
	Data: OutputData{
		TextPlain: []string{
			"multiline\n",
			"text\n",
			"data",
		},
	},
}

var testErrorOutput = Output{
	OutputType: "error",
	Traceback: []string{
		"Error message",
		"With \u001b[0;31mANSI colors\u001b[0m",
	},
}

var testStreamOutput = Output{
	OutputType: "stream",
	Text: []string{
		"multiline\n",
		"stream\n",
		"text\n",
	},
}

var testExecutionCount1 = 1
var testExecuteResultOutput = Output{
	OutputType:     "execute_result",
	ExecutionCount: &testExecutionCount1,
	Data: OutputData{
		TextPlain: []string{
			"multiline\n",
			"text\n",
			"data",
		},
	},
}

var testExecutionCount2 = 2
var testCodeCell = Cell{
	CellType:       "code",
	ExecutionCount: &testExecutionCount2,
	Source: []string{
		"print(\"Hello\")\n",
		"print(\"World\")",
	},
	Outputs: []Output{
		testDisplayDataHTMLOutput,
		testDisplayDataPDFOutput,
		testDisplayDataLaTeXOutput,
		testDisplayDataSVGOutput,
		testDisplayDataPNGOutput,
		testDisplayDataJPEGOutput,
		testDisplayDataMarkdownOutput,
		testDisplayDataPlainTextOutput,
		testErrorOutput,
		testExecuteResultOutput,
		testStreamOutput,
	},
}

var testRawCell = Cell{
	CellType: "raw",
	Source: []string{
		"This is a raw section, without formatting.\n",
		"This is the second line.",
	},
}

var testMetadata = Metadata{
	LanguageInfo: LanguageInfo{
		FileExtension: ".py",
	},
}

var testParsedNotebook = Notebook{
	Cells: []Cell{
		testMarkdownCell,
		testCodeCell,
		testRawCell,
	},
	Metadata: testMetadata,
}
