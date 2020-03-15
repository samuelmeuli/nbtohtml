package nbtohtml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Documentation of the Jupyter Notebook JSON format: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// (VCS: https://github.com/ipython/ipython-doc/blob/e9c83570cf3dea6d7a6b178ee59869b4f441220f/3/notebook/nbformat.html)
const notebookString = `{
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
				"some *markdown*"
			]
		},
		{
			"cell_type": "code",
			"execution_count": 1,
			"metadata": {
				"collapsed": true,
				"autoscroll": false
			},
			"source": [
				"some code"
			],
			"outputs": [
				{
					"output_type": "stream",
					"name": "stdout",
					"text": [
						"multiline stream text"
					]
				},
				{
					"output_type": "display_data",
					"data": {
						"image/png": "base64-encoded-png-data"
					},
					"metadata": {
						"image/png": {
							"width": 640,
							"height": 480
						}
					}
				},
				{
					"output_type": "execute_result",
					"execution_count": 42,
					"data": {
						"text/plain": [
							"multiline text data"
						]
					}
				},
				{
          "output_type": "error",
					"ename": "Some error name",
					"evalue": "Some error value",
					"traceback": [
						"Trace part 1",
						"Trace part 2"
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
				"some nbformat mime-type data"
			]
		}
	]
}`

var base64PngString = "base64-encoded-png-data"
var executionCount1 = 1
var executionCount42 = 42
var expected = Notebook{
	Cells: []Cell{
		{
			CellType: "markdown",
			Source:   []string{"some *markdown*"},
		},
		{
			CellType:       "code",
			ExecutionCount: &executionCount1,
			Source: []string{
				"some code",
			},
			Outputs: []Output{
				{
					OutputType: "stream",
					Text:       []string{"multiline stream text"},
				},
				{
					OutputType: "display_data",
					Data: OutputData{
						ImagePNG: &base64PngString,
					},
				},
				{
					OutputType:     "execute_result",
					ExecutionCount: &executionCount42,
					Data: OutputData{
						TextPlain: []string{"multiline text data"},
					},
				},
				{
					OutputType: "error",
					Traceback: []string{
						"Trace part 1",
						"Trace part 2",
					},
				},
			},
		},
		{
			CellType: "raw",
			Source:   []string{"some nbformat mime-type data"},
		},
	},
	Metadata: Metadata{
		LanguageInfo: LanguageInfo{
			FileExtension: ".py",
		},
	},
}

func TestParseNotebook(t *testing.T) {
	actual, err := parseNotebook(notebookString)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
