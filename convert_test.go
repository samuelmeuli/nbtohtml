package nbtohtml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertStreamOutputToHTML(t *testing.T) {
	expected := `<div class="output output-stream"><pre>multiline
stream
text
</pre></div>`
	actual, err := convertStreamOutputToHTML(testStreamOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertStreamOutputToHTMLMissingKey(t *testing.T) {
	expected := ""
	actual, err := convertStreamOutputToHTML(Output{OutputType: "stream"})
	assert.Errorf(t, err, "missing `text` key in output of type `stream`")
	assert.Equal(t, expected, actual)
}

func TestConvertDataHTMLOutputToHTML(t *testing.T) {
	expected := `<div class="output output-data-html"><div>
<p>Hello world</p>
</div></div>`
	actual, err := convertDataOutputToHTML(testDisplayDataHTMLOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertDataPDFOutputToHTML(t *testing.T) {
	expected := ""
	actual, err := convertDataOutputToHTML(testDisplayDataPDFOutput)
	assert.Errorf(t, err, "missing conversion logic for `application/pdf` data type")
	assert.Equal(t, expected, actual)
}

func TestConvertDataLaTeXOutputToHTML(t *testing.T) {
	expected := ""
	actual, err := convertDataOutputToHTML(testDisplayDataLaTeXOutput)
	assert.Errorf(t, err, "missing conversion logic for `text/latex` data type")
	assert.Equal(t, expected, actual)
}

func TestConvertDataSVGOutputToHTML(t *testing.T) {
	expected := `<div class="output output-data-svg"><svg id="star" xmlns="http://www.w3.org/2000/svg" width="255" height="240" viewBox="0 0 51 48">
<path d="M25 1l6 17h18L35 29l5 17-15-10-15 10 5-17L1 18h18z"/>
</svg></div>`
	actual, err := convertDataOutputToHTML(testDisplayDataSVGOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertDataPNGOutputToHTML(t *testing.T) {
	expected := fmt.Sprintf(
		`<div class="output output-data-png"><img src="data:image/png;base64,%s"></div>`,
		*testDisplayDataPNGOutput.Data.ImagePNG,
	)
	actual, err := convertDataOutputToHTML(testDisplayDataPNGOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertDataJPEGOutputToHTML(t *testing.T) {
	expected := fmt.Sprintf(
		`<div class="output output-data-jpeg"><img src="data:image/jpeg;base64,%s"></div>`,
		*testDisplayDataJPEGOutput.Data.ImageJPEG,
	)
	actual, err := convertDataOutputToHTML(testDisplayDataJPEGOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertDataMarkdownOutputToHTML(t *testing.T) {
	expected := `<div class="output output-data-markdown"><h1>Hello World</h1>

<p>This is <strong>bold</strong> and <em>italic</em></p>
</div>`
	actual, err := convertDataOutputToHTML(testDisplayDataMarkdownOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertDataPlainTextOutputToHTML(t *testing.T) {
	expected := `<div class="output output-data-plain-text"><pre>multiline
text
data</pre></div>`
	actual, err := convertDataOutputToHTML(testDisplayDataPlainTextOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputToHTMLMissingKey(t *testing.T) {
	expected := ""
	actual, err := convertDataOutputToHTML(Output{OutputType: "display_data"})
	assert.Errorf(t, err, "missing `text` key in output of type `stream`")
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutputToHTML(t *testing.T) {
	expected := `<div class="output output-error"><pre>Error message
With <span class="term-fg31">ANSI colors</span></pre></div>`
	actual, err := convertErrorOutputToHTML(testErrorOutput)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutputToHTMLMissingKey(t *testing.T) {
	expected := ""
	actual, err := convertErrorOutputToHTML(Output{OutputType: "error"})
	assert.Errorf(t, err, "missing `traceback` key in output of type `error`")
	assert.Equal(t, expected, actual)
}

func TestConvertMarkdownCellToHTML(t *testing.T) {
	expected := `<div class="cell cell-markdown"><h1>Hello World</h1>

<p>This is <strong>bold</strong> and <em>italic</em></p>
</div>`
	actual := convertMarkdownCellToHTML(testMarkdownCell)
	assert.Equal(t, expected, actual)
}

func TestConvertCodeCellToHTML(t *testing.T) {
	expected := `(?s)^<div class="cell cell-code">.*<pre class="chroma">.*print.*Hello.*print.*World.*</pre>.*</div><div class="output output-stream">.*</div>$`
	actual, err := convertCodeCellToHTML(testCodeCell, "py")
	assert.NoError(t, err)
	assert.Regexp(t, expected, actual)
}

func TestConvertRawCellToHTML(t *testing.T) {
	expected := `<div class="cell cell-raw"><pre>This is a raw section, without formatting.
This is the second line.</pre></div>`
	actual := convertRawCellToHTML(testRawCell)
	assert.Equal(t, expected, actual)
}
