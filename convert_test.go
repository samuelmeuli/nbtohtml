package nbtohtml

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStreamOutput(t *testing.T) {
	expected := template.HTML(`<pre>multiline
stream
text
</pre>`)
	actual := convertStreamOutput(testStreamOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertStreamOutputMissingKey(t *testing.T) {
	expected := template.HTML("")
	actual := convertStreamOutput(Output{OutputType: "stream"})
	assert.Equal(t, expected, actual)
}

func TestConvertDataHTMLOutput(t *testing.T) {
	expected := template.HTML(`<div>
<p>Hello world</p>
</div>`)
	actual := convertDataOutput(testDisplayDataHTMLOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataPDFOutput(t *testing.T) {
	expected := template.HTML("<pre>PDF output</pre>")
	actual := convertDataOutput(testDisplayDataPDFOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataLaTeXOutput(t *testing.T) {
	expected := template.HTML("<pre>LaTeX output</pre>")
	actual := convertDataOutput(testDisplayDataLaTeXOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataSVGOutput(t *testing.T) {
	expected := template.HTML(
		`<svg id="star" xmlns="http://www.w3.org/2000/svg" width="255" height="240" viewBox="0 0 51 48">
<path d="M25 1l6 17h18L35 29l5 17-15-10-15 10 5-17L1 18h18z"/>
</svg>`,
	)
	actual := convertDataOutput(testDisplayDataSVGOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataPNGOutput(t *testing.T) {
	expected := template.HTML(fmt.Sprintf(
		`<img src="data:image/png;base64,%s">`,
		*testDisplayDataPNGOutput.Data.ImagePNG,
	))
	actual := convertDataOutput(testDisplayDataPNGOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataJPEGOutput(t *testing.T) {
	expected := template.HTML(fmt.Sprintf(
		`<img src="data:image/jpeg;base64,%s">`,
		*testDisplayDataJPEGOutput.Data.ImageJPEG,
	))
	actual := convertDataOutput(testDisplayDataJPEGOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataMarkdownOutput(t *testing.T) {
	expected := template.HTML(`<h1>Hello World</h1>

<p>This is <strong>bold</strong> and <em>italic</em></p>
`)
	actual := convertDataOutput(testDisplayDataMarkdownOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataPlainTextOutput(t *testing.T) {
	expected := template.HTML(`<pre>multiline
text
data</pre>`)
	actual := convertDataOutput(testDisplayDataPlainTextOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputMissingKey(t *testing.T) {
	expected := template.HTML("")
	actual := convertDataOutput(Output{OutputType: "display_data"})
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutput(t *testing.T) {
	expected := template.HTML(`<pre>Error message
With <span class="term-fg31">ANSI colors</span></pre>`)
	actual := convertErrorOutput(testErrorOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutputMissingKey(t *testing.T) {
	expected := template.HTML("<pre>An unknown error occurred</pre>")
	actual := convertErrorOutput(Output{OutputType: "error"})
	assert.Equal(t, expected, actual)
}

func TestConvertMarkdownCell(t *testing.T) {
	expected := template.HTML(`<h1>Hello World</h1>

<p>This is <strong>bold</strong> and <em>italic</em></p>
`)
	actual := convertMarkdownCell(testMarkdownCell)
	assert.Equal(t, expected, actual)
}

func TestConvertCodeCell(t *testing.T) {
	expected := template.HTML(`(?s)<pre class="chroma">.*print.*Hello.*print.*World.*</pre>`)
	actual := convertCodeCell(testCodeCell, "py")
	assert.Regexp(t, expected, actual)
}

func TestConvertRawCell(t *testing.T) {
	expected := template.HTML(`<pre>This is a raw section, without formatting.
This is the second line.</pre>`)
	actual := convertRawCell(testRawCell)
	assert.Equal(t, expected, actual)
}
