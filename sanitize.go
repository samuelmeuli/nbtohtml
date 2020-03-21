package nbtohtml

import (
	"html"
	"html/template"

	"github.com/microcosm-cc/bluemonday"
)

var sanitizerPolicy = initSanitizerPolicy()

func initSanitizerPolicy() *bluemonday.Policy {
	var p = bluemonday.UGCPolicy()
	p.AllowDataURIImages()
	return p
}

// escapeHTML escapes special HTML characters like "<" to become "&lt;".
func escapeHTML(htmlString string) template.HTML {
	escapedString := html.EscapeString(htmlString)
	return template.HTML(escapedString) // nolint:gosec
}

// sanitizeHTML takes a string that contains a HTML fragment or document and applies the given
// bluemonday policy.
func sanitizeHTML(htmlString string) template.HTML {
	sanitizedString := sanitizerPolicy.Sanitize(htmlString)
	return template.HTML(sanitizedString) // nolint:gosec
}
