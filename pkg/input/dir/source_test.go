package dir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceNonexistentPath(t *testing.T) {
	doesNotExist := "/dev/null/null"
	_, err := NewFileSource(doesNotExist)
	assert.Error(t, err, "should not create for nonexistent path")
}

func TestSourceValidPath(t *testing.T) {
	exists := "/"
	_, err := NewFileSource(exists)
	assert.NoError(t, err, "valid path")

	relative := "."
	_, err = NewFileSource(relative)
	assert.NoError(t, err, "valid path")
}
