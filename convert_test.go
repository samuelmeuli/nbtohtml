// nolint:gosec
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

func TestConvertStreamOutputCodeInjection(t *testing.T) {
	expected := template.HTML(`<pre>multiline
stream
text
&lt;script&gt;window.alert(&#39;I&#39;m evil!&#39;);&lt;/script&gt;
</pre>`)
	actual := convertStreamOutput(testStreamOutputCodeInjection)
	assert.Equal(t, expected, actual)
}

func TestConvertStreamOutputMissingKey(t *testing.T) {
	expected := template.HTML("")
	actual := convertStreamOutput(testStreamOutputMissingKey)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputHTML(t *testing.T) {
	expected := template.HTML(`
<p>Hello world</p>
`)
	actual := convertDataOutput(testHTMLOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputHTMLCodeInjection(t *testing.T) {
	expected := template.HTML(`
<p>Hello world</p>
`)
	actual := convertDataOutput(testHTMLOutputCodeInjection)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputPDF(t *testing.T) {
	expected := template.HTML("<pre>PDF output</pre>")
	actual := convertDataOutput(testPDFOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputLaTeX(t *testing.T) {
	expected := template.HTML("<pre>LaTeX output</pre>")
	actual := convertDataOutput(testLaTeXOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputSVG(t *testing.T) {
	expected := template.HTML("<pre>SVG output</pre>")
	actual := convertDataOutput(testSVGOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputPNG(t *testing.T) {
	expected := template.HTML(fmt.Sprintf(`<img src="data:image/png;base64,%s">`, testPNGString))
	actual := convertDataOutput(testPNGOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputJPEG(t *testing.T) {
	expected := template.HTML(fmt.Sprintf(`<img src="data:image/jpeg;base64,%s">`, testJPEGString))
	actual := convertDataOutput(testJPEGOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputMarkdown(t *testing.T) {
	expected := template.HTML(`<h1>Hello World</h1>
<p>This is <strong>bold</strong> and <em>italic</em></p>
`)
	actual := convertDataOutput(testMarkdownOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputMarkdownCodeInjection(t *testing.T) {
	expected := template.HTML(`<h1>Hello World</h1>
<p>This is <strong>bold</strong> and <em>italic</em></p>
<!-- raw HTML omitted -->
`)
	actual := convertDataOutput(testMarkdownOutputCodeInjection)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputPlainText(t *testing.T) {
	expected := template.HTML(`<pre>multiline
text
data</pre>`)
	actual := convertDataOutput(testPlainTextOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputPlainTextCodeInjection(t *testing.T) {
	expected := template.HTML(`<pre>multiline
text
data
&lt;script&gt;window.alert(&#39;I&#39;m evil!&#39;);&lt;/script&gt;
</pre>`)
	actual := convertDataOutput(testPlainTextOutputCodeInjection)
	assert.Equal(t, expected, actual)
}

func TestConvertDataOutputMissingKey(t *testing.T) {
	expected := template.HTML("")
	actual := convertDataOutput(output{OutputType: "display_data"})
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutput(t *testing.T) {
	expected := template.HTML(`<pre>Error message
With <span class="term-fg31">ANSI colors</span></pre>`)
	actual := convertErrorOutput(testErrorOutput)
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutputCodeInjection(t *testing.T) {
	expected := template.HTML(`<pre>Error message
With <span class="term-fg31">ANSI colors</span>
&lt;script&gt;window.alert(&#39;I&#39;m evil!&#39;);&lt;&#47;script&gt;</pre>`)
	actual := convertErrorOutput(testErrorOutputCodeInjection)
	assert.Equal(t, expected, actual)
}

func TestConvertErrorOutputMissingKey(t *testing.T) {
	expected := template.HTML("<pre>An unknown error occurred</pre>")
	actual := convertErrorOutput(output{OutputType: "error"})
	assert.Equal(t, expected, actual)
}

func TestConvertMarkdownCell(t *testing.T) {
	expected := template.HTML(`<h1>Hello World</h1>
<p>This is <strong>bold</strong> and <em>italic</em></p>
`)
	actual := convertMarkdownCell(testMarkdownCell)
	assert.Equal(t, expected, actual)
}

func TestConvertMarkdownCellCodeInjection(t *testing.T) {
	expected := template.HTML(`<h1>Hello World</h1>
<p>This is <strong>bold</strong> and <em>italic</em></p>
<!-- raw HTML omitted -->
`)
	actual := convertMarkdownCell(testMarkdownCellCodeInjection)
	assert.Equal(t, expected, actual)
}

func TestConvertCodeCell(t *testing.T) {
	expected := template.HTML(`(?s)<pre class="chroma">.*print.*Hello.*print.*World.*</pre>`)
	actual := convertCodeCell(testCodeCell, "py")
	assert.Regexp(t, expected, actual)
}

func TestConvertCodeCellCodeInjection(t *testing.T) {
	expected := template.HTML(`(?s)<pre class="chroma">.*print.*Hello.*print.*World.*</pre>`)
	actual := convertCodeCell(testCodeCellCodeInjection, "py")
	assert.Regexp(t, expected, actual)
}

func TestConvertRawCell(t *testing.T) {
	expected := template.HTML(`<pre>This is a raw section, without formatting.
This is the second line.</pre>`)
	actual := convertRawCell(testRawCell)
	assert.Equal(t, expected, actual)
}

func TestConvertRawCellCodeInjection(t *testing.T) {
	expected := template.HTML(`<pre>This is a raw section, without formatting.
This is the second line.
&lt;script&gt;window.alert(&#39;I&#39;m evil!&#39;);&lt;/script&gt;</pre>`)
	actual := convertRawCell(testRawCellCodeInjection)
	assert.Equal(t, expected, actual)
}
