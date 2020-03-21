package nbtohtml

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeHTML(t *testing.T) {
	expected := template.HTML(
		"&lt;p&gt;Hello world&lt;/p&gt;\n&lt;script&gt;window.alert(&#39;I&#39;m evil!&#39;);&lt;/script&gt;",
	)
	actual := escapeHTML(`<p>Hello world</p>
<script>window.alert('I'm evil!');</script>`)
	assert.Equal(t, expected, actual)
}

func TestSanitizeHTML(t *testing.T) {
	expected := template.HTML("<p>Hello world</p>\n")
	actual := sanitizeHTML(`<p>Hello world</p>
<script>window.alert('I'm evil!');</script>`)
	assert.Equal(t, expected, actual)
}
