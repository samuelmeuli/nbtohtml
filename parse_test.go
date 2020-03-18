package nbtohtml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNotebook(t *testing.T) {
	actual, err := parseNotebook(testNotebookString)
	expected := testParsedNotebook
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
