package nbtohtml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseNotebook(t *testing.T) {
	actual, err := parseNotebook(testNotebookString)
	expected := testParsedNotebook
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
